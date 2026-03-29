package node

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

// Message types for communication
type MessageType int

const (
	Handshake MessageType = iota
	Data
)

type NetworkMessage struct {
	Type    MessageType
	FromID  string
	Payload interface{}
}

type ChatPayload struct {
	Sender string
	Text   string
}

type listenerTCP struct {
	listenAddr string
	localNode  Node
	NodeCh     chan Node // Channel for incoming node discovery
	MessageCh  chan NetworkMessage
	// outgoingCh <-chan Node // Channel for nodes to connect to
	cm *ConnMap // Connection map
	// nm         *NodeMap    // Node map
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	ln     net.Listener
	mu     sync.RWMutex // Protects the maps during concurrent access
	onChat func(sender string, text string)
}

func NewListenerTCP(node Node) *listenerTCP {
	ctx, cancel := context.WithCancel(context.Background())
	return &listenerTCP{
		listenAddr: fmt.Sprintf("0.0.0.0:%d", TCPPort),
		localNode:  node,
		NodeCh:     make(chan Node),
		MessageCh:  make(chan NetworkMessage),
		cm:         NewConnMap(),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start begins both listening for incoming connections and processing outgoing connection requests
func (l *listenerTCP) Start() error {
	// Start TCP listener
	ln, err := net.Listen("tcp", l.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to start listener on %s: %w", l.listenAddr, err)
	}
	l.ln = ln

	log.Printf("TCP listener started on %s", l.listenAddr)

	// Start accepting incoming connections
	l.wg.Add(1)
	go l.acceptLoop()

	// Start processing outgoing connection requests
	l.wg.Add(1)
	go l.outgoingLoop()

	return nil
}

// acceptLoop handles incoming connection requests from other nodes
func (l *listenerTCP) acceptLoop() {
	defer l.wg.Done()

	for {
		select {
		case <-l.ctx.Done():
			return
		default:
			// Accept new connection
			conn, err := l.ln.Accept()
			if err != nil {
				select {
				case <-l.ctx.Done():
					return
				default:
					log.Printf("Accept error: %v", err)
					time.Sleep(100 * time.Millisecond) // Prevent tight loop on persistent errors
					continue
				}
			}

			// Handle incoming connection
			l.wg.Add(1)
			go l.handleIncomingConnection(conn)
		}
	}
}

// handleIncomingConnection processes a new incoming connection
func (l *listenerTCP) handleIncomingConnection(conn net.Conn) {
	defer l.wg.Done()
	pc := NewPeerConn(conn)

	remoteNode, err := l.receiveHandshake(pc)
	if err != nil {
		log.Printf("Handshake failed from %s: %v", conn.RemoteAddr(), err)
		conn.Close()
		return
	}

	l.mu.Lock()
	l.cm.Store(remoteNode.ID, pc)
	// l.nm.Store(conn, remoteNode.ID)
	l.mu.Unlock()

	log.Printf("New incoming connection from node %s (%s)", remoteNode.ID, conn.RemoteAddr())

	// select {
	// case l.incomingCh <- remoteNode:
	// default:
	// 	// Channel full, don't block
	// }

	// Start managing the connection
	l.manageConnection(remoteNode, pc)
}

// outgoingLoop processes nodes from the channel and initiates outgoing connections
func (l *listenerTCP) outgoingLoop() {
	defer l.wg.Done()

	for {
		select {
		case <-l.ctx.Done():
			return
		case node, ok := <-l.NodeCh:
			if !ok {
				// Channel closed
				return
			}

			// Check for duplicate before attempting connection
			l.mu.RLock()
			_, exists := l.cm.Load(node.ID)
			l.mu.RUnlock()

			if exists {
				log.Printf("Already connected to node %s, skipping outgoing connection", node.ID)
				continue
			}

			// Don't connect to self
			if node.ID == l.localNode.ID {
				log.Printf("Skipping connection to self")
				continue
			}

			// Initiate outgoing connection
			l.wg.Add(1)
			go l.initiateOutgoingConnection(node)
		}
	}
}

// initiateOutgoingConnection connects to a remote node
func (l *listenerTCP) initiateOutgoingConnection(node Node) {
	defer l.wg.Done()

	// Double-check for duplicate (in case it was added while waiting)
	l.mu.RLock()
	_, exists := l.cm.Load(node.ID)
	l.mu.RUnlock()

	if exists {
		log.Printf("Node %s already connected, aborting outgoing connection", node.ID)
		return
	}

	// Connect to the node
	address := net.JoinHostPort(node.Addr.IP.String(), strconv.Itoa(TCPPort))
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		log.Printf("Failed to connect to node %s at %s: %v", node.ID, address, err)
		return
	}
	pc := NewPeerConn(conn)

	// Set handshake deadline
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	// Send handshake with our info
	err = l.sendHandshake(pc)
	if err != nil {
		log.Printf("Handshake failed with node %s: %v", node.ID, err)
		conn.Close()
		return
	}

	// Reset deadline
	conn.SetDeadline(time.Time{})

	// Store connection (final check with lock)
	l.mu.Lock()
	// Check again under lock
	if _, exists := l.cm.Load(node.ID); exists {
		l.mu.Unlock()
		log.Printf("Node %s was connected by another goroutine", node.ID)
		conn.Close()
		return
	}

	l.cm.Store(node.ID, pc)
	l.mu.Unlock()

	log.Printf("Successfully connected to node %s (%s)", node.ID, address)

	// Start managing the connection
	l.manageConnection(node, pc)
}

// manageConnection handles an established connection (both incoming and outgoing)
func (l *listenerTCP) manageConnection(node Node, pc *PeerConn) {
	defer func() {
		pc.Conn.Close()

		l.mu.Lock()
		l.cm.Delete(node.ID)
		l.mu.Unlock()

		log.Printf("Connection to node %s closed", node.ID)
	}()

	// Message processing loop
	for {
		select {
		case <-l.ctx.Done():
			return
		default:
			// Set read deadline
			pc.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

			msg, err := l.receiveMessage(pc)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// Timeout - connection is idle but still alive
					continue
				}
				if err != io.EOF {
					log.Printf("Error reading from node %s: %v", node.ID, err)
				}
				return
			}

			// Reset deadline after successful read
			pc.Conn.SetReadDeadline(time.Time{})

			// Handle message
			l.handleMessage(node, msg)
		}
	}
}

// Helper functions for sending/receiving messages
func (l *listenerTCP) sendHandshake(pc *PeerConn) error {
	// Only send our own info (the requesting node)
	tmp := l.localNode
	tmp.Addr = pc.Conn.LocalAddr().(*net.TCPAddr)
	return l.sendMessage(pc, NetworkMessage{
		Type:    Handshake,
		FromID:  l.localNode.ID,
		Payload: tmp,
	})
}

func (l *listenerTCP) receiveHandshake(pc *PeerConn) (Node, error) {
	msg, err := l.receiveMessage(pc)
	if err != nil {
		return Node{}, err
	}

	if msg.Type != Handshake {
		return Node{}, fmt.Errorf("expected handshake, got %v", msg.Type)
	}

	node, ok := msg.Payload.(Node)
	if !ok {
		return Node{}, fmt.Errorf("invalid handshake payload")
	}

	// Update with actual remote address if Addr is nil or IP is empty
	if node.Addr == nil || node.Addr.IP == nil {
		host, portStr, _ := net.SplitHostPort(pc.Conn.RemoteAddr().String())
		port := TCPPort
		fmt.Sscanf(portStr, "%d", &port)
		node.Addr = &net.TCPAddr{IP: net.ParseIP(host), Port: port}
	}

	return node, nil
}

func (l *listenerTCP) sendMessage(pc *PeerConn, msg NetworkMessage) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	return pc.Enc.Encode(msg)
}

func (l *listenerTCP) receiveMessage(pc *PeerConn) (NetworkMessage, error) {
	var msg NetworkMessage
	err := pc.Dec.Decode(&msg)
	return msg, err
}

func (l *listenerTCP) handleMessage(from Node, msg NetworkMessage) {
	switch msg.Type {
	case Data:
		switch v := msg.Payload.(type) {
		case ChatPayload:
			if l.onChat != nil {
				l.onChat(v.Sender, v.Text)
			}
		case string:
			// Backward compatibility if some node sends plain strings.
			if l.onChat != nil {
				l.onChat(from.Name, v)
			}
		default:
			log.Printf("Unknown data payload from %s: %T", from.ID, msg.Payload)
		}
	default:
		// Ignore other message types here (handshake handled separately).
	}
}

func (l *listenerTCP) SetChatHandler(handler func(sender string, text string)) {
	l.onChat = handler
}

func (l *listenerTCP) SendChat(sender string, text string) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.cm.Range(func(key, value any) bool {
		if pc, ok := value.(*PeerConn); ok {
			msg := NetworkMessage{
				Type:    Data,
				FromID:  l.localNode.ID,
				Payload: ChatPayload{Sender: sender, Text: text},
			}
			if err := l.sendMessage(pc, msg); err != nil {
				log.Printf("Failed to send to %v: %v", key, err)
			}
		}
		return true
	})
}

// IsConnected checks if a node is currently connected
func (l *listenerTCP) IsConnected(node Node) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	_, exists := l.cm.Load(node.ID)
	return exists
}

// Stop gracefully shuts down the listener
func (l *listenerTCP) Stop() {
	l.cancel()

	if l.ln != nil {
		l.ln.Close()
	}

	// Close all connections
	l.mu.Lock()
	l.cm.Range(func(key, value interface{}) bool {
		if pc, ok := value.(*PeerConn); ok {
			pc.Conn.Close()
		}
		return true
	})
	l.mu.Unlock()

	// Wait for all goroutines to finish
	l.wg.Wait()
	log.Println("TCP listener stopped")
}

// Initialize in init() or main
func init() {
	gob.Register(Node{})
	gob.Register(map[string]interface{}{})
	gob.Register(ChatPayload{})
}

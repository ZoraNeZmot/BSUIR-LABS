package iup

import (
	"chat/node"
	"chat/startup"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	BroadPort = 56789
)

type IUP struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	addrLocal *net.UDPAddr
	c         chan node.Node
	seen      sync.Map
}

func (peer *IUP) SetChanel(c chan node.Node) {
	peer.c = c
}

func NewIUP(Name string) *IUP {
	bytes := make([]byte, 8)
	// TODO handle an error
	rand.Read(bytes)
	return &IUP{
		Name: Name,
		ID:   hex.EncodeToString(bytes),
	}
}

func (peer *IUP) GetNode() node.Node {
	return node.Node{Name: peer.Name, ID: peer.ID}
}

func CreateIUP(message []byte, addr *net.UDPAddr) (*IUP, error) {
	peer := &IUP{}
	err := json.Unmarshal(message, peer)
	if err != nil {
		return nil, err
	}

	peer.addrLocal = addr

	return peer, nil
}

func (peer *IUP) String() string {
	res, _ := json.Marshal(peer)
	return string(res)
	// return fmt.Sprintf("iup:Name-%s,ID-%s", peer.Name, peer.ID)
}

func (peer *IUP) GetLocalAddr() *net.UDPAddr {
	return peer.addrLocal
}

func (peer *IUP) listenIncoming(conn *net.UDPConn) {
	defer conn.Close()
	buffer := make([]byte, 2048)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		message := buffer[:n]

		newPeer, err := CreateIUP(message, remoteAddr)
		if err != nil {
			log.Println(err)
		}
		if peer.ID != newPeer.ID {
			if _, loaded := peer.seen.LoadOrStore(newPeer.ID, struct{}{}); loaded {
				continue
			}

			nn := newPeer.GetNode()
			nn.Addr = &net.TCPAddr{IP: remoteAddr.IP, Port: node.TCPPort}

			if peer.c != nil {
				select {
				case peer.c <- nn:
				default:
					// Drop if receiver is slow; discovery repeats anyway.
				}
			}

			log.Printf("New node: %s:%s\n", newPeer.Name, newPeer.ID)
			// Help the newcomer discover us as well.
			go peer.SendGreetingMessage()

		}

	}

}

func (peer *IUP) ListenIncoming() {
	// Keep backward compatibility: -p still exists, but the broadcast discovery must listen on BroadPort.
	_ = startup.FlagPort
	listenAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", BroadPort))
	if err != nil {
		log.Fatal("Failed to resolve listen address:", err)
	}

	listenConn, err := net.ListenUDP("udp4", listenAddr)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	go peer.listenIncoming(listenConn)
}

func (peer *IUP) sendBroad(message []byte) error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	broadcastAddr := "255.255.255.255"
	count := 0

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		if iface.Flags&net.FlagBroadcast == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			var ipnet *net.IPNet

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				ipnet = v
			case *net.IPAddr:
				ip = v.IP
				if ip.To4() != nil {
					_, ipnet, _ = net.ParseCIDR(ip.String() + "/24")
				}
			}

			if ip == nil || ip.To4() == nil {
				continue
			}

			var broadcastIP net.IP
			if ipnet != nil {
				mask := ipnet.Mask
				broadcastIP = make(net.IP, len(ip))
				j := 0
				for i := 0; i < len(ip); i++ {
					if i > 11 {
						broadcastIP[i] = ip[i] | ^mask[j]
						j++
					} else {
						broadcastIP[i] = ip[i]
					}
				}
			} else {
				broadcastIP = net.ParseIP(broadcastAddr)
			}

			udpAddr := &net.UDPAddr{
				IP:   broadcastIP,
				Port: BroadPort,
			}

			sendConn, err := net.DialUDP("udp4", nil, udpAddr)
			if err != nil {
				log.Printf("Failed to create UDP connection for interface %s: %v", iface.Name, err)
				continue
			}
			go func(conn *net.UDPConn, ifaceName string) {
				defer conn.Close()

				for i := 0; i < 3; i++ {
					_, err = conn.Write(message)
					if err != nil {
						log.Printf("Broadcast error on %s: %v", ifaceName, err)
					}
					time.Sleep(time.Second)
				}
			}(sendConn, iface.Name)

			count++
		}
	}

	if count == 0 {
		return errors.New("no available broadcast-capable interfaces")
	}

	log.Printf("Broadcasting message on %d interfaces", count)
	return nil
}

func (peer *IUP) SendGreetingMessage() {
	message, err := json.Marshal(peer)
	if err != nil {
		log.Printf("Failed to marshal greeting message: %v", err)
		return
	}

	if err := peer.sendBroad(message); err != nil {
		log.Printf("Failed to send broadcast: %v", err)
	}
}

package node

import (
	"encoding/gob"
	"net"
	"sync"
)

type PeerConn struct {
	Conn net.Conn
	Enc  *gob.Encoder
	Dec  *gob.Decoder
	mu   sync.Mutex
}

func NewPeerConn(conn net.Conn) *PeerConn {
	return &PeerConn{
		Conn: conn,
		Enc:  gob.NewEncoder(conn),
		Dec:  gob.NewDecoder(conn),
	}
}

type ConnMap struct {
	sync.Map
}

func (m *ConnMap) Store(nodeID string, value *PeerConn) {
	m.Map.Store(nodeID, value)
}

func (m *ConnMap) Load(nodeID string) (*PeerConn, bool) {
	c, ok := m.Map.Load(nodeID)
	if !ok {
		return nil, false
	}
	conn, ok := c.(*PeerConn)
	return conn, ok
}

func (m *ConnMap) Delete(nodeID string) {
	m.Map.Delete(nodeID)
}

func NewConnMap() *ConnMap {
	return &ConnMap{}
}

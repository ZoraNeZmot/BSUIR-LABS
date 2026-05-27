package node

import (
	"net"
	"sync"
)

type NodeMap struct {
	sync.Map
}

func (m *NodeMap) Store(key net.Conn, value string) {
	m.Map.Store(key, value)
}

func (m *NodeMap) Load(key net.Conn) (string, bool) {
	v, ok := m.Map.Load(key)
	return v.(string), ok
}

func NewNodeMap() *NodeMap {
	return &NodeMap{}
}

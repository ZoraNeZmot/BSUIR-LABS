package node

import (
	"fmt"
	"net"
)

const (
	TCPPort = 56790
)

type Node struct {
	Name string
	ID   string
	Addr *net.TCPAddr
}

func (node Node) GetIdentifier() string {
	return node.Name + "-" + node.ID
}

func (node Node) SplitIdentifier(iden string) (name string, id string) {
	fmt.Sscanf(iden, "%s-%s", name, id)
	node.Name = name
	node.ID = id
	return name, id
}

func NewNode(iden string, addr *net.TCPAddr) Node {
	res := Node{
		Addr: addr,
	}
	res.SplitIdentifier(iden)
	return res
}

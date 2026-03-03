package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8080,
	})
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	var buf []byte = make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(buf[:n]))
}

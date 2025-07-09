package main

import (
	"fmt"
	"net"
	"slices"
	"time"
)

func main() {
	go udpServerOrdering()
	time.Sleep(time.Second)
	udpClientOrdering()
}

// UDP server: buffers incoming packets and reorders them.
func udpServerOrdering() {
	addr, _ := net.ResolveUDPAddr("udp", ":9601")
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()
	received := make(map[uint64]string)
	for range 5 {
		buf := make([]byte, 1024)
		n, _, _ := conn.ReadFromUDP(buf)
		var seq uint64
		var msg string
		fmt.Sscanf(string(buf[:n]), "SEQ:%d MSG:%s", &seq, &msg)
		received[seq] = msg
	}
	// After receiving, print in order:
	var keys []uint64
	for k := range received {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Printf("In order: SEQ:%d MSG:%s\n", k, received[k])
	}
}

// Client: sends shuffled sequence packets
func udpClientOrdering() {
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9601")
	conn, _ := net.DialUDP("udp", nil, addr)
	defer conn.Close()
	type packet struct {
		seq uint64
		msg string
	}
	packets := []packet{
		{3, "foo"}, {0, "bar"}, {2, "baz"}, {4, "hi"}, {1, "bye"},
	}
	for _, pkt := range packets {
		msg := fmt.Sprintf("SEQ:%d MSG:%s", pkt.seq, pkt.msg)
		conn.Write([]byte(msg))
		time.Sleep(200 * time.Millisecond)
	}
}

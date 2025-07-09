package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	go udpReliableServer()
	time.Sleep(time.Second)
	udpReliableClient()
}

// Server: listens, expects increasing sequence numbers, sends ACKs.
func udpReliableServer() {
	addr, _ := net.ResolveUDPAddr("udp", ":9501")
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()
	seqExpected := uint64(0)
	buf := make([]byte, 4096)
	for {
		n, clientAddr, _ := conn.ReadFromUDP(buf)
		var seq uint64
		fmt.Sscanf(string(buf[:n]), "SEQ:%d", &seq)
		if seq == seqExpected {
			fmt.Printf("Server: received SEQ:%d\n", seq)
			seqExpected++
		}
		// Send ACK
		conn.WriteToUDP([]byte(fmt.Sprintf("ACK:%d", seq)), clientAddr)
	}
}

// Client: sends N packets, waits for ACK after each.
func udpReliableClient() {
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9501")
	conn, _ := net.DialUDP("udp", nil, addr)
	defer conn.Close()
	for seq := range 5 {
		msg := fmt.Sprintf("SEQ:%d", seq)
		conn.Write([]byte(msg))
		buf := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Timeout/no ACK for SEQ:", seq)
			continue
		}
		var ack uint64
		fmt.Sscanf(string(buf[:n]), "ACK:%d", &ack)
		fmt.Printf("Client: got %s\n", string(buf[:n]))
		time.Sleep(300 * time.Millisecond)
	}
}

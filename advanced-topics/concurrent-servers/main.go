package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

func main() {
	go startConcurrentTCP()
	go startConcurrentUDP()
	select {} // Block forever
}

// Concurrent TCP Server: each client handled by a goroutine.
func startConcurrentTCP() {
	ln, err := net.Listen("tcp", ":9201")
	if err != nil {
		fmt.Println("TCP listen error:", err)
		os.Exit(1)
	}
	fmt.Println("Concurrent TCP server on :9201")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("TCP accept error:", err)
			continue
		}
		go handleTCP(conn)
	}
}

func handleTCP(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("TCP read error:", err)
			}
			return
		}
		// Echo message
		conn.Write(buf[:n])
	}
}

// Concurrent UDP Server: worker pool (goroutines) for incoming packets
func startConcurrentUDP() {
	addr, _ := net.ResolveUDPAddr("udp", ":9202")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("UDP listen error:", err)
		os.Exit(1)
	}
	fmt.Println("Concurrent UDP server on :9202")
	const workers = 8
	packetChan := make(chan udpPacket, 100)
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pkt := range packetChan {
				conn.WriteToUDP(pkt.data, pkt.addr)
			}
		}()
	}
	for {
		buf := make([]byte, 4096)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		packetChan <- udpPacket{data: buf[:n], addr: addr}
	}
}

type udpPacket struct {
	data []byte
	addr *net.UDPAddr
}

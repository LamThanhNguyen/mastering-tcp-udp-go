package main

import (
	"flag"
	"fmt"
	"io"
	"net"
)

const (
	tcpAddr = "0.0.0.0:9101"
	udpAddr = "0.0.0.0:9102"
)

func main() {
	proto := flag.String("proto", "tcp", "tcp or udp")
	flag.Parse()

	if *proto == "tcp" {
		runTCPServer()
	} else {
		runUDPServer()
	}
}

func runTCPServer() {
	ln, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	fmt.Println("TCP server listening on", tcpAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			buf := make([]byte, 4096)
			for {
				n, err := c.Read(buf)
				if err != nil {
					if err != io.EOF {
						fmt.Println("TCP read error:", err)
					}
					return
				}
				// Echo back for RTT
				_, err = c.Write(buf[:n])
				if err != nil {
					fmt.Println("TCP write error:", err)
					return
				}
			}
		}(conn)
	}
}

func runUDPServer() {
	addr, err := net.ResolveUDPAddr("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("UDP server listening on", udpAddr)
	buf := make([]byte, 4096)
	for {
		n, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("UDP read error:", err)
			continue
		}
		// Echo back for RTT
		_, err = conn.WriteToUDP(buf[:n], remote)
		if err != nil {
			fmt.Println("UDP write error:", err)
		}
	}
}

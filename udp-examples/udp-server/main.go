package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	addr := net.UDPAddr{
		Port: 9001,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Failed to start UDP server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Printf("UDP server listening on %s\n", addr.String())

	buf := make([]byte, 2048)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("Error reading from UDP: %v\n", err)
			continue
		}

		msg := string(buf[:n])
		fmt.Printf("Received from %s: %s\n", remoteAddr.String(), msg)
		// Echo reply
		_, err = conn.WriteToUDP([]byte("Echo: "+msg), remoteAddr)
		if err != nil {
			fmt.Printf("Error sending response: %v\n", err)
		}
	}
}

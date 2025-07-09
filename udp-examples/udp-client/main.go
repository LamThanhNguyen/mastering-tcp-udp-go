package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9001")
	if err != nil {
		fmt.Printf("Failed to resolve address: %v\n", err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Printf("Failed to dial UDP server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connected to UDP server at", serverAddr.String())

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message (or 'quit'): ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		fmt.Printf("text: %v\n", text)
		if text == "quit" {
			break
		}
		_, err = conn.Write([]byte(text))
		if err != nil {
			fmt.Printf("Send error: %v\n", err)
			continue
		}
		buf := make([]byte, 2048)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("Read error: %v\n", err)
			continue
		}
		fmt.Printf("Server response: %s\n", string(buf[:n]))
	}
}

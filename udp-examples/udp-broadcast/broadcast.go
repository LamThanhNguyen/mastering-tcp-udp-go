package main

import (
	"fmt"
	"net"
	"os"
)

const (
	broadcastAddr = "255.255.255.255:9002"
	listenPort    = 9002
	maxMsgLen     = 2048
)

// ListenBroadcast listens for broadcasted chat messages and prints them.
func ListenBroadcast(selfName string) {
	addr := net.UDPAddr{
		Port: listenPort,
		IP:   net.IPv4zero,
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Failed to listen for broadcasts: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	buf := make([]byte, maxMsgLen)
	for {
		n, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		fmt.Printf("read from src: %v\n", src)
		msg, err := DecodeMessage(buf[:n])
		if err != nil {
			continue
		}
		if msg.From != selfName {
			fmt.Printf("\n[%s %s]: %s\n> ", msg.Time, msg.From, msg.Content)
		}
	}
}

// SendBroadcast sends a chat message via UDP broadcast.
func SendBroadcast(msg Message) error {
	bAddr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, bAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	data, err := EncodeMessage(msg)
	if err != nil {
		return err
	}
	_, err = conn.Write(data)
	return err
}

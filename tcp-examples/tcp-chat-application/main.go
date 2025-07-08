package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	port := "9000"
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
	defer ln.Close()
	log.Printf("Chat server started on port %s\n", port)

	chatRoom := NewChatRoom()
	go chatRoom.Run()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Connection error: %v\n", err)
			continue
		}
		go handleClient(conn, chatRoom)
	}
}

func handleClient(conn net.Conn, chatRoom *ChatRoom) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	conn.Write([]byte("Enter your username: "))
	username, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	username = strings.TrimSpace(username)
	if username == "" {
		username = fmt.Sprintf("User-%d", time.Now().UnixNano()%10000)
	}

	client := &Client{
		Conn:     conn,
		Username: username,
		Outbound: make(chan Message, 10),
	}

	chatRoom.Register <- client

	// Writer Goroutine: send outbound messages to this client
	go func() {
		encoder := json.NewEncoder(conn)
		for msg := range client.Outbound {
			encoder.Encode(msg)
		}
	}()

	// Reader loop: receive messages from client and broadcast
	decoder := json.NewDecoder(reader)
	for {
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			log.Printf("User %s disconnected or sent invalid data.", client.Username)
			break
		}
		msg.From = client.Username
		chatRoom.Broadcast <- msg
	}
	chatRoom.Unregister <- client
}

package main

import (
	"log"
	"net"
	"sync"
)

type Client struct {
	Conn     net.Conn
	Username string
	Outbound chan Message
}

type ChatRoom struct {
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (cr *ChatRoom) Run() {
	for {
		select {
		case client := <-cr.Register:
			cr.mu.Lock()
			cr.Clients[client] = true
			cr.mu.Unlock()
			log.Printf("User joined: %s\n", client.Username)
			cr.sendSystemMessage(client, "Welcome to the chat, "+client.Username+"!")
		case client := <-cr.Unregister:
			cr.mu.Lock()
			if _, ok := cr.Clients[client]; ok {
				close(client.Outbound)
				delete(cr.Clients, client)
				log.Printf("User left: %s\n", client.Username)
			}
			cr.mu.Unlock()
		case msg := <-cr.Broadcast:
			cr.mu.Lock()
			for client := range cr.Clients {
				// Don't send message back to sender
				if msg.From != client.Username {
					select {
					case client.Outbound <- msg:
					default:
						// Drop message if outbound is full (client slow)
					}
				}
			}
			cr.mu.Unlock()
		}
	}
}

func (cr *ChatRoom) sendSystemMessage(client *Client, text string) {
	sysMsg := Message{
		From:    "System",
		Content: text,
		Time:    NowStr(),
	}
	client.Outbound <- sysMsg
}

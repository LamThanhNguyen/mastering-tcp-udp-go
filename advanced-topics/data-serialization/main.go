package main

import (
	"encoding/json"
	"fmt"
	"net"
)

// Message for serialization (use JSON or protobuf)
type Message struct {
	From    string `json:"from"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

func main() {
	go startJSONServer()
	select {}
}

func startJSONServer() {
	ln, _ := net.Listen("tcp", ":9301")
	fmt.Println("JSON TCP server on :9301")
	for {
		conn, _ := ln.Accept()
		go handleJSON(conn)
	}
}

func handleJSON(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)
	for {
		var msg Message
		if err := dec.Decode(&msg); err != nil {
			fmt.Println("JSON decode error:", err)
			return
		}
		fmt.Printf("Received: %+v\n", msg)
		_ = enc.Encode(msg) // echo back
	}
}

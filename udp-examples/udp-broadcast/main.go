package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Print("Enter your username: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if username == "" {
		username = "user"
	}

	// Start broadcast listener in a goroutine
	go ListenBroadcast(username)

	fmt.Println("Type messages and hit Enter to broadcast. Type 'quit' to exit.")
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		if text == "quit" {
			fmt.Println("Exiting.")
			break
		}
		msg := Message{
			From:    username,
			Content: text,
			Time:    NowStr(),
		}
		if err := SendBroadcast(msg); err != nil {
			fmt.Printf("Send error: %v\n", err)
		}
	}
}

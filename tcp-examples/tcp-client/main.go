package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("Hello from TCP client!"))
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}
	fmt.Printf("Server response: %s\n", string(buf[:n]))
}

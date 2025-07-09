package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	ln, _ := net.Listen("tcp", ":9401")
	fmt.Println("TCP server (graceful shutdown) on :9401")

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		fmt.Println("Shutdown signal received")
		cancel()
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				fmt.Println("Server closed")
				wg.Wait()
				os.Exit(0)
			default:
				fmt.Println("Accept error:", err)
			}
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			handle(conn)
		}()
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 4096)
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		fmt.Printf("%v", n)
	}
}

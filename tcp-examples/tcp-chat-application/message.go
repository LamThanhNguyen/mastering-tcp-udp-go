package main

import "time"

type Message struct {
	From    string `json:"from"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

func NowStr() string {
	return time.Now().Format("15:04:05")
}

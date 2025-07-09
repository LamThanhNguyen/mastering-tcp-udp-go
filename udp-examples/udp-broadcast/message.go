package main

import (
	"encoding/json"
	"time"
)

type Message struct {
	From    string `json:"from"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

// NowStr returns the current time formatted as a string.
func NowStr() string {
	return time.Now().Format("15:04:05")
}

// EncodeMessage marshals a Message to JSON bytes.
func EncodeMessage(msg Message) ([]byte, error) {
	return json.Marshal(msg)
}

// DecodeMessage unmarshals JSON bytes to a Message.
func DecodeMessage(data []byte) (Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return msg, err
}

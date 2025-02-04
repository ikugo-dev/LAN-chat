package main

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type     string            `json:"type"`     // "chat", "vote", "file" ...
	Payload  json.RawMessage   `json:"payload"`  // The actual message content
	Metadata map[string]string `json:"metadata"` // Extra info (e.g., filename)
}

const (
	writeWait      = 5 * time.Minute
	pongWait       = 1 * time.Minute
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte // JSON format
}

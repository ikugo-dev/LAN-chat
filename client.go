package main

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type     string            `json:"type"`     // e.g., "chat", "vote", "file"
	Payload  json.RawMessage   `json:"payload"`  // the actual message content
	Metadata map[string]string `json:"metadata"` // Extra info (e.g., filename)
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
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

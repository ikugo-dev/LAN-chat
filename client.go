package main

import (
	"encoding/json"
	"log"
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
	send chan []byte
}

func (c *Client) writeMessage(msgType int, message []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	w, err := c.conn.NextWriter(msgType)
	if err != nil {
		return err
	}
	defer w.Close()

	w.Write(message)

	// Add queued chat messages to the current WebSocket message.
	n := len(c.send)
	for i := 0; i < n; i++ {
		w.Write(newline)
		w.Write(<-c.send)
	}
	return nil
}

// hub -> websocket
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok { // The hub closed the channel.
				c.writeMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.writeMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// websocket -> hub
func (c *Client) readPump(handler func(msg Message, c *Client)) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}
		var msg Message

		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}
		log.Printf("Received raw message: %s", message)
		log.Printf("Received unmarshaled: %+v", msg)

		handler(msg, c)
	}
}

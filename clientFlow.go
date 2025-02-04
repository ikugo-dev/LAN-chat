package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// hub -> websocket
func (c *Client) writePump() { // Sends data
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

// websocket -> hub
func (c *Client) readPump() { // Receives data
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
		// log.Printf("Received raw message: %s", message)
		// log.Printf("Received unmarshaled: %+v", msg)

		c.hub.broadcast <- message
	}
}

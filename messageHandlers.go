package main

import (
	"encoding/json"
	"log"
	"os"
)

func handleTextMessage(msg Message, c *Client) {
	if msg.Type != "chat" {
		log.Printf("Wrong message type: chat != %v", msg.Type)
		return
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	c.hub.broadcast <- data
}

func handleVoteMessage(msg Message, c *Client) {
	log.Printf("Received vote: %s", msg.Payload)
	c.hub.broadcast <- msg.Payload
}

func handleFileMessage(msg Message, c *Client) {
	filename, ok := msg.Metadata["filename"]
	if !ok {
		log.Println("Error: No filename provided in metadata")
		return
	}

	err := os.WriteFile("uploads/"+filename, msg.Payload, 0644)
	if err != nil {
		log.Printf("Failed to save file: %v", err)
		return
	}

	log.Printf("File received: %s", filename)
	c.hub.broadcast <- msg.Payload
}

func handleBinaryFileMessage(data []byte, c *Client) {
	filename := "uploaded_file.dat" // Generate a proper filename

	err := os.WriteFile("uploads/"+filename, data, 0644)
	if err != nil {
		log.Printf("Failed to save binary file: %v", err)
		return
	}

	log.Printf("Binary file received: %s", filename)
}

package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func pseudo_uuid() (uuid string) {
	// Generate randoms string
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	uuid = fmt.Sprintf("%04x-%04x-%04x-%04x-%04x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func main() {
	n := maelstrom.NewNode()

	// Register a handle for the "unique-ids" messages that responds with a unique id
	n.Handle("unique-ids", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type
		body["type"] = "generate_ok"

		// Updat the message with id
		body["id"] = pseudo_uuid()

		// Echo original message back with updated message type
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
}

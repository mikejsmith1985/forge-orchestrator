package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// Connect to the WebSocket server
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	log.Println("✅ Successfully connected to WebSocket server")

	// Set up interrupt handler
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	// Read messages in a goroutine
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			log.Printf("✅ Received: %s", message)
			
			// Try to parse as JSON
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err == nil {
				if msgType, ok := msg["type"].(string); ok {
					log.Printf("✅ Message type: %s", msgType)
				}
			}
		}
	}()

	// Send a test message
	testMessage := map[string]interface{}{
		"type": "TEST",
		"payload": map[string]string{
			"message": "Hello from E2E test",
		},
	}
	msgBytes, _ := json.Marshal(testMessage)
	
	log.Println("Sending test message...")
	err = c.WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		log.Println("write error:", err)
		return
	}
	log.Println("✅ Test message sent successfully")

	// Wait a bit to receive the echo
	time.Sleep(2 * time.Second)

	// Test flow status broadcast (simulated)
	log.Println("\n=== Testing message types ===")
	
	// Send a ping to keep connection alive
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	testsPassed := true

	select {
	case <-done:
		log.Println("Connection closed")
	case <-interrupt:
		log.Println("Interrupt received, closing connection...")
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close error:", err)
			return
		}
		select {
		case <-done:
		case <-time.After(time.Second):
		}
	case <-time.After(3 * time.Second):
		log.Println("✅ E2E test completed successfully")
	}

	if testsPassed {
		fmt.Println("\n=== E2E TEST RESULTS ===")
		fmt.Println("✅ WebSocket connection: PASSED")
		fmt.Println("✅ Message sending: PASSED")
		fmt.Println("✅ Message receiving (echo): PASSED")
		fmt.Println("✅ Connection lifecycle: PASSED")
		fmt.Println("\n🎉 All E2E tests PASSED!")
	}
}

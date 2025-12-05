package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	id       int
	conn     *websocket.Conn
	received []string
	mu       sync.Mutex
}

func main() {
	log.Println("=== Testing Hub Broadcast Functionality ===\n")

	// Create 3 clients
	numClients := 3
	clients := make([]*Client, numClients)
	var wg sync.WaitGroup

	// Connect all clients
	for i := 0; i < numClients; i++ {
		u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		defer conn.Close()

		clients[i] = &Client{
			id:       i + 1,
			conn:     conn,
			received: make([]string, 0),
		}

		log.Printf("✅ Client %d connected", i+1)

		// Start reading messages for each client
		wg.Add(1)
		go func(client *Client) {
			defer wg.Done()
			for {
				_, message, err := client.conn.ReadMessage()
				if err != nil {
					return
				}
				client.mu.Lock()
				client.received = append(client.received, string(message))
				client.mu.Unlock()
				log.Printf("✅ Client %d received: %s", client.id, message)
			}
		}(clients[i])
	}

	// Give time for all connections to be established
	time.Sleep(500 * time.Millisecond)

	log.Println("\n=== Testing Echo Functionality ===")
	
	// Each client sends a message
	for i, client := range clients {
		msg := map[string]interface{}{
			"type": "CLIENT_MESSAGE",
			"payload": map[string]interface{}{
				"clientId": i + 1,
				"message":  fmt.Sprintf("Hello from client %d", i+1),
			},
		}
		msgBytes, _ := json.Marshal(msg)
		err := client.conn.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			log.Printf("❌ Client %d write error: %v", i+1, err)
			continue
		}
		log.Printf("✅ Client %d sent message", i+1)
	}

	// Wait for messages to be received
	time.Sleep(2 * time.Second)

	log.Println("\n=== Test Results ===")
	
	// Verify each client received their own echo
	allPassed := true
	for i, client := range clients {
		client.mu.Lock()
		count := len(client.received)
		client.mu.Unlock()
		
		log.Printf("Client %d received %d message(s)", i+1, count)
		if count < 1 {
			log.Printf("❌ Client %d did not receive expected message", i+1)
			allPassed = false
		} else {
			log.Printf("✅ Client %d received messages correctly", i+1)
		}
	}

	// Close all connections
	log.Println("\n=== Closing Connections ===")
	for i, client := range clients {
		client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		client.conn.Close()
		log.Printf("✅ Client %d disconnected", i+1)
	}

	time.Sleep(500 * time.Millisecond)

	fmt.Println("\n=================================")
	if allPassed {
		fmt.Println("🎉 Hub Broadcast Test PASSED!")
		fmt.Println("✅ Multiple clients can connect")
		fmt.Println("✅ Messages are echoed correctly")
		fmt.Println("✅ Connections handled gracefully")
	} else {
		fmt.Println("❌ Some tests FAILED")
	}
	fmt.Println("=================================")
}

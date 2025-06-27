package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type Ack struct {
	MessageID int64  `json:"messageId"`
	Status    string `json:"status"`
}

func main() {
	user := flag.String("user", "", "Username of the client (required)")
	recipient := flag.String("recipient", "", "Username of the recipient (required)")
	server := flag.String("server", "ws://localhost:8080/ws", "WebSocket server URL")
	flag.Parse()

	if *user == "" {
		log.Fatal("Missing required -user flag")
	}

	if *recipient == "" {
		log.Fatal("Missing required -recipient flag")
	}

	url := fmt.Sprintf("%s?user=%s", *server, *user)

	log.Printf("🔌 Connecting to %s as user '%s'...\n", url, *user)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("❌ Connection error: %v", err)
	}
	defer conn.Close()

	// Handle OS interrupt signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Graceful shutdown coordination
	done := make(chan struct{})

	// Reader goroutine
	go func() {
		defer close(done)
		for {
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				log.Printf("📴 Connection closed or error: %v", err)
				return
			}

			// Try to parse as Message or Ack
			var msg Message
			if err := json.Unmarshal(msgBytes, &msg); err == nil && msg.Content != "" {
				log.Printf("📥 [%s] %s", msg.From, msg.Content)
				continue
			}

			var ack Ack
			if err := json.Unmarshal(msgBytes, &ack); err == nil {
				log.Printf("✅ Ack received for msg ID %d", ack.MessageID)
			}
		}
	}()

	// Sender loop
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if !scanner.Scan() {
				break
			}

			input := scanner.Text()
			parts := strings.Fields(input)
			if len(parts) < 1 {
				log.Println("⚠️ Invalid input format. Use: <message>")
				continue
			}

			content := strings.Join(parts[:], " ")

			msg := Message{
				From:      *user,
				To:        *recipient,
				Content:   content,
				Timestamp: time.Now().UnixNano(),
			}
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("⚠️  Failed to send: %v", err)
				return
			}
		}
	}()

	select {
	case <-interrupt:
		log.Println("👋 Interrupt received. Disconnecting...")
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Client closed"))
		time.Sleep(time.Second)
	case <-done:
		log.Println("🔌 Server disconnected.")
	}
}

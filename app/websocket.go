package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(store *MessageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.URL.Query().Get("user")
		if user == "" {
			http.Error(w, "Missing user param", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			return
		}
		defer conn.Close()

		session := &UserSession{Send: make(chan Message, 100)}
		store.AddSession(user, session)
		defer store.RemoveSession(user)

		// Writer goroutine
		go func() {
			for msg := range session.Send {
				if err := conn.WriteJSON(msg); err != nil {
					log.Printf("Error writing to %s: %v", user, err)
					return
				}
			}
		}()

		// Reader loop
		for {
			var msg Message
			if err := conn.ReadJSON(&msg); err != nil {
				log.Printf("Error reading from %s: %v", user, err)
				break
			}
			store.SendMessage(msg)

			ack := Ack{MessageID: msg.Timestamp, Status: "delivered"}
			if err := conn.WriteJSON(ack); err != nil {
				log.Printf("Failed to send ack to %s: %v", user, err)
			}
		}
	}
}

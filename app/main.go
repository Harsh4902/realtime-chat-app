package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	store := NewMessageStore()

	http.HandleFunc("/ws", wsHandler(store))
	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		user1 := r.URL.Query().Get("user1")
		user2 := r.URL.Query().Get("user2")
		if user1 == "" || user2 == "" {
			http.Error(w, "Missing user1 or user2", http.StatusBadRequest)
			return
		}
		history := store.GetHistory(user1, user2)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(history)
	})

	log.Println("Chat service started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

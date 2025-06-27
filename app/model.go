package main

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

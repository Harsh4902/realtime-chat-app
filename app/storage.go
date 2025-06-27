package main

import (
	"log"
	"sync"
	"time"
)

type UserSession struct {
	Send chan Message
}

type MessageStore struct {
	sessions map[string]*UserSession
	history  map[string][]Message
	buffer   map[string][]Message
	lock     sync.RWMutex
}

func NewMessageStore() *MessageStore {
	return &MessageStore{
		sessions: make(map[string]*UserSession),
		history:  make(map[string][]Message),
		buffer:   make(map[string][]Message),
	}
}

func (s *MessageStore) AddSession(user string, session *UserSession) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sessions[user] = session
	log.Printf("User connected: %s", user)

	// Deliver buffered messages if any
	if msgs, ok := s.buffer[user]; ok {
		for _, msg := range msgs {
			log.Printf("Delivering buffered message to %s: %s", user, msg.Content)
			session.Send <- msg
		}
		delete(s.buffer, user)
	}
}

func (s *MessageStore) RemoveSession(user string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.sessions, user)
	log.Printf("User disconnected: %s", user)
}

func (s *MessageStore) SendMessage(msg Message) {
	s.lock.Lock()
	defer s.lock.Unlock()

	msg.Timestamp = time.Now().UnixNano()

	key := chatKey(msg.From, msg.To)
	s.history[key] = append(s.history[key], msg)
	log.Printf("Message from %s to %s: %s", msg.From, msg.To, msg.Content)

	if session, ok := s.sessions[msg.To]; ok {
		log.Printf("Delivering message to %s", msg.To)
		session.Send <- msg
	} else {
		log.Printf("Buffering message for offline user %s", msg.To)
		s.buffer[msg.To] = append(s.buffer[msg.To], msg)
	}
}

func (s *MessageStore) GetHistory(user1, user2 string) []Message {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.history[chatKey(user1, user2)]
}

func chatKey(u1, u2 string) string {
	if u1 < u2 {
		return u1 + ":" + u2
	}
	return u2 + ":" + u1
}

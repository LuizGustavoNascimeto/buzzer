package domain

import "time"

// Message é a Entidade de negócio. Não sabe se vem de DynamoDB, Postgres ou HTTP.
type Message struct {
	ID          string
	GroupID     string
	SenderID    string
	DisplayName string
	Handle      string
	Content     string
	SentAt      time.Time
}

type MessageGroup struct {
	ID          string
	UserID      string
	DisplayName string
	Handle      string
	Content     string
	LastMessage string
	LastSentAt  time.Time
}

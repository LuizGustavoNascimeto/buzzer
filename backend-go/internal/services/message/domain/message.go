package domain

import "time"

type CreateMessage struct {
	GroupID        string `json:"group_id"`
	UserID         string `json:"user_id"`
	DisplayName    string `json:"display_name"`
	Handle         string `json:"handle"`
	Message        string `json:"message"`
	MessageGroupID string `json:"message_group_id"`
}

type CreateMessageGroup struct {
	UserID         string `json:"user_id"`
	DisplayName    string `json:"display_name"`
	Handle         string `json:"handle"`
	Message        string `json:"message"`
	MessageGroupID string `json:"message_group_id"`
}

type MessageGroup struct {
	ID          string    `dynamodbav:"message_group_uuid" json:"message_group_uuid"`
	DisplayName string    `dynamodbav:"user_display_name" json:"user_display_name"`
	Handle      string    `dynamodbav:"user_handle" json:"user_handle"`
	Message     string    `dynamodbav:"message" json:"message"`
	LastSentAt  time.Time `dynamodbav:"sk" json:"created_at"`
}

type Message struct {
	ID          string    `dynamodbav:"message_uuid" json:"message_uuid"`
	DisplayName string    `dynamodbav:"user_display_name" json:"user_display_name"`
	Handle      string    `dynamodbav:"user_handle" json:"user_handle"`
	Message     string    `dynamodbav:"message" json:"message"`
	LastSentAt  time.Time `dynamodbav:"sk" json:"created_at"`
}

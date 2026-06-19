package infra

import (
	"backend-go/internal/services/message/domain"
	"time"
)

type messageModel struct {
	ID               string    `dynamodbav:"message_uuid"`
	OtherUserID      string    `dynamodbav:"user_uuid"`
	OtherDisplayName string    `dynamodbav:"user_display_name"`
	OtherHandle      string    `dynamodbav:"user_handle"`
	Content          string    `dynamodbav:"message"`
	LastSentAt       time.Time `dynamodbav:"sk"`
}

type messageGroupModel struct {
	ID               string    `dynamodbav:"message_group_uuid"`
	OtherUserID      string    `dynamodbav:"user_uuid"`
	OtherDisplayName string    `dynamodbav:"user_display_name"`
	OtherHandle      string    `dynamodbav:"user_handle"`
	Content          string    `dynamodbav:"message"`
	LastSentAt       time.Time `dynamodbav:"sk"`
}

func (m *messageGroupModel) toDomain() domain.MessageGroup {
	return domain.MessageGroup{
		ID:          m.ID,
		UserID:      m.OtherUserID,
		DisplayName: m.OtherDisplayName,
		Handle:      m.OtherHandle,
		Content:     m.Content,
		LastMessage: m.Content,
		LastSentAt:  m.LastSentAt,
	}
}

func (m *messageModel) toDomain(groupID string) domain.Message {
	return domain.Message{
		ID:          m.ID,
		GroupID:     groupID,
		SenderID:    m.OtherUserID,
		DisplayName: m.OtherDisplayName,
		Handle:      m.OtherHandle,
		Content:     m.Content,
		SentAt:      m.LastSentAt,
	}
}

package domain

import "context"

type IMessageRepository interface {
	ListMessageGroups(ctx context.Context, userID string) (*[]MessageGroup, error)
	ListMessages(ctx context.Context, userID string) (*[]Message, error)
	CreateMessage(ctx context.Context, input *CreateMessage) (*Message, error)
	CreateMessageGroup(ctxs context.Context, input *CreateMessageGroup) (*MessageGroup, error)
}

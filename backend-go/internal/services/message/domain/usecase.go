package domain

import (
	"context"
)

type IMessageService interface {
	ListMessageGroups(ctx context.Context, handle string) (*[]MessageGroup, error)
	ListMessages(ctx context.Context, groupID string) (*[]Message, error)
}

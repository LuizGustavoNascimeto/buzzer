package usecase

import "time"

type CreateActivityInput struct {
	UserHandle          string
	Message             string
	ReplyToActivityUUID *int
	ExpiresAt           *time.Time
}

type UpdateActivityInput struct {
	ID      string
	Message string
}

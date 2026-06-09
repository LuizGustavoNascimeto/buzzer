package domain

import (
	"errors"
	"time"
)

type Activity struct {
	ID                  string
	UserID              string
	Message             string
	RepliesCount        int
	RepostsCount        int
	LikesCount          int
	ReplyToActivityUUID int
	ExpiresAt           *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// Activity representa uma atividade/post no feed

func (a *Activity) Validate() error {
	if a.Message == "" {
		return errors.New("message required")
	}
	if a.UserID == "" {
		return errors.New("UserHandle required")
	}
	return nil
}

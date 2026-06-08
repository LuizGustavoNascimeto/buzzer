package domain

import (
	"errors"
	"time"
)

type User struct {
	ID            string
	DisplayName   string
	Handle        string
	CognitoUserID string
	CreatedAt     time.Time
}

// Activity representa uma atividade/post no feed

func (a *User) Validate() error {
	if a.ID == "" {
		return errors.New("ID is required")
	}
	return nil
}

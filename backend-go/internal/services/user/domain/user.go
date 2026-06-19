package domain

import (
	"errors"
	"time"
)

type User struct {
	ID            string
	DisplayName   string
	Handle        string
	Email         string
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

type CreateMessageUsers struct {
	ID          string `gorm:"column:uuid"`
	DisplayName string `gorm:"column:display_name"`
	Handle      string `gorm:"column:handle"`
	Kind        string `gorm:"column:kind"`
}

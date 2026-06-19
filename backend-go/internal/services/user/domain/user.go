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

// MessageParticipant é uma projeção resumida de User, usada quando outro
// módulo (message) precisa apenas dos dados de identidade de um usuário
// participando de uma troca de mensagens — não é uma entidade própria.
type MessageParticipant struct {
	ID          string
	DisplayName string
	Handle      string
	Kind        string // "sender" | "receiver"
}

// Activity representa uma atividade/post no feed

func (a *User) Validate() error {
	if a.ID == "" {
		return errors.New("ID is required")
	}
	return nil
}

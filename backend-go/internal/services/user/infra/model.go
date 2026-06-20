package infra

import (
	"backend-go/internal/services/user/domain"
	gormutil "backend-go/pkg/gormutil/base"

	"github.com/google/uuid"
)

type UserModel struct {
	gormutil.Base
	DisplayName   string `gorm:"type:text;not null"`
	Handle        string `gorm:"type:text;not null"`
	Email         string `gorm:"type:text;not null"`
	CognitoUserID string `gorm:"type:text;not null"`
}

func (UserModel) TableName() string {
	return "user"
}

type messageParticipantModel struct {
	ID          string `gorm:"column:id"`
	DisplayName string `gorm:"column:display_name"`
	Handle      string `gorm:"column:handle"`
	Kind        string `gorm:"column:kind"`
}

func toModel(u *domain.User) *UserModel {
	var id uuid.UUID
	if u.ID != "" {
		id = uuid.MustParse(u.ID)
	}

	return &UserModel{
		Base: gormutil.Base{
			ID:        id,
			CreatedAt: u.CreatedAt,
		},
		DisplayName:   u.DisplayName,
		Handle:        u.Handle,
		CognitoUserID: u.CognitoUserID,
	}
}

func toDomain(m *UserModel) *domain.User {
	return &domain.User{
		ID:            m.ID.String(),
		CreatedAt:     m.CreatedAt,
		DisplayName:   m.DisplayName,
		Handle:        m.Handle,
		CognitoUserID: m.CognitoUserID,
	}
}

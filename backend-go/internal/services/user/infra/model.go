package infra

import (
	"backend-go/internal/services/user/domain"
	gormutil "backend-go/pkg/gormutil/base"

	"github.com/google/uuid"
)

type userModel struct {
	gormutil.Base
	DisplayName   string `gorm:"type:text;not null"`
	Handle        string `gorm:"type:text;not null"`
	CognitoUserID string `gorm:"type:text;not null"`
}

func (userModel) TableName() string {
	return "user"
}

func toModel(u *domain.User) *userModel {
	var id uuid.UUID
	if u.ID != "" {
		id = uuid.MustParse(u.ID)
	}

	return &userModel{
		Base: gormutil.Base{
			ID:        id,
			CreatedAt: u.CreatedAt,
		},
		DisplayName:   u.DisplayName,
		Handle:        u.Handle,
		CognitoUserID: u.CognitoUserID,
	}
}

func toDomain(m *userModel) *domain.User {
	return &domain.User{
		ID:            m.ID.String(),
		CreatedAt:     m.CreatedAt,
		DisplayName:   m.DisplayName,
		Handle:        m.Handle,
		CognitoUserID: m.CognitoUserID,
	}
}

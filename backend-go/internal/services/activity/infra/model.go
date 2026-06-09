package infra

import (
	"backend-go/internal/services/activity/domain"
	gormutil "backend-go/pkg/gormutil/base"
	"time"

	"github.com/google/uuid"
)

type ActivityModel struct {
	gormutil.Base
	UserID              string `gorm:"type:text"`
	Message             string `gorm:"type:text;not null"`
	RepliesCount        int    `gorm:"default:0"`
	RepostsCount        int    `gorm:"default:0"`
	LikesCount          int    `gorm:"default:0"`
	ReplyToActivityUUID int    `gorm:"type:int;"`
	ExpiresAt           *time.Time
}

func (ActivityModel) TableName() string {
	return "activity"
}

func toModel(a *domain.Activity) *ActivityModel {
	var id uuid.UUID
	if a.ID != "" {
		id = uuid.MustParse(a.ID)
	}

	return &ActivityModel{
		Base: gormutil.Base{
			ID:        id,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		},
		UserID:              a.UserID,
		Message:             a.Message,
		RepliesCount:        a.RepliesCount,
		RepostsCount:        a.RepostsCount,
		LikesCount:          a.LikesCount,
		ReplyToActivityUUID: a.ReplyToActivityUUID,
		ExpiresAt:           a.ExpiresAt,
	}
}

func toDomain(m *ActivityModel) *domain.Activity {
	return &domain.Activity{
		ID:                  m.ID.String(),
		UserID:              m.UserID,
		Message:             m.Message,
		RepliesCount:        m.RepliesCount,
		RepostsCount:        m.RepostsCount,
		LikesCount:          m.LikesCount,
		ReplyToActivityUUID: m.ReplyToActivityUUID,
		ExpiresAt:           m.ExpiresAt,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}

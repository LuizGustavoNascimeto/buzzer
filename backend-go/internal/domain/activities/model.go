package activities

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Activity representa uma atividade/post no feed
type Activity struct {
	UUID                string         `gorm:"type:uuid;primaryKey"   json:"uuid"`
	Handle              string         `gorm:"size:50;not null;index" json:"handle"` //vai virar um user no futuro
	Message             string         `gorm:"type:text;not null"     json:"message"`
	LikesCount          int            `gorm:"default:0"              json:"likes_count"`
	RepliesCount        int            `gorm:"default:0"              json:"replies_count"`
	RepostsCount        int            `gorm:"default:0"              json:"reposts_count"`
	ReplyToActivityUUID *string        `gorm:"type:uuid;index"        json:"reply_to_activity_uuid,omitempty"`
	Replies             []Activity     `gorm:"-"                      json:"replies"`
	ExpiresAt           *time.Time     `gorm:"index"                  json:"expires_at,omitempty"`
	CreatedAt           time.Time      `                              json:"created_at"`
	UpdatedAt           time.Time      `                              json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index"                  json:"deleted_at"`
}

func (a *Activity) Validate() error {
	if a.Message == "" {
		return errors.New("message required")
	}
	if a.Handle == "" {
		return errors.New("handle required")
	}
	return nil
}

package users

import (
	"time"

	"gorm.io/gorm"
)

type MessageGroup struct {
	UUID        string         `gorm:"type:uuid;primaryKey"   json:"uuid"`
	DisplayName string         `gorm:"size:100;not null"      json:"display_name"`
	Handle      string         `gorm:"size:50;not null;index" json:"handle"`
	CreatedAt   time.Time      `                              json:"created_at"`
	UpdatedAt   time.Time      `                              json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                  json:"deleted_at"`
}

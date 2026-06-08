package infra

// cmd/api/main.go
import (
	activityInfra "backend-go/internal/services/activity/infra"
	userInfra "backend-go/internal/services/user/infra"
	"fmt"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	if err := activityInfra.Migrate(db); err != nil {
		return fmt.Errorf("activity migration: %w", err)
	}
	if err := userInfra.Migrate(db); err != nil {
		return fmt.Errorf("user migration: %w", err)
	}
	return nil
}

package domain

import (
	"context"

	"github.com/google/uuid"
)

type IActivityRepository interface {
	Create(ctx context.Context, activity *Activity) error
	FindByID(ctx context.Context, id string) (*Activity, error)
	FindAll(ctx context.Context) ([]*Activity, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*Activity, error)
	Update(ctx context.Context, activity *Activity) error
	Delete(ctx context.Context, id string) error
}

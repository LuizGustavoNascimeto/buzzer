// internal/services/activity/infra/gorm_repository.go
package infra

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend-go/internal/services/activity/domain"
)

type ActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) Create(ctx context.Context, activity *domain.Activity) error {
	if activity == nil {
		return domain.ErrActivityRequired
	}
	if err := activity.Validate(); err != nil {
		return err
	}

	model := toModel(activity)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	activity.ID = model.ID.String()
	activity.CreatedAt = model.CreatedAt
	activity.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *ActivityRepository) FindByID(ctx context.Context, id string) (*domain.Activity, error) {
	var model activityModel

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrActivityNotFound
	}
	if err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (r *ActivityRepository) FindAll(ctx context.Context) ([]*domain.Activity, error) {
	var models []activityModel

	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	activities := make([]*domain.Activity, len(models))
	for i, m := range models {
		activities[i] = toDomain(&m)
	}
	return activities, nil
}

func (r *ActivityRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Activity, error) {
	var models []activityModel

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	activities := make([]*domain.Activity, len(models))
	for i, m := range models {
		activities[i] = toDomain(&m)
	}
	return activities, nil
}

func (r *ActivityRepository) Update(ctx context.Context, activity *domain.Activity) error {
	if activity == nil {
		return domain.ErrActivityRequired
	}
	if activity.ID == "" {
		return domain.ErrActivityIDRequired
	}
	if err := activity.Validate(); err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Model(&activityModel{}).
		Where("id = ?", activity.ID).
		Updates(toModel(activity))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrActivityNotFound
	}

	return nil
}

func (r *ActivityRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&activityModel{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrActivityNotFound
	}

	return nil
}

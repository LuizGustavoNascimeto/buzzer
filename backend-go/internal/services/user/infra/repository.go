// internal/services/user/infra/gorm_repository.go
package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"backend-go/internal/services/user/domain"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	if user == nil {
		return domain.ErrUserRequired
	}
	if err := user.Validate(); err != nil {
		return err
	}

	model := toModel(user)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	user.ID = model.ID.String()
	user.CreatedAt = model.CreatedAt
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var model UserModel

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
	var models []UserModel

	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(models))
	for i, m := range models {
		users[i] = toDomain(&m)
	}
	return users, nil
}

func (r *UserRepository) FindByHandle(ctx context.Context, handle string) (*domain.User, error) {
	var model UserModel

	err := r.db.WithContext(ctx).
		Where("handle = ?", handle).
		First(&model).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	if user == nil {
		return domain.ErrUserRequired
	}
	if user.ID == "" {
		return domain.ErrUserIDRequired
	}
	if err := user.Validate(); err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("id = ?", user.ID).
		Updates(toModel(user))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&UserModel{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) CreateMessageUser(ctx context.Context, senderHandler string, receiverHandler string) ([]*domain.CreateMessageUsers, error) {

	var result []*domain.CreateMessageUsers

	err := r.db.
		Table("public.users").
		Select(`
            users.uuid,
            users.display_name,
            users.handle,
            CASE users.handle = ?
                WHEN TRUE THEN 'sender'
                WHEN FALSE THEN 'receiver'
                ELSE 'other'
            END as kind
        `, senderHandler).
		Where(`
            users.handle = ? OR users.handle = ?
        `, senderHandler, receiverHandler).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

// SELECT
//   users.uuid,
//   users.display_name,
//   users.handle,
//   CASE users.cognito_user_id = %(cognito_user_id)s
//   WHEN TRUE THEN
//     'sender'
//   WHEN FALSE THEN
//     'recv'
//   ELSE
//     'other'
//   END as kind
// FROM public.users
// WHERE
//   users.cognito_user_id = %(cognito_user_id)s
//   OR
//   users.handle = %(user_receiver_handle)s

// internal/services/user/domain/errors.go
package domain

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserRequired   = errors.New("user is required")
	ErrUserIDRequired = errors.New("user id is required")
	ErrHandleRequired = errors.New("handle is required")
	ErrHandleTooLong  = errors.New("handle must be at most 255 characters")
)

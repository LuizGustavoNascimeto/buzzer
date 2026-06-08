// internal/services/activity/domain/errors.go
package domain

import "errors"

var (
	ErrActivityNotFound   = errors.New("activity not found")
	ErrActivityRequired   = errors.New("activity is required")
	ErrActivityIDRequired = errors.New("activity id is required")
	ErrHandleRequired     = errors.New("handle is required")
	ErrHandleTooLong      = errors.New("handle must be at most 255 characters")
)

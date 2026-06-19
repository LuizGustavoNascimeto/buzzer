package usecase

import "backend-go/internal/services/activity/domain"

type ActivityOutput struct {
	domain.Activity
	UserHandle      string
	UserDisplayName string
}

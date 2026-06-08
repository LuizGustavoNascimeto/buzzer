package handler

import (
	"backend-go/internal/services/user/domain"
	"time"
)

// ─── requests ────────────────────────────────────────────────────────────────

type CreateUserRequest struct {
	DisplayName   string `json:"display_name"    binding:"required,max=100"`
	Email         string `json:"email"           binding:"required,max=100"`
	Handle        string `json:"handle"          binding:"required,max=50"`
	CognitoUserID string `json:"cognito_user_id" binding:"required"`
}

type UpdateUserRequest struct {
	DisplayName string `json:"display_name" binding:"required,max=100"`
}

// ─── response (domain → cliente) ─────────────────────────────────────────────

type UserResponse struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Handle      string    `json:"handle"`
	CreatedAt   time.Time `json:"created_at"`
}

// ─── conversões ──────────────────────────────────────────────────────────────

func toUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:          u.ID,
		DisplayName: u.DisplayName,
		Handle:      u.Handle,
		CreatedAt:   u.CreatedAt,
		Email:       u.Email,
	}
}

func toUserListResponse(users []*domain.User) []UserResponse {
	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = toUserResponse(u)
	}
	return result
}

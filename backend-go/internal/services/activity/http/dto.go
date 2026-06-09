package handler

import (
	"backend-go/internal/services/activity/usecase"
	"time"
)

// ─── requests ────────────────────────────────────────────────────────────────

type CreateActivityRequest struct {
	UserHandle          string     `json:"user_handle"            binding:"required"`
	Message             string     `json:"message"                binding:"required,max=500"`
	ReplyToActivityUUID *int       `json:"reply_to_activity_uuid" binding:"omitempty"`
	ExpiresAt           *time.Time `json:"expires_at"             binding:"required"`
}

type UpdateActivityRequest struct {
	Message string `json:"message" binding:"required,max=500"`
}

// ─── response (usecase → cliente) ────────────────────────────────────────────

type ActivityResponse struct {
	ID                  string     `json:"id"`
	UserHandle          string     `json:"user_handle"`
	UserDisplayName     string     `json:"user_display_name"`
	Message             string     `json:"message"`
	RepliesCount        int        `json:"replies_count"`
	RepostsCount        int        `json:"reposts_count"`
	LikesCount          int        `json:"likes_count"`
	ReplyToActivityUUID *int       `json:"reply_to_activity_uuid,omitempty"`
	ExpiresAt           *time.Time `json:"expires_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// ─── conversões ──────────────────────────────────────────────────────────────

func toActivityResponse(a *usecase.ActivityResponse) ActivityResponse {
	return ActivityResponse{
		ID:                  a.ID,
		UserHandle:          a.UserHandle,
		UserDisplayName:     a.UserDisplayName,
		Message:             a.Message,
		RepliesCount:        a.RepliesCount,
		RepostsCount:        a.RepostsCount,
		LikesCount:          a.LikesCount,
		ReplyToActivityUUID: &a.ReplyToActivityUUID,
		ExpiresAt:           a.ExpiresAt,
		CreatedAt:           a.CreatedAt,
		UpdatedAt:           a.UpdatedAt,
	}
}

func toActivityListResponse(activities []*usecase.ActivityResponse) []ActivityResponse {
	result := make([]ActivityResponse, len(activities))
	for i, a := range activities {
		result[i] = toActivityResponse(a)
	}
	return result
}

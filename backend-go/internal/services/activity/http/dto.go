package handler

import (
	"backend-go/internal/services/activity/domain"
	"time"
)

// ─── requests ────────────────────────────────────────────────────────────────

type CreateActivityRequest struct {
	UserID              string     `json:"user_id"                binding:"required"`
	Message             string     `json:"message"                binding:"required,max=500"`
	ReplyToActivityUUID *int       `json:"reply_to_activity_uuid" binding:"omitempty"`
	ExpiresAt           *time.Time `json:"expires_at"             binding:"omitempty"`
}

type UpdateActivityRequest struct {
	Message string `json:"message" binding:"required,max=500"`
}

// ─── response (domain → cliente) ─────────────────────────────────────────────

type ActivityResponse struct {
	ID                  string     `json:"id"`
	UserID              string     `json:"user_id"`
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

func toActivityResponse(a *domain.Activity) ActivityResponse {
	return ActivityResponse{
		ID:                  a.ID,
		UserID:              a.UserID,
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

func toActivityListResponse(activities []*domain.Activity) []ActivityResponse {
	result := make([]ActivityResponse, len(activities))
	for i, a := range activities {
		result[i] = toActivityResponse(a)
	}
	return result
}

package usecase

import (
	"backend-go/internal/services/activity/domain"
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("buzzer-go/activities")

type IActivityService interface {
	Create(ctx context.Context, input CreateActivityInput) (*domain.Activity, error)
	FindByID(ctx context.Context, id string) (*domain.Activity, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Activity, error)
	FindAll(ctx context.Context) ([]*domain.Activity, error)
	Update(ctx context.Context, input UpdateActivityInput) (*domain.Activity, error)
	Delete(ctx context.Context, id string) error
}

// ─── inputs / outputs ────────────────────────────────────────────────────────

type CreateActivityInput struct {
	UserID              string
	Message             string
	ReplyToActivityUUID *int
	ExpiresAt           *time.Time
}

type UpdateActivityInput struct {
	ID      string
	Message string
}

// ─── usecase ─────────────────────────────────────────────────────────────────

type ActivityUsecase struct {
	repo domain.IActivityRepository
}

func NewActivityUsecase(repo domain.IActivityRepository) *ActivityUsecase {
	return &ActivityUsecase{repo: repo}
}

func (s *ActivityUsecase) Create(ctx context.Context, input CreateActivityInput) (*domain.Activity, error) {
	ctx, span := tracer.Start(ctx, "activities.create")
	defer span.End()

	span.SetAttributes(attribute.String("activity.user_id", input.UserID))

	activity := &domain.Activity{
		UserID:              input.UserID,
		Message:             input.Message,
		ReplyToActivityUUID: *input.ReplyToActivityUUID,
		ExpiresAt:           input.ExpiresAt,
	}

	if err := activity.Validate(); err != nil {
		recordSpanError(span, err, "err-activities-validate")
		return nil, err
	}

	if err := s.repo.Create(ctx, activity); err != nil {
		recordSpanError(span, err, "err-activities-create")
		return nil, err
	}

	span.SetAttributes(attribute.String("activity.id", activity.ID))
	return activity, nil
}

func (s *ActivityUsecase) FindByID(ctx context.Context, id string) (*domain.Activity, error) {
	ctx, span := tracer.Start(ctx, "activities.find_by_id")
	defer span.End()

	span.SetAttributes(attribute.String("activity.id", id))

	activity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		recordSpanError(span, err, "err-activities-find-by-id")
		return nil, err
	}

	return activity, nil
}

func (s *ActivityUsecase) FindAll(ctx context.Context) ([]*domain.Activity, error) {
	ctx, span := tracer.Start(ctx, "activities.list")
	defer span.End()

	activities, err := s.repo.FindAll(ctx)
	if err != nil {
		recordSpanError(span, err, "err-activities-find-all")
		return nil, err
	}

	span.SetAttributes(attribute.Int("activities.count", len(activities)))
	return activities, nil
}

func (s *ActivityUsecase) FindByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Activity, error) {
	ctx, span := tracer.Start(ctx, "activities.list_by_handle")
	defer span.End()

	span.SetAttributes(attribute.String("user.handle", userID.String()))

	activities, err := s.repo.FindByUser(ctx, userID)
	if err != nil {
		recordSpanError(span, err, "err-activities-find-by-handle")
		return nil, err
	}

	span.SetAttributes(attribute.Int("activities.count", len(activities)))
	return activities, nil
}

func (s *ActivityUsecase) Update(ctx context.Context, input UpdateActivityInput) (*domain.Activity, error) {
	ctx, span := tracer.Start(ctx, "activities.update")
	defer span.End()

	span.SetAttributes(attribute.String("activity.id", input.ID))

	activity, err := s.repo.FindByID(ctx, input.ID)
	if err != nil {
		recordSpanError(span, err, "err-activities-find-by-id")
		return nil, err
	}

	activity.Message = input.Message

	if err := activity.Validate(); err != nil {
		recordSpanError(span, err, "err-activities-validate")
		return nil, err
	}

	if err := s.repo.Update(ctx, activity); err != nil {
		recordSpanError(span, err, "err-activities-update")
		return nil, err
	}

	return activity, nil
}

func (s *ActivityUsecase) Delete(ctx context.Context, id string) error {
	ctx, span := tracer.Start(ctx, "activities.delete")
	defer span.End()

	span.SetAttributes(attribute.String("activity.id", id))

	if err := s.repo.Delete(ctx, id); err != nil {
		recordSpanError(span, err, "err-activities-delete")
		return err
	}

	return nil
}

func recordSpanError(span trace.Span, err error, slug string) {
	if err == nil {
		return
	}
	span.RecordError(err)
	span.SetAttributes(
		attribute.Bool("error", true),
		attribute.String("exception.slug", slug),
	)
	span.SetStatus(codes.Error, err.Error())
}

package usecase

import (
	"backend-go/internal/services/activity/domain"
	userService "backend-go/internal/services/user/usecase"

	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("buzzer-go/activities")

type IActivityRepository interface {
	Create(ctx context.Context, activity *domain.Activity) error
	FindByID(ctx context.Context, id string) (*domain.Activity, error)
	FindAll(ctx context.Context) ([]*domain.Activity, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Activity, error)
	Update(ctx context.Context, activity *domain.Activity) error
	Delete(ctx context.Context, id string) error
}
type IActivityService interface {
	Create(ctx context.Context, input CreateActivityInput) (*ActivityOutput, error)
	FindByID(ctx context.Context, id string) (*ActivityOutput, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*ActivityOutput, error)
	FindAll(ctx context.Context) ([]*ActivityOutput, error)
	Update(ctx context.Context, input UpdateActivityInput) (*ActivityOutput, error)
	Delete(ctx context.Context, id string) error
}

type ActivityUsecase struct {
	repo     IActivityRepository
	userRepo userService.IUserRepository
}

func NewActivityUsecase(repo IActivityRepository, userRepo userService.IUserRepository) *ActivityUsecase {
	return &ActivityUsecase{repo: repo, userRepo: userRepo}
}

func (s *ActivityUsecase) enrich(ctx context.Context, activity *domain.Activity) (*ActivityOutput, error) {
	user, err := s.userRepo.FindByID(ctx, activity.UserID)
	if err != nil {
		return nil, err
	}
	return &ActivityOutput{
		Activity:        *activity,
		UserHandle:      user.Handle,
		UserDisplayName: user.DisplayName,
	}, nil
}

func (s *ActivityUsecase) enrichMany(ctx context.Context, activities []*domain.Activity) ([]*ActivityOutput, error) {
	result := make([]*ActivityOutput, len(activities))
	for i, a := range activities {
		enriched, err := s.enrich(ctx, a)
		if err != nil {
			return nil, err
		}
		result[i] = enriched
	}
	return result, nil
}

func (s *ActivityUsecase) Create(ctx context.Context, input CreateActivityInput) (*ActivityOutput, error) {
	ctx, span := tracer.Start(ctx, "activities.create")
	defer span.End()

	span.SetAttributes(attribute.String("activity.user_handle", input.UserHandle))

	user, err := s.userRepo.FindByHandle(ctx, input.UserHandle)
	if err != nil {
		return nil, err
	}

	reply := 0
	if input.ReplyToActivityUUID != nil {
		reply = *input.ReplyToActivityUUID
	}

	activity := &domain.Activity{
		UserID:              user.ID,
		Message:             input.Message,
		ReplyToActivityUUID: reply,
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

	enriched := &ActivityOutput{
		Activity:        *activity,
		UserHandle:      user.Handle,
		UserDisplayName: user.DisplayName,
	}
	return enriched, nil
}

func (s *ActivityUsecase) FindByID(ctx context.Context, id string) (*ActivityOutput, error) {
	ctx, span := tracer.Start(ctx, "activities.find_by_id")
	defer span.End()

	span.SetAttributes(attribute.String("activity.id", id))

	activity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		recordSpanError(span, err, "err-activities-find-by-id")
		return nil, err
	}

	enriched, err := s.enrich(ctx, activity)
	if err != nil {
		recordSpanError(span, err, "err-activities-enrich")
		return nil, err
	}

	return enriched, nil
}

func (s *ActivityUsecase) FindAll(ctx context.Context) ([]*ActivityOutput, error) {
	ctx, span := tracer.Start(ctx, "activities.list")
	defer span.End()

	activities, err := s.repo.FindAll(ctx)
	if err != nil {
		recordSpanError(span, err, "err-activities-find-all")
		return nil, err
	}

	result, err := s.enrichMany(ctx, activities)
	if err != nil {
		recordSpanError(span, err, "err-activities-enrich")
		return nil, err
	}

	span.SetAttributes(attribute.Int("activities.count", len(result)))
	return result, nil
}

func (s *ActivityUsecase) FindByUser(ctx context.Context, userID uuid.UUID) ([]*ActivityOutput, error) {
	ctx, span := tracer.Start(ctx, "activities.list_by_handle")
	defer span.End()

	span.SetAttributes(attribute.String("user.handle", userID.String()))

	activities, err := s.repo.FindByUser(ctx, userID)
	if err != nil {
		recordSpanError(span, err, "err-activities-find-by-handle")
		return nil, err
	}

	result, err := s.enrichMany(ctx, activities)
	if err != nil {
		recordSpanError(span, err, "err-activities-enrich")
		return nil, err
	}

	span.SetAttributes(attribute.Int("activities.count", len(result)))
	return result, nil
}

func (s *ActivityUsecase) Update(ctx context.Context, input UpdateActivityInput) (*ActivityOutput, error) {
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

	enriched, err := s.enrich(ctx, activity)
	if err != nil {
		recordSpanError(span, err, "err-activities-enrich")
		return nil, err
	}

	return enriched, nil
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

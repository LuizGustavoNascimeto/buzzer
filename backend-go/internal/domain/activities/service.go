package activities

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("backend-go/activities")

type ActivitiesService struct {
	repo ActivityRepository
}

func NewActivitiesService(repo ActivityRepository) *ActivitiesService {
	return &ActivitiesService{repo: repo}
}

func (s *ActivitiesService) CreateActivity(ctx context.Context, activity *Activity) error {
	ctx, span := tracer.Start(ctx, "activities.create")
	defer span.End()

	if activity != nil {
		attributes := []attribute.KeyValue{
			attribute.String("activity.handle", activity.Handle),
		}
		if activity.ReplyToActivityUUID != nil {
			attributes = append(attributes, attribute.String("activity.reply_to_uuid", *activity.ReplyToActivityUUID))
		}
		span.SetAttributes(attributes...)
	}

	if err := activity.Validate(); err != nil {
		recordSpanError(span, err, "err-activities-validate")
		return err
	}
	if err := s.repo.Create(ctx, activity); err != nil {
		recordSpanError(span, err, "err-activities-create")
		return err
	}
	span.SetAttributes(attribute.String("activity.uuid", activity.UUID))
	return nil
}

func (s *ActivitiesService) FindActivityByID(ctx context.Context, id string) (*Activity, error) {
	ctx, span := tracer.Start(ctx, "activities.find_by_id")
	defer span.End()
	span.SetAttributes(attribute.String("activity.uuid", id))

	activity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		recordSpanError(span, err, "err-activities-find-by-id")
		return nil, err
	}
	return activity, nil
}

func (s *ActivitiesService) FindActivities(ctx context.Context) ([]*Activity, error) {
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

func (s *ActivitiesService) FindActivitiesByHandle(ctx context.Context, handle string) ([]*Activity, error) {
	ctx, span := tracer.Start(ctx, "activities.list_by_handle")
	defer span.End()
	span.SetAttributes(attribute.String("user.handle", handle))

	activities, err := s.repo.FindByHandle(ctx, handle)
	if err != nil {
		recordSpanError(span, err, "err-activities-find-by-handle")
		return nil, err
	}
	span.SetAttributes(attribute.Int("activities.count", len(activities)))
	return activities, nil
}

func (s *ActivitiesService) UpdateActivity(ctx context.Context, activity *Activity) error {
	ctx, span := tracer.Start(ctx, "activities.update")
	defer span.End()

	if activity != nil {
		span.SetAttributes(
			attribute.String("activity.uuid", activity.UUID),
			attribute.String("activity.handle", activity.Handle),
		)
	}

	if err := activity.Validate(); err != nil {
		recordSpanError(span, err, "err-activities-validate")
		return err
	}
	if err := s.repo.Update(ctx, activity); err != nil {
		recordSpanError(span, err, "err-activities-update")
		return err
	}
	return nil
}

func (s *ActivitiesService) DeleteActivity(ctx context.Context, id string) error {
	ctx, span := tracer.Start(ctx, "activities.delete")
	defer span.End()
	span.SetAttributes(attribute.String("activity.uuid", id))

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

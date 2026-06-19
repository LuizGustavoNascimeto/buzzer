package usecase

import (
	"backend-go/internal/services/user/domain"
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("buzzer-go/users")

type IUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindAll(ctx context.Context) ([]domain.User, error)
	FindByHandle(ctx context.Context, handle string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
	CreateMessageUser(ctx context.Context, senderHandle string, receiverHandle string) ([]domain.MessageParticipant, error)
}

type IUserService interface {
	Create(ctx context.Context, input CreateUserInput) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByHandle(ctx context.Context, handle string) (*domain.User, error)
	FindAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, input UpdateUserInput) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}
type UserUsecase struct {
	repo IUserRepository
}

func NewUserUsecase(repo IUserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (s *UserUsecase) Create(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "users.create")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.handle", input.Handle),
		attribute.String("user.cognito_id", input.CognitoUserID),
	)

	user := &domain.User{
		DisplayName:   input.DisplayName,
		Handle:        input.Handle,
		CognitoUserID: input.CognitoUserID,
	}

	if err := user.Validate(); err != nil {
		recordSpanError(span, err, "err-users-validate")
		return nil, err
	}

	if err := s.repo.Create(ctx, user); err != nil {
		recordSpanError(span, err, "err-users-create")
		return nil, err
	}

	span.SetAttributes(attribute.String("user.id", user.ID))
	return user, nil
}

func (s *UserUsecase) FindByID(ctx context.Context, id string) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "users.find_by_id")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", id))

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		recordSpanError(span, err, "err-users-find-by-id")
		return nil, err
	}

	return user, nil
}

func (s *UserUsecase) FindByHandle(ctx context.Context, handle string) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "users.find_by_handle")
	defer span.End()

	span.SetAttributes(attribute.String("user.handle", handle))

	user, err := s.repo.FindByHandle(ctx, handle)
	if err != nil {
		recordSpanError(span, err, "err-users-find-by-handle")
		return nil, err
	}

	return user, nil
}

func (s *UserUsecase) FindAll(ctx context.Context) ([]domain.User, error) {
	ctx, span := tracer.Start(ctx, "users.list")
	defer span.End()

	users, err := s.repo.FindAll(ctx)
	if err != nil {
		recordSpanError(span, err, "err-users-find-all")
		return nil, err
	}

	span.SetAttributes(attribute.Int("users.count", len(users)))
	return users, nil
}

func (s *UserUsecase) Update(ctx context.Context, input UpdateUserInput) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "users.update")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", input.ID))

	user, err := s.repo.FindByID(ctx, input.ID)
	if err != nil {
		recordSpanError(span, err, "err-users-find-by-id")
		return nil, err
	}

	user.DisplayName = input.DisplayName

	if err := user.Validate(); err != nil {
		recordSpanError(span, err, "err-users-validate")
		return nil, err
	}

	if err := s.repo.Update(ctx, user); err != nil {
		recordSpanError(span, err, "err-users-update")
		return nil, err
	}

	return user, nil
}

func (s *UserUsecase) Delete(ctx context.Context, id string) error {
	ctx, span := tracer.Start(ctx, "users.delete")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", id))

	if err := s.repo.Delete(ctx, id); err != nil {
		recordSpanError(span, err, "err-users-delete")
		return err
	}

	return nil
}

// ─── helpers ─────────────────────────────────────────────────────────────────

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

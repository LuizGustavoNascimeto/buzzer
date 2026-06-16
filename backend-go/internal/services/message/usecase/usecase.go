package usecase

import (
	"backend-go/internal/services/message/domain"
	userDomain "backend-go/internal/services/user/domain"

	"context"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("buzzer-go/messages")

type MessageUsecase struct {
	repo     domain.IMessageRepository
	userRepo userDomain.IUserRepository
}

func NewMessageUsecase(repo domain.IMessageRepository, userRepo userDomain.IUserRepository) *MessageUsecase {
	return &MessageUsecase{repo: repo, userRepo: userRepo}
}

func (m MessageUsecase) ListMessageGroups(ctx context.Context, handle string) (*[]domain.MessageGroup, error) {
	ctx, span := tracer.Start(ctx, "groups.list")
	defer span.End()
	user, err := m.userRepo.FindByHandle(ctx, handle)
	if err != nil {
		return nil, err
	}
	groups, err := m.repo.ListMessageGroups(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (m MessageUsecase) ListMessages(ctx context.Context, groupID string) (*[]domain.Message, error) {
	ctx, span := tracer.Start(ctx, "messages.list")
	defer span.End()
	groups, err := m.repo.ListMessages(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return groups, nil

}

// func (s *ActivityUsecase) FindAll(ctx context.Context) ([]*ActivityResponse, error) {
// 	ctx, span := tracer.Start(ctx, "activities.list")
// 	defer span.End()

// 	activities, err := s.repo.FindAll(ctx)
// 	if err != nil {
// 		recordSpanError(span, err, "err-activities-find-all")
// 		return nil, err
// 	}

// 	result, err := s.enrichMany(ctx, activities)
// 	if err != nil {
// 		recordSpanError(span, err, "err-activities-enrich")
// 		return nil, err
// 	}

// 	span.SetAttributes(attribute.Int("activities.count", len(result)))
// 	return result, nil
// }

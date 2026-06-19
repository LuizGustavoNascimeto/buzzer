package usecase

import (
	"context"

	"backend-go/internal/services/message/domain"
	userDomain "backend-go/internal/services/user/domain"
	userService "backend-go/internal/services/user/usecase"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("buzzer-go/messages")

// Contrato que o usecase exige do repositório. Só tipos de domain ou
// inputs definidos aqui — nunca um tipo de infra (dynamodbav) ou HTTP (json).
type IMessageRepository interface {
	ListMessageGroups(ctx context.Context, userID string) ([]domain.MessageGroup, error)
	ListMessages(ctx context.Context, groupID string) ([]domain.Message, error)
	CreateMessage(ctx context.Context, msg *domain.Message) (*domain.Message, error)
	CreateMessageGroupsInBatch(ctx context.Context, my, other *CreateMessageGroupInput, msg *CreateMessageInput) (*domain.Message, error)
}

// Contrato que o handler exige do usecase (MessageUsecase implementa).
type IMessageService interface {
	ListMessageGroups(ctx context.Context, handle string) ([]domain.MessageGroup, error)
	ListMessages(ctx context.Context, groupID string) ([]domain.Message, error)
	CreateMessage(ctx context.Context, input *CreateMessageInput) error
}

type MessageUsecase struct {
	repo     IMessageRepository
	userRepo userService.IUserRepository
}

func NewMessageUsecase(repo IMessageRepository, userRepo userService.IUserRepository) *MessageUsecase {
	return &MessageUsecase{repo: repo, userRepo: userRepo}
}

func (m MessageUsecase) ListMessageGroups(ctx context.Context, handle string) ([]domain.MessageGroup, error) {
	ctx, span := tracer.Start(ctx, "groups.list")
	defer span.End()

	user, err := m.userRepo.FindByHandle(ctx, handle)
	if err != nil {
		return nil, err
	}

	return m.repo.ListMessageGroups(ctx, user.ID)
}

func (m MessageUsecase) ListMessages(ctx context.Context, groupID string) ([]domain.Message, error) {
	ctx, span := tracer.Start(ctx, "messages.list")
	defer span.End()

	return m.repo.ListMessages(ctx, groupID)
}

func (m MessageUsecase) CreateMessage(ctx context.Context, input *CreateMessageInput) error {
	ctx, span := tracer.Start(ctx, "messages.create")
	defer span.End()

	receiverHandle := ""
	if input.ReceiverHandle != nil {
		receiverHandle = *input.ReceiverHandle
	}

	users, err := m.userRepo.CreateMessageUser(ctx, input.SenderHandle, receiverHandle)
	if err != nil {
		return err
	}

	var myUser, otherUser userDomain.MessageParticipant
	for i := range users {
		switch users[i].Kind {
		case "sender":
			myUser = users[i]
		case "receiver":
			otherUser = users[i]
		}
	}

	if input.GroupID == nil {
		groupSender := &CreateMessageGroupInput{
			UserID:           myUser.ID,
			Content:          input.Content,
			OtherDisplayName: otherUser.DisplayName,
			OtherHandle:      otherUser.Handle,
			OtherUserID:      otherUser.ID,
		}
		groupReceiver := &CreateMessageGroupInput{
			UserID:           otherUser.ID,
			Content:          input.Content,
			OtherDisplayName: myUser.DisplayName,
			OtherHandle:      myUser.Handle,
			OtherUserID:      myUser.ID,
		}

		_, err := m.repo.CreateMessageGroupsInBatch(ctx, groupSender, groupReceiver, input)
		return err
	}

	msg := &domain.Message{
		GroupID:     *input.GroupID,
		SenderID:    myUser.ID,
		DisplayName: myUser.DisplayName,
		Handle:      myUser.Handle,
		Content:     input.Content,
	}

	_, err = m.repo.CreateMessage(ctx, msg)
	return err
}

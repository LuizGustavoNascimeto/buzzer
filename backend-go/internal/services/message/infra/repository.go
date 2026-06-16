package infra

import (
	"backend-go/internal/services/message/domain"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type MessageRepository struct {
	ddb       *dynamodb.Client
	tableName string
}

func NewMessageRepository(ddb *dynamodb.Client) *MessageRepository {
	return &MessageRepository{ddb: ddb, tableName: "buzzer-messages"}
}

func (r *MessageRepository) ListMessageGroups(ctx context.Context, userID string) (*[]domain.MessageGroup, error) {
	year := strconv.Itoa(time.Now().Year())
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :year)"),
		ScanIndexForward:       aws.Bool(false),
		Limit:                  aws.Int32(20),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":year": &types.AttributeValueMemberS{Value: year},
			":pk":   &types.AttributeValueMemberS{Value: "GRP#" + userID},
		},
	}

	res, err := r.ddb.Query(ctx, queryInput)
	if err != nil {
		return nil, fmt.Errorf("failed to query DynamoDB: %w", err)
	}
	var groups []domain.MessageGroup
	err = attributevalue.UnmarshalListOfMaps(res.Items, &groups)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return &groups, nil
}

func (r *MessageRepository) ListMessages(ctx context.Context, messageID string) (*[]domain.Message, error) {
	year := strconv.Itoa(time.Now().Year())
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :year)"),
		ScanIndexForward:       aws.Bool(false),
		Limit:                  aws.Int32(20),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":year": &types.AttributeValueMemberS{Value: year},
			":pk":   &types.AttributeValueMemberS{Value: "MSG#" + messageID},
		},
	}

	res, err := r.ddb.Query(ctx, queryInput)
	if err != nil {
		return nil, fmt.Errorf("failed to query DynamoDB: %w", err)
	}
	var messages []domain.Message
	err = attributevalue.UnmarshalListOfMaps(res.Items, &messages)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return &messages, nil
}

func (r *MessageRepository) CreateMessage(ctx context.Context, input *domain.CreateMessage) (*domain.Message, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	messageUUID := uuid.New().String()

	item := map[string]types.AttributeValue{
		"pk":                &types.AttributeValueMemberS{Value: "MSG#" + input.MessageGroupID},
		"sk":                &types.AttributeValueMemberS{Value: now},
		"message":           &types.AttributeValueMemberS{Value: input.Message},
		"message_uuid":      &types.AttributeValueMemberS{Value: messageUUID},
		"user_uuid":         &types.AttributeValueMemberS{Value: input.UserID},
		"user_display_name": &types.AttributeValueMemberS{Value: input.DisplayName},
		"user_handle":       &types.AttributeValueMemberS{Value: input.Handle},
	}

	_, err := r.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("createMessage: failed to put item: %w", err)
	}

	createdAt, _ := time.Parse(time.RFC3339, now)

	return &domain.Message{
		ID:          messageUUID,
		DisplayName: input.DisplayName,
		Handle:      input.Handle,
		Message:     input.Message,
		LastSentAt:  createdAt,
	}, nil
}

func (r *MessageRepository) CreateMessageGroup(ctx context.Context, input *domain.CreateMessageGroup) (*domain.MessageGroup, error) {
	now := time.Now().Format(time.RFC3339)
	messageGroupUUID := uuid.New().String()

	item := map[string]types.AttributeValue{
		"pk":                 &types.AttributeValueMemberS{Value: "GRP#" + input.UserID},
		"sk":                 &types.AttributeValueMemberS{Value: now},
		"message_group_uuid": &types.AttributeValueMemberS{Value: messageGroupUUID},
		"message":            &types.AttributeValueMemberS{Value: input.Message},
		"user_uuid":          &types.AttributeValueMemberS{Value: input.UserID},
		"user_display_name":  &types.AttributeValueMemberS{Value: input.DisplayName},
		"user_handle":        &types.AttributeValueMemberS{Value: input.Handle},
	}

	_, err := r.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("createMessageGroup: failed to put item: %w", err)
	}

	createdAt, _ := time.Parse(time.RFC3339, now)

	return &domain.MessageGroup{
		ID:          messageGroupUUID,
		DisplayName: input.DisplayName,
		Handle:      input.Handle,
		Message:     input.Message,
		LastSentAt:  createdAt,
	}, nil
}

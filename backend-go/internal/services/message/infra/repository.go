package infra

import (
	"backend-go/internal/services/message/domain"
	"backend-go/internal/services/message/usecase"
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

func (r *MessageRepository) ListMessageGroups(ctx context.Context, userID string) ([]domain.MessageGroup, error) {
	year := strconv.Itoa(time.Now().Year())
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :year)"),
		ScanIndexForward:       aws.Bool(false),
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

	return groups, nil
}

func (r *MessageRepository) ListMessages(ctx context.Context, groupID string) ([]domain.Message, error) {
	year := strconv.Itoa(time.Now().Year())
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :year)"),
		ScanIndexForward:       aws.Bool(false),
		Limit:                  aws.Int32(20),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":year": &types.AttributeValueMemberS{Value: year},
			":pk":   &types.AttributeValueMemberS{Value: "MSG#" + groupID},
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

	return messages, nil
}

func (r *MessageRepository) CreateMessage(ctx context.Context, msg *domain.Message) (*domain.Message, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	messageUUID := uuid.New().String()

	item := map[string]types.AttributeValue{
		"pk":                &types.AttributeValueMemberS{Value: "MSG#" + msg.GroupID},
		"sk":                &types.AttributeValueMemberS{Value: now},
		"message":           &types.AttributeValueMemberS{Value: msg.Content},
		"message_uuid":      &types.AttributeValueMemberS{Value: messageUUID},
		"user_uuid":         &types.AttributeValueMemberS{Value: msg.SenderID},
		"user_display_name": &types.AttributeValueMemberS{Value: msg.DisplayName},
		"user_handle":       &types.AttributeValueMemberS{Value: msg.Handle},
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
		DisplayName: msg.DisplayName,
		Handle:      msg.Handle,
		Content:     msg.Content,
		SentAt:      createdAt,
	}, nil
}

func (r *MessageRepository) CreateMessageGroupsInBatch(ctx context.Context, my, other *usecase.CreateMessageGroupInput, message *usecase.CreateMessageInput) (*domain.Message, error) {
	now := time.Now()
	nowStr := now.Format(time.RFC3339)
	messageGroupUUID := uuid.New().String()
	messageUUID := uuid.New().String()

	myItem := map[string]types.AttributeValue{
		"pk":                 &types.AttributeValueMemberS{Value: "GRP#" + my.UserID},
		"sk":                 &types.AttributeValueMemberS{Value: nowStr},
		"message_group_uuid": &types.AttributeValueMemberS{Value: messageGroupUUID},
		"message":            &types.AttributeValueMemberS{Value: my.Content},
		"user_uuid":          &types.AttributeValueMemberS{Value: my.OtherUserID},
		"user_display_name":  &types.AttributeValueMemberS{Value: my.OtherDisplayName},
		"user_handle":        &types.AttributeValueMemberS{Value: my.OtherHandle},
	}

	otherItem := map[string]types.AttributeValue{
		"pk":                 &types.AttributeValueMemberS{Value: "GRP#" + other.UserID},
		"sk":                 &types.AttributeValueMemberS{Value: nowStr},
		"message_group_uuid": &types.AttributeValueMemberS{Value: messageGroupUUID},
		"message":            &types.AttributeValueMemberS{Value: other.Content},
		"user_uuid":          &types.AttributeValueMemberS{Value: other.OtherUserID},
		"user_display_name":  &types.AttributeValueMemberS{Value: other.OtherDisplayName},
		"user_handle":        &types.AttributeValueMemberS{Value: other.OtherHandle},
	}

	messageItem := map[string]types.AttributeValue{
		"pk":                 &types.AttributeValueMemberS{Value: "MSG#" + messageGroupUUID},
		"sk":                 &types.AttributeValueMemberS{Value: nowStr},
		"message_uuid":       &types.AttributeValueMemberS{Value: messageUUID},
		"message_group_uuid": &types.AttributeValueMemberS{Value: messageGroupUUID},
		"message":            &types.AttributeValueMemberS{Value: message.Content},
		"sender_handle":      &types.AttributeValueMemberS{Value: message.SenderHandle},
	}
	if message.ReceiverHandle != nil {
		messageItem["receiver_handle"] = &types.AttributeValueMemberS{Value: *message.ReceiverHandle}
	}

	_, err := r.ddb.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{Put: &types.Put{TableName: aws.String(r.tableName), Item: myItem}},
			{Put: &types.Put{TableName: aws.String(r.tableName), Item: otherItem}},
			{Put: &types.Put{TableName: aws.String(r.tableName), Item: messageItem}},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("createMessageGroupsInBatch: failed to transact write items: %w", err)
	}

	return &domain.Message{
		ID:          messageGroupUUID,
		GroupID:     messageGroupUUID,
		SenderID:    other.UserID,
		Content:     message.Content,
		SentAt:      now,
		DisplayName: my.OtherDisplayName,
		Handle:      my.OtherHandle,
	}, nil
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	tableName        = "buzzer-messages"
	messageGroupUUID = "5ae290ed-55d1-47a0-bc6d-fe2bc2700399"
)

func main() {
	ctx := context.Background()

	isProd := len(os.Args) == 2 && contains(os.Args[1], "prod")

	client := buildClient(ctx, isProd)

	currentYear := fmt.Sprintf("%d", time.Now().Year())

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		ScanIndexForward:       aws.Bool(false),
		Limit:                  aws.Int32(20),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :year)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":year": &types.AttributeValueMemberS{Value: currentYear},
			":pk":   &types.AttributeValueMemberS{Value: "MSG#" + messageGroupUUID},
		},
	}

	response, err := client.Query(ctx, queryInput)
	if err != nil {
		log.Fatalf("failed to query DynamoDB: %v", err)
	}

	// Print full response as JSON
	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}
	fmt.Println(string(responseJSON))

	// Print consumed capacity
	if response.ConsumedCapacity != nil {
		capacityJSON, err := json.MarshalIndent(response.ConsumedCapacity, "", "  ")
		if err != nil {
			log.Fatalf("failed to marshal consumed capacity: %v", err)
		}
		fmt.Println(string(capacityJSON))
	}

	// Reverse items (query returned in reverse, we flip back)
	items := response.Items
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}

	for _, item := range items {
		senderHandle := attrString(item, "user_handle")
		message := attrString(item, "message")
		timestamp := attrString(item, "sk")

		dt, err := time.Parse("2006-01-02T15:04:05.999999-07:00", timestamp)
		if err != nil {
			// try alternative format
			dt, err = time.Parse("2006-01-02T15:04:05.999999+00:00", timestamp)
			if err != nil {
				dt = time.Time{}
			}
		}

		formattedTime := dt.Format("2006-01-02 03:04 PM")

		truncatedMessage := message
		if len(message) > 40 {
			truncatedMessage = message[:40] + "..."
		}

		fmt.Printf("%-12s%-22s%s\n", senderHandle, formattedTime, truncatedMessage)
	}
}

func buildClient(ctx context.Context, isProd bool) *dynamodb.Client {
	if isProd {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("failed to load AWS config: %v", err)
		}
		return dynamodb.NewFromConfig(cfg)
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8000")
	})
}

func attrString(item map[string]types.AttributeValue, key string) string {
	if v, ok := item[key]; ok {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			return s.Value
		}
	}
	return ""
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsRune(s, substr))
}

func containsRune(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

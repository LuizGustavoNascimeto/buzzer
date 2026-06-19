package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jackc/pgx/v5"
)

const (
	tableName    = "buzzer-messages"
	targetHandle = "andrewbrown"
)

func main() {
	ctx := context.Background()

	isProd := len(os.Args) == 2 && strings.Contains(os.Args[1], "prod")

	// Busca UUID do usuário no PostgreSQL
	myUserUUID, err := getMyUserUUID(ctx)
	if err != nil {
		log.Fatalf("failed to get user uuid: %v", err)
	}
	fmt.Printf("my-uuid: %s\n", myUserUUID)

	// Monta cliente DynamoDB
	client := buildDynamoClient(ctx, isProd)

	currentYear := fmt.Sprintf("%d", time.Now().Year())

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :year)"),
		ScanIndexForward:       aws.Bool(false),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":year": &types.AttributeValueMemberS{Value: currentYear},
			":pk":   &types.AttributeValueMemberS{Value: "GRP#" + myUserUUID},
		},
	}

	response, err := client.Query(ctx, queryInput)
	if err != nil {
		log.Fatalf("failed to query DynamoDB: %v", err)
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}
	fmt.Println(string(responseJSON))
}

// getMyUserUUID busca o UUID do usuário no PostgreSQL via variável de ambiente CONNECTION_URL
func getMyUserUUID(ctx context.Context) (string, error) {
	connURL := os.Getenv("POSTGRES_URL")
	if connURL == "" {
		// fallback local
		connURL = "postgresql://postgres:password@localhost:5432/buzzer"
	}

	conn, err := pgx.Connect(ctx, connURL)
	if err != nil {
		return "", fmt.Errorf("unable to connect to postgres: %w", err)
	}
	defer conn.Close(ctx)

	sql := `
		SELECT "user".id
		FROM "user"
		WHERE "user".handle = $1
	`

	var uuid string
	err = conn.QueryRow(ctx, sql, targetHandle).Scan(&uuid)
	if err != nil {
		return "", fmt.Errorf("query failed: %w", err)
	}

	return uuid, nil
}

func buildDynamoClient(ctx context.Context, isProd bool) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	if isProd {
		return dynamodb.NewFromConfig(cfg)
	}

	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8000")
	})
}

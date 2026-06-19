package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Estrutura mínima do evento do DynamoDB Stream
type Event struct {
	Records []Record `json:"Records"`
}

type Record struct {
	EventName string `json:"eventName"`
	DynamoDB  struct {
		Keys struct {
			PK struct {
				S string `json:"S"`
			} `json:"pk"`
			SK struct {
				S string `json:"S"`
			} `json:"sk"`
		} `json:"Keys"`

		NewImage map[string]struct {
			S string `json:"S"`
		} `json:"NewImage"`
	} `json:"dynamodb"`
}

var (
	client *dynamodb.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client = dynamodb.NewFromConfig(cfg)
}

func handler(ctx context.Context, event Event) error {
	log.Printf("event: %+v\n", event)

	if len(event.Records) == 0 {
		return nil
	}

	rec := event.Records[0]

	// 1. Ignora REMOVE
	if rec.EventName == "REMOVE" {
		log.Println("skip REMOVE event")
		return nil
	}

	pk := rec.DynamoDB.Keys.PK.S
	sk := rec.DynamoDB.Keys.SK.S

	// 2. Filtra MSG#
	if !strings.HasPrefix(pk, "MSG#") {
		return nil
	}

	groupUUID := strings.Replace(pk, "MSG#", "", 1)

	message := rec.DynamoDB.NewImage["message"].S

	log.Printf("GRUP ===> %s %s\n", groupUUID, message)

	tableName := "cruddur-messages"
	indexName := "message-group-sk-index"

	// 3. Query no GSI
	out, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String(indexName),
		KeyConditionExpression: aws.String("message_group_uuid = :g"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":g": &types.AttributeValueMemberS{Value: groupUUID},
		},
	})
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}

	log.Printf("RESP ===> %d items\n", len(out.Items))

	// 4. Recria itens
	for _, item := range out.Items {

		pkAttr := item["pk"].(*types.AttributeValueMemberS).Value
		skAttr := item["sk"].(*types.AttributeValueMemberS).Value

		// DELETE
		_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: pkAttr},
				"sk": &types.AttributeValueMemberS{Value: skAttr},
			},
		})
		if err != nil {
			log.Printf("DELETE error: %v", err)
			continue
		}

		log.Printf("DELETE ===> %s %s\n", pkAttr, skAttr)

		// RECREATE
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"pk":                 &types.AttributeValueMemberS{Value: pkAttr},
				"sk":                 &types.AttributeValueMemberS{Value: sk},
				"message_group_uuid": item["message_group_uuid"],
				"message":            &types.AttributeValueMemberS{Value: message},
				"user_display_name":  item["user_display_name"],
				"user_handle":        item["user_handle"],
				"user_uuid":          item["user_uuid"],
			},
		})
		if err != nil {
			log.Printf("PUT error: %v", err)
			continue
		}

		log.Printf("CREATE ===> %s\n", pkAttr)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(os.Getenv("AWS_DEFAULT_REGION")),
	)
	if err != nil {
		log.Fatal("erro ao carregar config:", err)
	}

	// Local por padrão
	isProd := false

	if len(os.Args) == 2 {
		if strings.Contains(strings.ToLower(os.Args[1]), "prod") {
			isProd = true
		}
	}

	var client *dynamodb.Client

	if isProd {
		fmt.Println("Usando DynamoDB AWS (produção)")
		client = dynamodb.NewFromConfig(cfg)
	} else {
		fmt.Println("Usando DynamoDB Local")
		client = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String("http://localhost:8000")
		})
	}

	res, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String("buzzer-messages"),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("message_group_uuid"),
				AttributeType: "S",
			},
			{
				AttributeName: aws.String("pk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("sk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("pk"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("sk"),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("message-group-sk-index"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("message_group_uuid"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("sk"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
		},

		BillingMode: types.BillingModeProvisioned,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}

package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"log"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8000")
	})

	result, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("buzzer-messages"),
	})
	if err != nil {
		log.Fatal("erro ao fazer scan:", err)
	}

	for _, item := range result.Items {
		var m map[string]interface{}
		if err := attributevalue.UnmarshalMap(item, &m); err != nil {
			log.Printf("erro ao deserializar item: %v", err)
			continue
		}
		fmt.Println(m)
	}
}

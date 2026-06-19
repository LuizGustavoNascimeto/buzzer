package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const tableName = "buzzer-messages"

type User struct {
	UUID        string
	DisplayName string
	Handle      string
}

func getDynamoClient(prod bool) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatal(err)
	}

	if prod {
		return dynamodb.NewFromConfig(cfg)
	}

	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8000")
	})
}

func getUserUUIDs(db *sql.DB) (myUser, otherUser User) {
	rows, err := db.Query(`
		SELECT "user".id, "user".display_name, "user".handle
		FROM "user"
		WHERE "user".handle IN ($1, $2)
	`, "gustavo", "luiz")
	if err != nil {
		log.Fatal("erro ao buscar usuários:", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UUID, &u.DisplayName, &u.Handle); err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}

	for _, u := range users {
		switch u.Handle {
		case "gustavo":
			myUser = u
		case "luiz":
			otherUser = u
		}
	}

	fmt.Printf("get_user_uuids: %+v %+v\n", myUser, otherUser)
	return
}

func createMessageGroup(ctx context.Context, client *dynamodb.Client, messageGroupUUID, myUserUUID, lastMessageAt, message, otherUserUUID, otherUserDisplayName, otherUserHandle string) {
	_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"pk":                 &types.AttributeValueMemberS{Value: "GRP#" + myUserUUID},
			"sk":                 &types.AttributeValueMemberS{Value: lastMessageAt},
			"message_group_uuid": &types.AttributeValueMemberS{Value: messageGroupUUID},
			"message":            &types.AttributeValueMemberS{Value: message},
			"user_uuid":          &types.AttributeValueMemberS{Value: otherUserUUID},
			"user_display_name":  &types.AttributeValueMemberS{Value: otherUserDisplayName},
			"user_handle":        &types.AttributeValueMemberS{Value: otherUserHandle},
		},
	})
	if err != nil {
		log.Fatal("erro ao criar message group:", err)
	}
	fmt.Println("message group criado")
}

func createMessage(ctx context.Context, client *dynamodb.Client, messageGroupUUID, createdAt, message, userUUID, userDisplayName, userHandle string) {
	_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"pk":                &types.AttributeValueMemberS{Value: "MSG#" + messageGroupUUID},
			"sk":                &types.AttributeValueMemberS{Value: createdAt},
			"message_uuid":      &types.AttributeValueMemberS{Value: uuid.New().String()},
			"message":           &types.AttributeValueMemberS{Value: message},
			"user_uuid":         &types.AttributeValueMemberS{Value: userUUID},
			"user_display_name": &types.AttributeValueMemberS{Value: userDisplayName},
			"user_handle":       &types.AttributeValueMemberS{Value: userHandle},
		},
	})
	if err != nil {
		log.Fatal("erro ao criar mensagem:", err)
	}
}

var conversation = `
Person 1: Have you ever watched Babylon 5? It's one of my favorite TV shows!
Person 2: Yes, I have! I love it too. What's your favorite season?
` // ... resto da conversa

func main() {
	prod := true // mude para true para usar produção

	ddbClient := getDynamoClient(prod)

	db, err := sql.Open("pgx", os.Getenv("PROD_POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	messageGroupUUID := "5ae290ed-55d1-47a0-bc6d-fe2bc2700399"
	now := time.Now().UTC()

	myUser, otherUser := getUserUUIDs(db)

	createMessageGroup(ctx, ddbClient, messageGroupUUID,
		myUser.UUID, now.Format(time.RFC3339), "this is a filler message",
		otherUser.UUID, otherUser.DisplayName, otherUser.Handle,
	)

	createMessageGroup(ctx, ddbClient, messageGroupUUID,
		otherUser.UUID, now.Format(time.RFC3339), "this is a filler message",
		myUser.UUID, myUser.DisplayName, myUser.Handle,
	)

	lines := strings.Split(strings.TrimSpace(conversation), "\n")
	for i, line := range lines {
		var key string
		var message string

		if strings.HasPrefix(line, "Person 1: ") {
			key = "my_user"
			message = strings.TrimPrefix(line, "Person 1: ")
		} else if strings.HasPrefix(line, "Person 2: ") {
			key = "other_user"
			message = strings.TrimPrefix(line, "Person 2: ")
		} else {
			log.Fatalf("linha inválida: %s", line)
		}

		var u User
		if key == "my_user" {
			u = myUser
		} else {
			u = otherUser
		}

		createdAt := now.Add(time.Duration(i) * time.Minute).Format(time.RFC3339)
		createMessage(ctx, ddbClient, messageGroupUUID, createdAt, message, u.UUID, u.DisplayName, u.Handle)
	}

	fmt.Println("seed concluído")
}

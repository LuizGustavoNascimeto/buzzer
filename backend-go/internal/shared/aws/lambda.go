package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

type CognitoEvent struct {
	Version       string                 `json:"version"`
	TriggerSource string                 `json:"triggerSource"`
	Region        string                 `json:"region"`
	UserPoolID    string                 `json:"userPoolId"` // ← adicione esse campo
	UserName      string                 `json:"userName"`
	Request       CognitoRequest         `json:"request"`
	Response      map[string]interface{} `json:"response"` // ← mude para map
}

type CognitoRequest struct {
	UserAttributes map[string]string `json:"userAttributes"`
}

func handler(ctx context.Context, event CognitoEvent) (ret CognitoEvent, err error) {
	// Captura qualquer panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC recuperado: %v", r)
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	log.Printf("=== INICIO version=%s triggerSource=%s userPoolId=%s",
		event.Version, event.TriggerSource, event.UserPoolID)

	attrs := event.Request.UserAttributes
	if attrs == nil {
		log.Printf("ERRO: userAttributes é nil")
		return event, nil
	}

	displayName := attrs["name"]
	email := attrs["email"]
	handle := attrs["preferred_username"]
	cognitoID := attrs["sub"]

	log.Printf("displayName=%s email=%s handle=%s cognitoID=%s",
		displayName, email, handle, cognitoID)

	connStr := os.Getenv("CONNECTION_URL")
	if connStr == "" {
		log.Printf("ERRO: CONNECTION_URL vazia")
		return event, nil
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("ERRO sql.Open: %v", err)
		return event, nil
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Printf("ERRO db.Ping: %v", err)
		return event, nil
	}
	log.Printf("DB conectado OK")

	query := `INSERT INTO "user" (display_name, email, handle, cognito_user_id)
              VALUES ($1, $2, $3, $4)`

	_, err = db.ExecContext(ctx, query, displayName, email, handle, cognitoID)
	if err != nil {
		log.Printf("ERRO insert: %v", err)
		return event, nil // ← não retorna erro para não bloquear Cognito
	}

	log.Printf("=== INSERT OK ===")
	return event, nil
}

func main() {
	lambda.Start(handler)
}

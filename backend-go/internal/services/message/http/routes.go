package handler

import (
	"backend-go/internal/middleware"
	"backend-go/internal/services/message/infra"
	"backend-go/internal/services/message/usecase"
	userInfra "backend-go/internal/services/user/infra"
	"backend-go/pkg/auth"
	ddb "backend-go/pkg/dynamo"
	db "backend-go/pkg/gormutil/db"
	"backend-go/pkg/logger"
	"log"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, cwLogger *logger.CloudWatchLogger, validator *auth.Validator) {
	connDDB, err := ddb.GetClient()
	if err != nil {
		log.Fatalf("DynamoDB not initialized before registering message routes: %v", err)
	}
	connDB, err := db.GetConn()
	if err != nil {
		log.Fatal(err)
	}
	repo := infra.NewMessageRepository(connDDB)
	userRepo := userInfra.NewUserRepository(connDB.Gorm)

	service := usecase.NewMessageUsecase(repo, userRepo)
	handler := NewMessageHandler(service, cwLogger)

	protected := rg.Group("")
	protected.Use(middleware.RequireAuth(validator))
	{
		messageRoute := protected.Group("/messages")
		messageRoute.POST("", handler.CreateMessage)
		messageRoute.GET("/:group_id", handler.ListMessages)
	}
	{
		groupRoute := protected.Group("/message_groups")
		groupRoute.GET("/:handle", handler.ListMessageGroups)
	}

}

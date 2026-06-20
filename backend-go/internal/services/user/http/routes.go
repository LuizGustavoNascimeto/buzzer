package handler

import (
	"backend-go/internal/services/user/infra"
	"backend-go/internal/services/user/usecase"
	"backend-go/pkg/auth"
	db "backend-go/pkg/gormutil/db"
	"backend-go/pkg/logger"
	"log"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, cwLogger *logger.CloudWatchLogger, validator *auth.Validator) {
	conn, err := db.GetConn()
	if err != nil {
		log.Fatal(err)
	}

	repo := infra.NewUserRepository(conn.Gorm)
	service := usecase.NewUserUsecase(repo)
	handler := NewUserHandler(service, cwLogger)

	rg.GET("", handler.ListUsers)
	rg.GET("/findByHandle/:handle", handler.FindByHandle)
}

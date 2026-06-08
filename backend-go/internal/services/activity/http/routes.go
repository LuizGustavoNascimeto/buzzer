package handler

import (
	"backend-go/internal/services/activity/infra"
	"backend-go/internal/services/activity/usecase"
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

	repo := infra.NewActivityRepository(conn.Gorm)
	service := usecase.NewActivityUsecase(repo)
	handler := NewActivityHandler(service, cwLogger)

	rg.GET("/home", handler.FindByAll)
	rg.GET("/notifications", handler.FindByHandle)
}

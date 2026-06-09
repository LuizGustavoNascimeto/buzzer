package handler

import (
	"backend-go/internal/middleware"
	"backend-go/internal/services/activity/infra"
	"backend-go/internal/services/activity/usecase"
	userInfra "backend-go/internal/services/user/infra"
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
	userRepo := userInfra.NewUserRepository(conn.Gorm)
	service := usecase.NewActivityUsecase(repo, userRepo)
	handler := NewActivityHandler(service, cwLogger)

	rg.GET("/home", handler.FindByAll)
	protected := rg.Group("")
	protected.Use(middleware.RequireAuth(validator))
	{
		protected.GET("/notifications", handler.FindByHandle)
		protected.POST("", handler.Create)
	}

}

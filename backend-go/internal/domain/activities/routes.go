package activities

import (
	"backend-go/internal/logger"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, cwLogger *logger.CloudWatchLogger) {

	repo := NewMockActivitiesRepo(NewMockActivities()...)
	service := NewActivitiesService(repo)
	handler := NewActivitiesHandler(service, cwLogger)

	rg.GET("/home", handler.ListActivities)
	rg.GET("/notifications", handler.ListNotifications)
}

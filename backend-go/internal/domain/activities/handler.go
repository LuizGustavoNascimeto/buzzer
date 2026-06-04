package activities

import (
	"backend-go/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ActivitiesHandler struct {
	service *ActivitiesService
	log     *logger.CloudWatchLogger
}

func NewActivitiesHandler(service *ActivitiesService, log *logger.CloudWatchLogger) *ActivitiesHandler {
	return &ActivitiesHandler{
		service: service,
		log:     log,
	}
}

func (h *ActivitiesHandler) ListActivities(c *gin.Context) {
	// start := time.Now()
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("feed.type", "home"))

	activities, err := h.service.FindActivities(c.Request.Context())
	if err != nil {
		// h.log.SendLog(logger.LogEntry{
		// 	Timestamp:  c.Request.Context().Value("timestamp").(string),
		// 	Level:      "ERROR",
		// 	Method:     c.Request.Method,
		// 	Path:       c.Request.URL.Path,
		// 	StatusCode: http.StatusInternalServerError,
		// 	Message:    "Failed to list activities",
		// 	Error:      err.Error(),
		// 	ClientIP:   c.ClientIP(),
		// 	Latency:    c.Request.Context().Value("latency").(string),
		// })

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activities)
}

func (h *ActivitiesHandler) ListNotifications(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(
		attribute.String("feed.type", "notifications"),
		attribute.String("user.handle", "Andrew Brown"),
	)

	notifications, err := h.service.FindActivitiesByHandle((c.Request.Context()), "Andrew Brown")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

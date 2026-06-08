package handler

import (
	"backend-go/internal/services/user/usecase"
	"backend-go/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rollbar/rollbar-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UserHandler struct {
	service usecase.IUserService
	log     *logger.CloudWatchLogger
}

func NewUserHandler(service usecase.IUserService, log *logger.CloudWatchLogger) *UserHandler {
	return &UserHandler{
		service: service,
		log:     log,
	}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	// start := time.Now()
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("feed.type", "home"))

	activities, err := h.service.FindAll(c.Request.Context())
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
		rollbar.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activities)
}

func (h *UserHandler) ListByHandle(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(
		attribute.String("feed.type", "notifications"),
		attribute.String("user.handle", "Andrew Brown"),
	)

	notifications, err := h.service.FindByHandle((c.Request.Context()), "Andrew Brown")
	if err != nil {
		rollbar.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

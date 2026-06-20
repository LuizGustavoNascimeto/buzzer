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
	span.SetAttributes(attribute.String("feed.type", "users"))

	users, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		rollbar.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) FindByHandle(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(
		attribute.String("feed.type", "notifications"),
		attribute.String("user.handle", "Andrew Brown"),
	)

	user, err := h.service.FindByHandle((c.Request.Context()), c.Param("handle"))
	if err != nil {
		rollbar.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

package activities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ActivitiesHandler struct {
	service *ActivitiesService
}

func NewActivitiesHandler(service *ActivitiesService) *ActivitiesHandler {
	return &ActivitiesHandler{service: service}
}

func (h *ActivitiesHandler) ListActivities(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("feed.type", "home"))

	activities, err := h.service.FindActivities(c.Request.Context())
	if err != nil {
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

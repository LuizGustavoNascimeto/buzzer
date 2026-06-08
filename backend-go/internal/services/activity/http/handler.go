package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend-go/internal/services/activity/usecase"
	"backend-go/pkg/logger"
)

type ActivityHandler struct {
	service usecase.IActivityService
	log     *logger.CloudWatchLogger // interface do logger também, não o concreto
}

func NewActivityHandler(service usecase.IActivityService, log *logger.CloudWatchLogger) *ActivityHandler {
	return &ActivityHandler{service: service, log: log}
}

func (h *ActivityHandler) Create(c *gin.Context) {
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activity, err := h.service.Create(c.Request.Context(), usecase.CreateActivityInput{
		UserID:              req.UserID,
		Message:             req.Message,
		ReplyToActivityUUID: req.ReplyToActivityUUID,
		ExpiresAt:           req.ExpiresAt,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toActivityResponse(activity))
}

func (h *ActivityHandler) FindByID(c *gin.Context) {
	activity, err := h.service.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toActivityResponse(activity))
}

func (h *ActivityHandler) FindByHandle(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		handleError(c, err)
		return
	}
	activities, err := h.service.FindByUser(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toActivityListResponse(activities))
}

func (h *ActivityHandler) FindByAll(c *gin.Context) {
	activities, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toActivityListResponse(activities))
}

func (h *ActivityHandler) Update(c *gin.Context) {
	var req UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activity, err := h.service.Update(c.Request.Context(), usecase.UpdateActivityInput{
		ID:      c.Param("id"),
		Message: req.Message,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toActivityResponse(activity))
}

func (h *ActivityHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func handleError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

}

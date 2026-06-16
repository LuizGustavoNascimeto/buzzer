package handler

import (
	"backend-go/internal/services/message/domain"

	"backend-go/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	usecase domain.IMessageService
	log     *logger.CloudWatchLogger
}

func NewMessageHandler(repo domain.IMessageService, log *logger.CloudWatchLogger) *MessageHandler {
	return &MessageHandler{usecase: repo, log: log}
}

// func (h *MessageHandler) CreateMessage(c *gin.Context) {
// 	var req *domain.CreateMessage
// 	if err := c.ShouldBindJSON(req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	message, err := h.usecase.CreateMessage(c.Request.Context(), req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.JSON(http.StatusCreated, message)
// }

// func (h *MessageHandler) CreateMessageGroup(c *gin.Context) {
// 	var req *domain.CreateMessageGroup
// 	if err := c.ShouldBindJSON(req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	message, err := h.usecase.CreateMessageGroup(c.Request.Context(), req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, message)
// }

func (h *MessageHandler) ListMessage(c *gin.Context) {
	messages, err := h.usecase.ListMessages(c.Request.Context(), c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (h *MessageHandler) ListMessageGroup(c *gin.Context) {
	messagesGroups, err := h.usecase.ListMessageGroups(c.Request.Context(), c.Param("handle"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messagesGroups)
}

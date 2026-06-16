package handler

import (
	"backend-go/internal/services/message/domain"
	"log"

	"backend-go/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	repo domain.IMessageRepository
	log  *logger.CloudWatchLogger
}

func NewMessageHandler(repo domain.IMessageRepository, log *logger.CloudWatchLogger) *MessageHandler {
	return &MessageHandler{repo: repo, log: log}
}

func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var req *domain.CreateMessage
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message, err := h.repo.CreateMessage(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, message)
}

func (h *MessageHandler) CreateMessageGroup(c *gin.Context) {
	var req *domain.CreateMessageGroup
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message, err := h.repo.CreateMessageGroup(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, message)
}

func (h *MessageHandler) ListMessage(c *gin.Context) {
	var userID string
	if err := c.ShouldBindJSON(&userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	messages, err := h.repo.ListMessages(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (h *MessageHandler) ListMessageGroup(c *gin.Context) {
	var groupID string
	log.Println("um trem ai ", groupID)
	if err := c.ShouldBindJSON(&groupID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "inform a valid group_id"})
		return
	}
	messagesGroups, err := h.repo.ListMessages(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messagesGroups)
}

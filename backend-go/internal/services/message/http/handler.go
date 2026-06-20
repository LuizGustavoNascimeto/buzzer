package handler

import (
	"net/http"

	"backend-go/internal/services/message/domain"
	"backend-go/internal/services/message/usecase"
	"backend-go/pkg/logger"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	usecase usecase.IMessageService
	log     *logger.CloudWatchLogger
}

func NewMessageHandler(uc usecase.IMessageService, log *logger.CloudWatchLogger) *MessageHandler {
	return &MessageHandler{usecase: uc, log: log}
}

func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	input := toCreateMessageInput(req)

	res, err := h.usecase.CreateMessage(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := MessageResponse{
		ID:          res.ID,
		Content:     res.Content,
		DisplayName: res.DisplayName,
		Handle:      res.Handle,
		GroupID:     res.GroupID,
		CreatedAt:   res.SentAt.String()}

	c.JSON(http.StatusCreated, response)
}

func (h *MessageHandler) ListMessages(c *gin.Context) {
	messages, err := h.usecase.ListMessages(c.Request.Context(), c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toMessageResponseList(messages))
}

func (h *MessageHandler) ListMessageGroups(c *gin.Context) {
	groups, err := h.usecase.ListMessageGroups(c.Request.Context(), c.Param("handle"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toMessageGroupResponseList(groups))
}

// --- tradução: Request (HTTP) -> Input (usecase) ---

func toCreateMessageInput(req CreateMessageRequest) usecase.CreateMessageInput {
	return usecase.CreateMessageInput{
		GroupID:        req.GroupID,
		SenderHandle:   req.SenderHandle,
		ReceiverHandle: req.ReceiverHandle,
		Content:        req.Content,
	}
}

// --- tradução: domain (entidade) -> Response (HTTP) ---

func toMessageResponse(m domain.Message) MessageResponse {
	return MessageResponse{
		ID:          m.ID,
		GroupID:     m.GroupID,
		DisplayName: m.DisplayName,
		Handle:      m.Handle,
		Content:     m.Content,
		CreatedAt:   m.SentAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toMessageResponseList(messages []domain.Message) []MessageResponse {
	out := make([]MessageResponse, 0, len(messages))
	for _, m := range messages {
		out = append(out, toMessageResponse(m))
	}
	return out
}

func toMessageGroupResponse(g domain.MessageGroup) MessageGroupResponse {
	return MessageGroupResponse{
		ID:          g.ID,
		DisplayName: g.DisplayName,
		Handle:      g.Handle,
		LastMessage: g.LastMessage,
		LastSentAt:  g.LastSentAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toMessageGroupResponseList(groups []domain.MessageGroup) []MessageGroupResponse {
	out := make([]MessageGroupResponse, 0, len(groups))
	for _, g := range groups {
		out = append(out, toMessageGroupResponse(g))
	}
	return out
}

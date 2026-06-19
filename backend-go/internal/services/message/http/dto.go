package handler

// CreateMessageRequest é o shape exato que chega no corpo do POST.
// json tags ficam SÓ aqui — o usecase nunca vê isso.
type CreateMessageRequest struct {
	GroupID        *string `json:"message_group_uuid"`
	SenderHandle   string  `json:"sender_handle" binding:"required"`
	ReceiverHandle *string `json:"receiver_handle"`
	Content        string  `json:"message" binding:"required"`
}

// MessageResponse é o shape exato que volta no corpo da resposta.
// Você decide aqui o que expor — não precisa ser 1:1 com domain.Message.
type MessageResponse struct {
	ID          string `json:"message_uuid"`
	GroupID     string `json:"message_group_uuid"`
	DisplayName string `json:"user_display_name"`
	Handle      string `json:"user_handle"`
	Content     string `json:"message"`
	CreatedAt   string `json:"created_at"`
}

type MessageGroupResponse struct {
	ID          string `json:"message_group_uuid"`
	DisplayName string `json:"user_display_name"`
	Handle      string `json:"user_handle"`
	LastMessage string `json:"message"`
	LastSentAt  string `json:"created_at"`
}

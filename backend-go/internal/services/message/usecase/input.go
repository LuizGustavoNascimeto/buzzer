package usecase

type CreateMessageGroupInput struct {
	UserID           string
	OtherUserID      string
	OtherDisplayName string
	OtherHandle      string
	Content          string
}

type CreateMessageInput struct {
	GroupID        *string
	SenderHandle   string
	ReceiverHandle *string
	Content        string
}

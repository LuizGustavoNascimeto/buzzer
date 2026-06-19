package infra

type messageItem struct {
	GroupID     string `dynamodbav:"message_group_uuid"`
	UserID      string `dynamodbav:"user_uuid"`
	DisplayName string `dynamodbav:"user_display_name"`
	Handle      string `dynamodbav:"user_handle"`
	Message     string `dynamodbav:"message"`
}

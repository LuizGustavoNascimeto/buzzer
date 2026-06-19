package usecase

type CreateUserInput struct {
	DisplayName   string
	Handle        string
	CognitoUserID string
}

type UpdateUserInput struct {
	ID          string
	DisplayName string
}

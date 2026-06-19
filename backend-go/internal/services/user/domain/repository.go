package domain

import "context"

type IUserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindAll(ctx context.Context) ([]*User, error)
	FindByHandle(ctx context.Context, handle string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	CreateMessageUser(ctx context.Context, senderHandle string, receiverHandle string) ([]*CreateMessageUsers, error)
}

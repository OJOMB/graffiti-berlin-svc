package domain

import "context"

type Repo interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, userID string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUserName(ctx context.Context, userName string) (*User, error)
	UpdateUser(ctx context.Context, user User) error
}

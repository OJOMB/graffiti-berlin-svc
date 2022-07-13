package domain

import "context"

type Repo interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, userID string) (*User, error)
	UpdateUser(ctx context.Context, user User) error
}

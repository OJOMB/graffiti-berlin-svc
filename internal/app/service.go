package app

import (
	"context"

	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/domain"
)

type Service interface {
	CreateUser(ctx context.Context, UserName, Email, Password string) (*domain.User, *domain.Error)
	GetUser(ctx context.Context, userID string) (*domain.User, *domain.Error)
	PatchUser(ctx context.Context, userID string, patch []byte) *domain.Error
}

package app

import "github.com/OJOMB/graffiti-berlin-svc/internal/pkg/auth"

type TokenAuth interface {
	TokenGenerator
	TokenDecoder
}

type TokenGenerator interface {
	GenerateTokenString(userID string) (string, error)
}

type TokenDecoder interface {
	GetClaims(tokenString string) (auth.JWTClaims, error)
}

package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTTool struct {
	secretKey    []byte
	expiresAfter time.Duration
	issuer       string
	idTool       IDGenerator
}

type IDGenerator interface {
	New() (string, error)
}

type JWTClaims struct {
	jwt.RegisteredClaims
}

func NewJWTTool(secretKey string, expiresAfter time.Duration, issuer string) *JWTTool {
	return &JWTTool{
		secretKey:    []byte(secretKey),
		expiresAfter: expiresAfter,
		issuer:       issuer,
	}
}

func (jt *JWTTool) GenerateTokenString(subject string) (string, error) {
	tokenID, err := jt.idTool.New()
	if err != nil {
		return "", fmt.Errorf("failed to generate token ID: %v", err)
	}

	expirationTime := time.Now().Add(jt.expiresAfter)
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			Issuer:    jt.issuer,
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jt.secretKey)
}

func (jt *JWTTool) GetClaims(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) { return jt.secretKey, nil },
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(JWTClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims from token")
	}

	return &claims, nil
}

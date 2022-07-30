package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestJWTToolGenerateSignedTokenString_successPath(t *testing.T) {
	jt := NewJWTTool("supersecretkey", time.Hour, "graffiti-berlin-svc")

	tokenSigned, err := jt.GenerateTokenString("user_id")
	assert.NoError(t, err)
	assert.Regexp(t, `^(?:[\w-]*\.){2}[\w-]*$`, tokenSigned)

	tok, err := jwt.ParseWithClaims(tokenSigned, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jt.secretKey), nil
	})
	assert.NoError(t, err)

	claims, ok := tok.Claims.(*JWTClaims)
	assert.True(t, ok)
	assert.Equal(t, "graffiti-berlin-svc", claims.Issuer)
	assert.Equal(t, "user_id", claims.Subject)
}

func TestJWTToolGetClaims_successPath(t *testing.T) {
	jt := NewJWTTool("supersecretkey", time.Hour, "graffiti-berlin-svc")

	// JWT NumericDate doesn't deal in nanoseconds so we truncate to nearest second
	issuedAt := time.Now().Truncate(time.Second)
	expiresAt := issuedAt.Add(2 * time.Hour)
	userID := "user_id"
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "graffiti-berlin-svc",
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := token.SignedString(jt.secretKey)
	assert.NoError(t, err)

	decodedClaims, err := jt.GetClaims(signedStr)
	assert.NoError(t, err)
	assert.Equal(t, "graffiti-berlin-svc", decodedClaims.Issuer)
	assert.Equal(t, "user_id", decodedClaims.Subject)
	assert.Equal(t, issuedAt, decodedClaims.IssuedAt.Time)
	assert.Equal(t, expiresAt, decodedClaims.ExpiresAt.Time)
}

func TestJWTToolGetClaims_expiredToken_successPath(t *testing.T) {
	jt := NewJWTTool("supersecretkey", time.Hour, "graffiti-berlin-svc")

	// JWT NumericDate doesn't deal in nanoseconds so we truncate to nearest second
	issuedAt := time.Time{}
	expiresAt := issuedAt.Add(time.Hour)
	userID := "user_id"
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "graffiti-berlin-svc",
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := token.SignedString(jt.secretKey)
	assert.NoError(t, err)

	decodedClaims, err := jt.GetClaims(signedStr)
	assert.Nil(t, decodedClaims)
	assert.Error(t, err)

	assert.True(t, strings.HasPrefix(err.Error(), "failed to parse token: token is expired by"))
}

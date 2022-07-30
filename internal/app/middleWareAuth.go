package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type TokenValidator struct {
	logger  *logrus.Entry
	decoder TokenDecoder
}

func NewTokenValidator(l *logrus.Entry) *TokenValidator {
	return &TokenValidator{logger: l.WithField("middleware", "TokenValidator")}
}

func (tv *TokenValidator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeaderVal := r.Header.Get("Authorization")
		if authHeaderVal == "" {
			tv.logger.Info("no token found in request")
			appErr := newAppErr("no token found in request", http.StatusUnauthorized)
			http.Error(w, appErr.Error(), appErr.code)
			return
		}

		tokenString := strings.TrimPrefix(authHeaderVal, "Bearer ")
		if authHeaderVal == tokenString {
			appErr := newAppErr("auth header value in unexpected format", http.StatusUnauthorized)
			http.Error(w, appErr.Error(), appErr.code)
			return
		}

		claims, err := tv.decoder.GetClaims(tokenString)
		if err != nil {
			appErr := newAppErr(fmt.Sprintf("invalid token: %s", err.Error()), http.StatusUnauthorized)
			http.Error(w, appErr.Error(), appErr.code)
			return
		}

		// include the subject in the context as this is needed at the domain level
		r.WithContext(context.WithValue(r.Context(), "subject", claims.Subject))

		next.ServeHTTP(w, r)
	})
}

package app

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

const reqIDHeader = "X-Request-ID"

type middlewareRequestID struct {
	logger *logrus.Entry
	idTool nanoID
}

func NewmiddlewareRequestID(l *logrus.Entry) *middlewareRequestID {
	return &middlewareRequestID{logger: l.WithField("middleware", "middlewareRequestID")}
}

func (mw *middlewareRequestID) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(reqIDHeader)
		if reqID == "" {
			generatedID, err := mw.idTool.New()
			if err != nil {
				mw.logger.Errorf("failed to generate request ID: %s", err.Error())
				next.ServeHTTP(w, r)
				return
			}

			reqID = generatedID
		}

		r.WithContext(context.WithValue(r.Context(), reqIDHeader, reqID))
	})
}

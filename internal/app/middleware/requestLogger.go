package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const reqIDHeader = "X-Request-ID"

type responseWriterRecorder struct {
	http.ResponseWriter
	status int
	body   string
}

func (w *responseWriterRecorder) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterRecorder) Write(bs []byte) (int, error) {
	w.body += string(bs)
	return w.ResponseWriter.Write(bs)
}

type RequestResponseLogger struct {
	logger          *logrus.Entry
	logResponseBody bool
}

func NewRequestResponseLogger(l *logrus.Entry) *RequestResponseLogger {
	return &RequestResponseLogger{logger: l.WithField("middleware", "RequestResponseLogger")}
}

func (rrl *RequestResponseLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(reqIDHeader)

		rrl.logger.WithFields(
			logrus.Fields{
				"type":             "request",
				"path":             r.URL.Path,
				"query_parameters": r.URL.RawQuery,
				"method":           r.Method,
				"request_id":       reqID,
				"src_host":         r.RemoteAddr,
			},
		).Info("incoming request")

		start := time.Now().UTC()

		customResponseWriter := &responseWriterRecorder{ResponseWriter: w}
		next.ServeHTTP(customResponseWriter, r)

		end := time.Now().UTC()

		rrl.logger.WithFields(
			logrus.Fields{
				"type":          "response",
				"status":        customResponseWriter.status,
				"response_time": end.Sub(start).Milliseconds(),
				"path":          r.URL.Path,
				"method":        r.Method,
				"body":          customResponseWriter.body,
				"request_id":    reqID,
				"tgt_host":      r.Host,
			},
		).Info("outgoing response")
	})
}

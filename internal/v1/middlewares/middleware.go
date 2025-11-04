package middewares

import (
	"context"
	"net/http"
	"orderService/pkg/logger"

	"github.com/google/uuid"
)

func LoggingMiddleware(ctx context.Context) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqIDHeader := w.Header().Get("x-request-id")
			if reqIDHeader != "" {
				logger.WithRequestID(ctx, reqIDHeader)
			} else {
				logger.WithRequestID(ctx, uuid.NewString())
			}
			next.ServeHTTP(w, r)
		})
	}
}

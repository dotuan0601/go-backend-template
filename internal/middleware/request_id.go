package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	PathKey      contextKey = "path"
	MethodKey    contextKey = "method"
)

func RequestID(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		requestID := r.Header.Get(
			"X-Request-ID",
		)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		r.Header.Set(
			"X-Request-ID",
			requestID,
		)

		ctx := context.WithValue(
			r.Context(),
			RequestIDKey,
			requestID,
		)

		r = r.WithContext(ctx)

		w.Header().Set(
			"X-Request-ID",
			requestID,
		)

		next.ServeHTTP(
			w,
			r,
		)
	})
}

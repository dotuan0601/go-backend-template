package middleware

import (
	"net/http"
	"time"

	"log/slog"
)

type responseWriter struct {
	http.ResponseWriter

	statusCode int
}

func (rw *responseWriter) WriteHeader(
	code int,
) {

	rw.statusCode = code

	rw.ResponseWriter.WriteHeader(
		code,
	)
}

func Logger(
	log *slog.Logger,
) func(http.Handler) http.Handler {

	return func(
		next http.Handler,
	) http.Handler {

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			start := time.Now()

			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(
				rw,
				r,
			)

			latency := time.Since(
				start,
			)

			requestID := r.Header.Get(
				"X-Request-ID",
			)

			log.Info(
				"request completed",
				"request_id",
				requestID,
				"method",
				r.Method,
				"path",
				r.URL.Path,
				"status",
				rw.statusCode,
				"latency",
				latency.String(),
			)

		})
	}
}

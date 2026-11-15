package handler

import (
	"log/slog"
	"net/http"

	"github.com/dotuan0601/go-backend-template/internal/middleware"
	"github.com/dotuan0601/go-backend-template/internal/response"
)

type HealthHandler struct {
	logger *slog.Logger
}

func NewHealthHandler(
	logger *slog.Logger,
) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

func (h *HealthHandler) Health(
	w http.ResponseWriter,
	r *http.Request,
) {
	requestID := r.Context().Value(
		middleware.RequestIDKey,
	)
	h.logger.Info(
		"health check called",
		"method",
		r.Method,
		"path",
		r.URL.Path,
		"request_id",
		requestID,
	)

	data := map[string]string{
		"status": "ok",
	}

	response.Success(
		w,
		data,
	)

}

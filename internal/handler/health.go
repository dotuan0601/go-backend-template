package handler

import (
	"log/slog"
	"net/http"

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

	h.logger.Info(
		"health check called",
	)

	data := map[string]interface{}{
		"status": "ok",
	}

	response.Success(
		w,
		data,
	)
}

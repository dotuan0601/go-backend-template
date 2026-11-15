package router

import (
	"net/http"

	"github.com/dotuan0601/go-backend-template/internal/handler"
)

func RegisterRoutes(
	mux *http.ServeMux,
	healthHandler *handler.HealthHandler,
) {
	mux.HandleFunc("GET /health", healthHandler.Health)
}

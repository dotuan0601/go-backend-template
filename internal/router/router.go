package router

import (
	"net/http"

	"github.com/dotuan0601/go-backend-template/internal/handler"
)

func RegisterRoutes(
	mux *http.ServeMux,
	healthHandler *handler.HealthHandler,
	userHandler *handler.UserHandler,
) {
	mux.HandleFunc("GET /health", healthHandler.Health)
	// list users
	mux.HandleFunc(
		"GET /users",
		userHandler.GetUsers,
	)

	// user detail
	mux.HandleFunc(
		"GET /users/{id}",
		userHandler.GetUser,
	)
}

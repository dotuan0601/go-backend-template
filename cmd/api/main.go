package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dotuan0601/go-backend-template/internal/config"
	"github.com/dotuan0601/go-backend-template/internal/handler"
	"github.com/dotuan0601/go-backend-template/internal/logger"
	"github.com/dotuan0601/go-backend-template/internal/middleware"
	"github.com/dotuan0601/go-backend-template/internal/router"
)

func main() {
	env := flag.String(
		"env",
		"dev",
		"application environment",
	)

	flag.Parse()

	configPath := fmt.Sprintf(
		"configs/config-%s.yml",
		*env,
	)

	cfg, err := config.Load(configPath)
	if err != nil {
		panic(err)
	}

	log := logger.Init(*env)

	mux := http.NewServeMux()

	healthHandler := handler.NewHealthHandler(log)

	router.RegisterRoutes(
		mux,
		healthHandler,
	)

	port := cfg.Server.Port
	addr := fmt.Sprintf(":%d", port)

	log.Info(
		"starting api server",
		"addr",
		addr,
	)

	server := &http.Server{
		Addr:    addr,
		Handler: middleware.RequestID(mux),
	}

	serverErr := make(chan error)

	go func() {
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {
	case err := <-serverErr:
		fmt.Printf("server error: %v\n", err)

	case sig := <-quit:
		fmt.Printf("received signal: %v\n", sig)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error(
			"failed to shutdown server",
			"error",
			err,
		)
	}

	fmt.Println("server stopped")
}

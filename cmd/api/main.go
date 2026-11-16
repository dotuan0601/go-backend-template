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
	"github.com/dotuan0601/go-backend-template/internal/datasource/mysql"
	"github.com/dotuan0601/go-backend-template/internal/datasource/redis"
	"github.com/dotuan0601/go-backend-template/internal/handler"
	"github.com/dotuan0601/go-backend-template/internal/logger"
	"github.com/dotuan0601/go-backend-template/internal/middleware"
	"github.com/dotuan0601/go-backend-template/internal/repository"
	"github.com/dotuan0601/go-backend-template/internal/router"
	"github.com/dotuan0601/go-backend-template/internal/service"
)

func main() {

	env := flag.String(
		"env",
		"dev",
		"application environment",
	)

	flag.Parse()

	// =========================
	// Config
	// =========================

	configPath := fmt.Sprintf(
		"configs/config-%s.yml",
		*env,
	)

	cfg, err := config.Load(configPath)

	if err != nil {
		panic(err)
	}

	// =========================
	// Logger
	// =========================

	log := logger.Init(*env)

	// =========================
	// Datasource
	// =========================

	mysqlDS, err := mysql.New(
		cfg.MySQL,
	)

	if err != nil {
		log.Error(
			"failed to init mysql",
			"error",
			err,
		)

		panic(err)
	}

	defer func() {

		if err := mysqlDS.Close(); err != nil {

			log.Error(
				"failed to close mysql",
				"error",
				err,
			)

		}

	}()

	redisDS, err := redis.New(
		cfg.Redis.Host,
		cfg.Redis.Port,
	)

	if err != nil {

		log.Error(
			"failed to init redis",
			"error",
			err,
		)

		panic(err)
	}

	defer redisDS.Close()

	// =========================
	// Repository
	// =========================

	userRepository := repository.NewUserRepository(
		mysqlDS.DB(),
	)

	// =========================
	// Service
	// =========================

	userService := service.NewUserService(
		userRepository,
		redisDS.Client(),
	)

	// =========================
	// Handler
	// =========================

	healthHandler := handler.NewHealthHandler(
		log,
	)

	userHandler := handler.NewUserHandler(
		userService,
		log,
	)

	// =========================
	// Router
	// =========================

	mux := http.NewServeMux()

	router.RegisterRoutes(
		mux,
		healthHandler,
		userHandler,
	)

	// =========================
	// Middleware
	// =========================

	httpHandler := middleware.RequestID(
		mux,
	)

	httpHandler = middleware.Logger(
		log,
	)(
		httpHandler,
	)

	// =========================
	// HTTP Server
	// =========================

	addr := fmt.Sprintf(
		":%d",
		cfg.Server.Port,
	)

	server := &http.Server{
		Addr:    addr,
		Handler: httpHandler,
	}

	log.Info(
		"starting api server",
		"addr",
		addr,
	)

	serverErr := make(chan error, 1)

	go func() {

		err := server.ListenAndServe()

		if err != nil &&
			err != http.ErrServerClosed {

			serverErr <- err
		}

	}()

	// =========================
	// Shutdown signal
	// =========================

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {

	case err := <-serverErr:

		log.Error(
			"server error",
			"error",
			err,
		)

	case sig := <-quit:

		log.Info(
			"shutdown signal received",
			"signal",
			sig,
		)

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

	log.Info(
		"server stopped",
	)

}

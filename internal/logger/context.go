package logger

import (
	"context"
	"log/slog"

	"github.com/dotuan0601/go-backend-template/internal/middleware"
)

func Error(
	ctx context.Context,
	log *slog.Logger,
	msg string,
	err error,
) {

	requestID, _ := ctx.Value(
		middleware.RequestIDKey,
	).(string)

	log.Error(
		msg,
		"request_id",
		requestID,
		"error",
		err,
	)
}

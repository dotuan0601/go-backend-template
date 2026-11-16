package logger

import (
	"log/slog"
	"os"
)

func Init(
	env string,
) *slog.Logger {

	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	)

	log := slog.New(handler)

	slog.SetDefault(log)

	return log
}

package logger

import (
	"log/slog"
	"os"
)

func Init(env string) *slog.Logger {

	var level slog.Level

	if env == "dev" {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: level,
		},
	)

	log := slog.New(handler)

	slog.SetDefault(log)

	return log
}

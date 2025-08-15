package log

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Level     string
	AppName   string
	AddSource bool
}

func New(cfg Config) *slog.Logger {
	var level slog.Level

	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
	})).With("app_name", cfg.AppName)

	logger.Info("Logger initialized", "level", cfg.Level, "app_name", cfg.AppName)
	return logger
}

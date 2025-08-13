package log

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"
	"strings"
)

type Config struct {
	Level     string
	AppName   string
	AddSource bool
}

const logDir = "./logs"

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

	timestamp := time.Now().Format("2006-01-02")
	logFile := filepath.Join(logDir, cfg.AppName+"-"+timestamp+".log")

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}

	logger := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
	})).With("app_name", cfg.AppName)

	logger.Info("Logger initialized", "level", cfg.Level, "app_name", cfg.AppName, "log_file", logFile)
	return logger
}
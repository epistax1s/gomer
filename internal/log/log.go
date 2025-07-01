package log

import (
	"io"
	"log/slog"
	"os"

	"github.com/epistax1s/gomer/internal/config"
	"github.com/epistax1s/gomer/internal/log/roll"
)

var log *slog.Logger

func InitLogger(config *config.LogConfig) {
	logger := &roll.Logger{
		Filename:   "./log/gomer.log",
		MaxSize:    2,
		MaxBackups: 1,
		MaxAge:     360,
		Compress:   true,
	}

	var writers []io.Writer

	writers = append(writers, logger)

	if config.Stdout {
		writers = append(writers, os.Stdout)
	}

	multiWriter := io.MultiWriter(writers...)

	handler := slog.NewTextHandler(
		multiWriter,
		&slog.HandlerOptions{
			Level: parseLogLevel(config.Level),
		},
	)

	log = slog.New(handler)
}

func parseLogLevel(logLevel string) slog.Level {
	switch logLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Debug(msg string, args ...any) {
	log.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	log.Error(msg, args...)
}

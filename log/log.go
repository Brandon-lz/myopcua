package log

import (
	"log/slog"
	"os"
)

var logger *slog.Logger
var logLevel slog.Level

func Init(level ...slog.Level) {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	for _,l := range level {
		SetLevel(l)
	}
	slog.SetDefault(logger)
}

func SetLevel(level slog.Level) {
	logLevel = level
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)
}
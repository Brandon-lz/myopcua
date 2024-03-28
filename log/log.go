package log

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger
var LogLevel slog.Level

func Init(level ...slog.Level) {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	for _,l := range level {
		SetLevel(l)
	}
}

func SetLevel(level slog.Level) {
	LogLevel = level
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: LogLevel}))
}
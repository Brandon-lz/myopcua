package log

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func Init() {
    Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
    // Logger.Info("hello, world", "user", os.Getenv("USER"))
}
package log

import (
	"errors"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger

// var loggerWriter io.Writer = os.Stdout

func Init(levelName string) error {

	var level = slog.Level(0)
	switch strings.ToUpper(levelName) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		return errors.New("unknown name")
	}

	// Create a multi-writer that writes to both the log file and DataDog.
	loggerWriter := getLogWriter()
	logger = slog.New(slog.NewJSONHandler(loggerWriter, &slog.HandlerOptions{Level: level}))

	// 添加全局字段
	logger = logger.With("deploy_env", os.Getenv("deploy_env"))
	slog.SetDefault(logger)

	return nil
}

func getLogWriter() io.Writer {
	fileName := "sys.log"
	var logFile *os.File
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		logFile, err = os.Create("sys.log")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		logFile, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatal(err)
		}
	}
	// defer logFile.Close()

	// Create a multi-writer that writes to both the log file and DataDog.
	// loggerWriter := io.MultiWriter(os.Stdout, logFile, NewInfluxWriter())
	loggerWriter := io.MultiWriter(os.Stdout, logFile, NewElasticsearchWriter())
	return loggerWriter
	// logger = slog.New(slog.NewJSONHandler(multiWriter, nil))
}

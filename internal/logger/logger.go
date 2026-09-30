package logger

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Init starts a new logger.
// This logger is set as the default.
// Returns the logger and a function to close the rotator. The returned function
// should be deferred from the caller so that the rotator is closed when caller
// funciton exits.
func Init(logPath string, appName string) (*slog.Logger, func() error) {
	rotator := &lumberjack.Logger{
		Filename:  logPath,
		MaxSize:   10, // MB before rotation
		MaxAge:    30, // days
		LocalTime: true,
		Compress:  true, // gzip old files
	}

	// Write to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, rotator)

	logger := slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})).With("app", appName)

	slog.SetDefault(logger)
	logger.Info("Initialized logger")

	return logger, rotator.Close
}

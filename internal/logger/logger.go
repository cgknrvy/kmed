package logger

import (
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

	// Since the app is compiled with `-ldflags -H=windowsgui` in windows it errors
	// when it tries to compile to the console. Hence having two handlers aids
	// in ensuring that logging still happens to the file and isn't blocked when the console
	// writes fail.

	fileHandler := slog.NewJSONHandler(rotator, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})

	logger := slog.New(slog.NewMultiHandler(fileHandler, consoleHandler)).With("app", appName)

	slog.SetDefault(logger)
	logger.Info("Initialized logger")

	return logger, rotator.Close
}

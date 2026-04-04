// internal/infra/logger/logger.go
package logger

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(filePath string) (*slog.Logger, io.Writer) {
	lumberjackLog := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, lumberjackLog)

	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	return slog.New(handler), multiWriter
}

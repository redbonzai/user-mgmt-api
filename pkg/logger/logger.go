package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// InitLogger initializes the logger
func InitLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	log, err = config.Build()
	if err != nil {
		panic(err)
	}
}

// Info logs an info message
func Info(message string, fields ...zapcore.Field) {
	log.Info(message, fields...)
}

// Error logs an error message
func Error(message string, fields ...zapcore.Field) {
	log.Error(message, fields...)
}

// Debug logs a debug message
func Debug(message string, fields ...zapcore.Field) {
	log.Debug(message, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(message string, fields ...zapcore.Field) {
	log.Fatal(message, fields...)
}

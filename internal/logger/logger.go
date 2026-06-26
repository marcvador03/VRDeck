package logger

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *zap.Logger
	once          sync.Once
	filePath      string
	logLevel      zapcore.Level = zap.InfoLevel
)

func InitLogger(path string, level zapcore.Level) {
	filePath = path
	logLevel = level
}

func GetDefaultLogger() *zap.Logger {
	once.Do(func() {
		// Open the log file
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Errorf("failed to open log file %s: %w", filePath, err))
		}

		// Configure the encoder (JSON)
		config := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(config),
			zapcore.AddSync(file),
			logLevel, // Use the configured log level
		)

		defaultLogger = zap.New(core, zap.AddCaller())
	})
	return defaultLogger
}

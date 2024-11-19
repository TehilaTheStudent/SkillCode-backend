package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)



func InitLogger() *zap.Logger {
	// Set up log rotation with lumberjack
	logFile := &lumberjack.Logger{
		Filename:   "./logs/app.log", // Log file path
		MaxSize:    10,               // Max size in MB before rotation
		MaxBackups: 5,                // Number of old log files to keep
		MaxAge:     7,                // Max age in days before deletion
		Compress:   true,             // Compress rotated files
	}

	// Set up Zap core with file writer only
	fileWriteSyncer := zapcore.AddSync(logFile)

	// Configure human-readable log format for development
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Use readable time format
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "message"
	encoderConfig.CallerKey = "caller"

	// Create core for file
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // Human-readable format
		fileWriteSyncer,                          // Write to file only
		zapcore.DebugLevel,                       // Log level (capture everything in development)
	)

	// Create and return the logger
	return zap.New(core, zap.AddCaller())
}

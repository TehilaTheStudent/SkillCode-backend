package config

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

const LogDir = "./logs"


func InitLogger() (*zap.Logger, error) {
	// Log file path
	logFilePath := LogDir + "/skillcode.log"

	// Set up log rotation with lumberjack
	logFile := &lumberjack.Logger{
		Filename:   logFilePath, // Log file path
		MaxSize:    10,          // Max size in MB before rotation
		MaxBackups: 5,           // Number of old log files to keep
		MaxAge:     7,           // Max age in days before deletion
		Compress:   true,        // Compress rotated files
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

	// Create the logger
	logger := zap.New(core, zap.AddCaller())

	// Ensure the log file has the correct permissions
	if err := ensureFilePermissions(logFilePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to set file permissions: %w", err)
	}

	// Check if logger is nil
	if logger == nil {
		return nil, fmt.Errorf("failed to initialize logger")
	}

	return logger, nil
}

// ensureFilePermissions ensures the log file has the desired permissions
func ensureFilePermissions(filePath string, perms os.FileMode) error {
	if err := os.MkdirAll(LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	if _, err := os.Stat(filePath); err == nil {
		// File exists; set permissions
		if err := os.Chmod(filePath, perms); err != nil {
			return fmt.Errorf("failed to set file permissions: %w", err)
		}
	}
	return nil
}

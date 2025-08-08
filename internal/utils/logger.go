package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Logger struct {
	errorLogger *log.Logger
	infoLogger  *log.Logger
	errorFile   *os.File
}

var AppLogger *Logger

func InitLogger() error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Open error log file
	errorFile, err := os.OpenFile(
		filepath.Join("logs", "error.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		return fmt.Errorf("failed to open error log file: %v", err)
	}

	// Create loggers
	errorLogger := log.New(errorFile, "", 0)
	infoLogger := log.New(os.Stdout, "", 0)

	AppLogger = &Logger{
		errorLogger: errorLogger,
		infoLogger:  infoLogger,
		errorFile:   errorFile,
	}

	return nil
}

func (l *Logger) LogError(message string, err error, context map[string]interface{}) {
	// Get caller information
	_, file, line, _ := runtime.Caller(1)
	filename := filepath.Base(file)

	// Build log entry
	logEntry := fmt.Sprintf("[ERROR] %s | %s:%d | %s",
		time.Now().Format("2006-01-02 15:04:05"),
		filename,
		line,
		message,
	)

	if err != nil {
		logEntry += fmt.Sprintf(" | Error: %v", err)
	}

	if context != nil {
		logEntry += fmt.Sprintf(" | Context: %+v", context)
	}

	// Write to both error file and stdout
	l.errorLogger.Println(logEntry)
	l.infoLogger.Println(logEntry)
}

func (l *Logger) LogInfo(message string, context map[string]interface{}) {
	// Get caller information
	_, file, line, _ := runtime.Caller(1)
	filename := filepath.Base(file)

	// Build log entry
	logEntry := fmt.Sprintf("[INFO] %s | %s:%d | %s",
		time.Now().Format("2006-01-02 15:04:05"),
		filename,
		line,
		message,
	)

	if context != nil {
		logEntry += fmt.Sprintf(" | Context: %+v", context)
	}

	l.infoLogger.Println(logEntry)
}

func (l *Logger) LogWarning(message string, context map[string]interface{}) {
	// Get caller information
	_, file, line, _ := runtime.Caller(1)
	filename := filepath.Base(file)

	// Build log entry
	logEntry := fmt.Sprintf("[WARNING] %s | %s:%d | %s",
		time.Now().Format("2006-01-02 15:04:05"),
		filename,
		line,
		message,
	)

	if context != nil {
		logEntry += fmt.Sprintf(" | Context: %+v", context)
	}

	l.infoLogger.Println(logEntry)
}

func (l *Logger) Close() {
	if l.errorFile != nil {
		l.errorFile.Close()
	}
}

// Convenience functions for global logger
func LogError(message string, err error, context map[string]interface{}) {
	if AppLogger != nil {
		AppLogger.LogError(message, err, context)
	}
}

func LogInfo(message string, context map[string]interface{}) {
	if AppLogger != nil {
		AppLogger.LogInfo(message, context)
	}
}

func LogWarning(message string, context map[string]interface{}) {
	if AppLogger != nil {
		AppLogger.LogWarning(message, context)
	}
}

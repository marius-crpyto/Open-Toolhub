package logger

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	levelStr := "debug"
	outputPath := "./log"
	logFileName := "test.log"

	logger, err := New(levelStr, outputPath, logFileName)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if logger.Logger == nil {
		t.Fatalf("New() logger.Logger is nil")
	}

	if logger.atomicLevel.Level() != zapcore.DebugLevel {
		t.Fatalf("New() logger.atomicLevel.Level() = %v, want %v", logger.atomicLevel.Level(), zapcore.DebugLevel)
	}

	if logger.file == nil {
		t.Fatalf("New() logger.file is nil")
	}

	// Check if the log file exists
	if _, err := os.Stat(filepath.Join(outputPath, logFileName)); os.IsNotExist(err) {
		t.Fatalf("New() log file %s does not exist", filepath.Join(outputPath, logFileName))
	}

	logger.Logger.Sync()

	logger.Info("test log")

	logContent, err := os.ReadFile(filepath.Join(outputPath, logFileName))
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if !bytes.Contains(logContent, []byte("test log")) {
		t.Fatalf("Log file does not contain expected log message")
	}

}

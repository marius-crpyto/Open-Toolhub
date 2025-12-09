package common

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelError LogLevel = "error"
	LogLevelDebug LogLevel = "debug"
	LogLevelWarn  LogLevel = "warn"
)

type LoggerCallback func(level LogLevel, format string, v ...interface{})

func WaitSignal(cancelFuncs []context.CancelFunc, logger LoggerCallback) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	defer close(signals)

	for sig := range signals {
		logger(LogLevelInfo, "Received signal: %s", sig.String())
		switch sig {
		case os.Interrupt, os.Kill:
			exitStatus := 0
			for _, cancel := range cancelFuncs {
				cancel()
			}

			logger(LogLevelInfo, "waiting for sub-tasks to clean up...")
			time.Sleep(2 * time.Second)

			logger(LogLevelInfo, "application exited gracefully")
			os.Exit(exitStatus)
		default:
			logger(LogLevelInfo, "unexpected signal, ignoring...")
		}
	}
}

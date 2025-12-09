package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Logger      *zap.Logger
	atomicLevel zap.AtomicLevel
	file        *os.File
}

func New(levelStr, outputPath, logFileName string) (*Logger, error) {
	al := zap.NewAtomicLevel()
	if err := al.UnmarshalText([]byte(levelStr)); err != nil {
		return nil, err
	}

	consoleEncCfg := zap.NewProductionEncoderConfig()
	consoleEncCfg.TimeKey = "timestamp"
	consoleEncCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncCfg.CallerKey = "caller"
	consoleEncCfg.EncodeCaller = zapcore.ShortCallerEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncCfg)

	var cores []zapcore.Core
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		al,
	)
	cores = append(cores, consoleCore)

	var file *os.File
	if outputPath != "" {
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return nil, err
		}

		if logFileName == "" {
			logFileName = "app.log"
		}

		fp := filepath.Join(outputPath, logFileName)
		var err error
		file, err = os.OpenFile(fp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}

		fileEncCfg := zap.NewProductionEncoderConfig()
		fileEncCfg.TimeKey = "timestamp"
		fileEncCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		fileEncCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		fileEncCfg.CallerKey = "caller"
		fileEncCfg.EncodeCaller = zapcore.ShortCallerEncoder
		fileEncoder := zapcore.NewConsoleEncoder(fileEncCfg)

		fileCore := zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(file),
			al,
		)
		cores = append(cores, fileCore)
	}

	core := zapcore.NewTee(cores...)
	z := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return &Logger{Logger: z, atomicLevel: al, file: file}, nil
}

func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

func (l *Logger) Close() error {
	_ = l.Logger.Sync()
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *Logger) SetLevel(level zapcore.Level) {
	l.atomicLevel.SetLevel(level)
}

func (l *Logger) Info(msg string, fields ...zap.Field)  { l.Logger.Info(msg, fields...) }
func (l *Logger) Debug(msg string, fields ...zap.Field) { l.Logger.Debug(msg, fields...) }
func (l *Logger) Warn(msg string, fields ...zap.Field)  { l.Logger.Warn(msg, fields...) }
func (l *Logger) Error(msg string, fields ...zap.Field) { l.Logger.Error(msg, fields...) }
func (l *Logger) Fatal(msg string, fields ...zap.Field) { l.Logger.Fatal(msg, fields...) }
func (l *Logger) Panic(msg string, fields ...zap.Field) { l.Logger.Panic(msg, fields...) }

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger:      l.Logger.With(fields...),
		atomicLevel: l.atomicLevel,
		file:        l.file,
	}
}

func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.Logger.Sugar()
}

package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger to provide structured logging capabilities.
type Logger struct {
	*zap.Logger
}

// Config contains logger configuration.
type Config struct {
	Level      string // debug, info, warn, error
	Format     string // json, console
	OutputPath string // stdout, stderr, or file path
}

// New creates a new Logger instance based on the provided configuration.
func New(cfg Config) (*Logger, error) {
	// Parse log level
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	// Build encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Adjust encoder config based on format
	if cfg.Format != "json" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Create output sink
	outputPath := cfg.OutputPath
	if outputPath == "" {
		outputPath = "stdout"
	}

	// Build zap config
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         cfg.Format,
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{outputPath},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Build logger
	zapLogger, err := zapConfig.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &Logger{Logger: zapLogger}, nil
}

// parseLevel converts a string level to zapcore.Level.
func parseLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("invalid log level: %s", level)
	}
}

// WithFields returns a new logger with the given fields.
func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.Logger.With(fields...)}
}

// WithError adds an error field to the logger.
func (l *Logger) WithError(err error) *Logger {
	return &Logger{Logger: l.Logger.With(zap.Error(err))}
}

// WithComponent adds a component field to the logger.
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{Logger: l.Logger.With(zap.String("component", component))}
}

// WithRequestID adds a request ID field to the logger.
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{Logger: l.Logger.With(zap.String("request_id", requestID))}
}

// WithSpecID adds a spec ID field to the logger.
func (l *Logger) WithSpecID(specID string) *Logger {
	return &Logger{Logger: l.Logger.With(zap.String("spec_id", specID))}
}

// NewNop creates a no-op logger for testing.
func NewNop() *Logger {
	return &Logger{Logger: zap.NewNop()}
}

// Sync flushes any buffered log entries.
// This should be called before the application exits.
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// Helper functions for creating zap fields

// String creates a string field.
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

// Int creates an int field.
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Int64 creates an int64 field.
func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

// Bool creates a bool field.
func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

// Duration creates a duration field.
func Duration(key string, val time.Duration) zap.Field {
	return zap.Duration(key, val)
}

// Any creates a field with any type.
func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

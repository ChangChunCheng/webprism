package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "json format info level",
			config: Config{
				Level:      "info",
				Format:     "json",
				OutputPath: "stdout",
			},
			wantErr: false,
		},
		{
			name: "console format debug level",
			config: Config{
				Level:      "debug",
				Format:     "console",
				OutputPath: "stdout",
			},
			wantErr: false,
		},
		{
			name: "warn level",
			config: Config{
				Level:      "warn",
				Format:     "json",
				OutputPath: "stdout",
			},
			wantErr: false,
		},
		{
			name: "error level",
			config: Config{
				Level:      "error",
				Format:     "json",
				OutputPath: "stdout",
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			config: Config{
				Level:      "invalid",
				Format:     "json",
				OutputPath: "stdout",
			},
			wantErr: true,
			errMsg:  "invalid log level",
		},
		{
			name: "empty output path defaults to stdout",
			config: Config{
				Level:      "info",
				Format:     "json",
				OutputPath: "",
			},
			wantErr: false,
		},
		{
			name: "stderr output",
			config: Config{
				Level:      "info",
				Format:     "json",
				OutputPath: "stderr",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)
				assert.NotNil(t, logger.Logger)

				// Clean up
				if logger != nil {
					_ = logger.Sync()
				}
			}
		})
	}
}

func TestNew_FileOutput(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	config := Config{
		Level:      "info",
		Format:     "json",
		OutputPath: logFile,
	}

	logger, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Write a log message
	logger.Info("test message")
	_ = logger.Sync()

	// Verify file was created
	_, err = os.Stat(logFile)
	assert.NoError(t, err, "log file should be created")
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		wantLevel zapcore.Level
		wantErr   bool
	}{
		{
			name:      "debug level",
			level:     "debug",
			wantLevel: zapcore.DebugLevel,
			wantErr:   false,
		},
		{
			name:      "info level",
			level:     "info",
			wantLevel: zapcore.InfoLevel,
			wantErr:   false,
		},
		{
			name:      "warn level",
			level:     "warn",
			wantLevel: zapcore.WarnLevel,
			wantErr:   false,
		},
		{
			name:      "error level",
			level:     "error",
			wantLevel: zapcore.ErrorLevel,
			wantErr:   false,
		},
		{
			name:      "invalid level",
			level:     "invalid",
			wantLevel: zapcore.InfoLevel,
			wantErr:   true,
		},
		{
			name:      "empty level",
			level:     "",
			wantLevel: zapcore.InfoLevel,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := parseLevel(tt.level)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLevel, level)
			}
		})
	}
}

func TestLogger_WithFields(t *testing.T) {
	logger := NewNop()

	newLogger := logger.WithFields(
		zap.String("key1", "value1"),
		zap.Int("key2", 42),
	)

	assert.NotNil(t, newLogger)
	// Should not panic when logging
	assert.NotPanics(t, func() {
		newLogger.Info("test")
	})
}

func TestLogger_WithError(t *testing.T) {
	logger := NewNop()
	err := assert.AnError

	newLogger := logger.WithError(err)

	assert.NotNil(t, newLogger)
	// Should not panic when logging
	assert.NotPanics(t, func() {
		newLogger.Info("test")
	})
}

func TestLogger_WithComponent(t *testing.T) {
	logger := NewNop()

	newLogger := logger.WithComponent("test-component")

	assert.NotNil(t, newLogger)
	// Should not panic when logging
	assert.NotPanics(t, func() {
		newLogger.Info("test")
	})
}

func TestLogger_WithRequestID(t *testing.T) {
	logger := NewNop()

	newLogger := logger.WithRequestID("req-123")

	assert.NotNil(t, newLogger)
	// Should not panic when logging
	assert.NotPanics(t, func() {
		newLogger.Info("test")
	})
}

func TestLogger_WithSpecID(t *testing.T) {
	logger := NewNop()

	newLogger := logger.WithSpecID("spec-456")

	assert.NotNil(t, newLogger)
	// Should not panic when logging
	assert.NotPanics(t, func() {
		newLogger.Info("test")
	})
}

func TestNewNop(t *testing.T) {
	logger := NewNop()

	assert.NotNil(t, logger)
	assert.NotNil(t, logger.Logger)

	// No-op logger should not panic on any operation
	assert.NotPanics(t, func() {
		logger.Info("test")
		logger.Error("test")
		logger.Debug("test")
		logger.Warn("test")
		_ = logger.Sync()
	})
}

func TestLogger_Sync(t *testing.T) {
	config := Config{
		Level:      "info",
		Format:     "json",
		OutputPath: "stdout",
	}

	logger, err := New(config)
	require.NoError(t, err)

	err = logger.Sync()
	// Sync may return an error on stdout/stderr, but should not panic
	// We just verify it doesn't panic
	assert.NotPanics(t, func() {
		_ = logger.Sync()
	})
}

func TestFieldHelpers(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		field := String("key", "value")
		assert.Equal(t, "key", field.Key)
		assert.Equal(t, zapcore.StringType, field.Type)
	})

	t.Run("Int", func(t *testing.T) {
		field := Int("key", 42)
		assert.Equal(t, "key", field.Key)
		assert.Equal(t, zapcore.Int64Type, field.Type)
	})

	t.Run("Int64", func(t *testing.T) {
		field := Int64("key", int64(42))
		assert.Equal(t, "key", field.Key)
		assert.Equal(t, zapcore.Int64Type, field.Type)
	})

	t.Run("Bool", func(t *testing.T) {
		field := Bool("key", true)
		assert.Equal(t, "key", field.Key)
		assert.Equal(t, zapcore.BoolType, field.Type)
	})

	t.Run("Duration", func(t *testing.T) {
		field := Duration("key", time.Second)
		assert.Equal(t, "key", field.Key)
		assert.Equal(t, zapcore.DurationType, field.Type)
	})

	t.Run("Any", func(t *testing.T) {
		field := Any("key", map[string]string{"nested": "value"})
		assert.Equal(t, "key", field.Key)
	})
}

func TestLogger_LogLevels(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	config := Config{
		Level:      "debug",
		Format:     "json",
		OutputPath: logFile,
	}

	logger, err := New(config)
	require.NoError(t, err)

	// Test all log levels
	assert.NotPanics(t, func() {
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")
	})

	_ = logger.Sync()

	// Verify log file has content
	content, err := os.ReadFile(logFile)
	require.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestLogger_ChainedWithMethods(t *testing.T) {
	logger := NewNop()

	// Test chaining multiple With methods
	newLogger := logger.
		WithComponent("test").
		WithRequestID("req-123").
		WithSpecID("spec-456").
		WithError(assert.AnError).
		WithFields(zap.String("custom", "field"))

	assert.NotNil(t, newLogger)

	// Should not panic
	assert.NotPanics(t, func() {
		newLogger.Info("test message")
	})
}

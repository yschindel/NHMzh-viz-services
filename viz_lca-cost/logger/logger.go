package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Fields is a map of key-value pairs for structured logging
type Fields map[string]interface{}

func init() {
	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	// Configure console writer with more detailed output
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("| %s ", i)
		},
	}

	// Create new global logger
	log.Logger = zerolog.New(consoleWriter).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()

	if os.Getenv("ENVIRONMENT") == "development" {
		log.Debug().Msg("Development environment detected, enabling debug logging")
	} else {
		if os.Getenv("ENVIRONMENT") == "" {
			log.Debug().Msg("ENVIRONMENT variable not set, using TRACE level")
		}
	}
}

// Logger wraps zerolog logger
type Logger struct {
	Logger zerolog.Logger
}

// New creates a new logger instance
func New() *Logger {
	return &Logger{
		Logger: log.With().Logger(),
	}
}

// WithContext creates a new logger with context fields
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// TODO: Add OpenTelemetry trace and span IDs when integrating
	return l
}

// WithFields creates a new logger with additional fields
func (l *Logger) WithFields(fields Fields) *Logger {
	newLogger := l.Logger.With().Fields(fields).Logger()
	return &Logger{Logger: newLogger}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.Logger.Debug().Msg(fmt.Sprintf(format, args...))
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.Logger.Info().Msg(fmt.Sprintf(format, args...))
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.Logger.Warn().Msg(fmt.Sprintf(format, args...))
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.Logger.Error().Msg(fmt.Sprintf(format, args...))
}

// Global logging functions
func Debug(format string, args ...interface{}) {
	log.Debug().Msg(fmt.Sprintf(format, args...))
}

func Info(format string, args ...interface{}) {
	log.Info().Msg(fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	log.Warn().Msg(fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	log.Error().Msg(fmt.Sprintf(format, args...))
}

func WithContext(ctx context.Context) *Logger {
	return &Logger{Logger: log.With().Logger()}
}

func WithFields(fields Fields) *Logger {
	return &Logger{Logger: log.With().Fields(fields).Logger()}
}

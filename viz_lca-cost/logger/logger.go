package logger

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

// Fields is a map of key-value pairs for structured logging
type Fields map[string]any

type Logger struct {
	logger *slog.Logger
}

func init() {
	// Create a new tinted logger with custom options
	w := os.Stderr
	opts := &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.RFC3339,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time from logs
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}

			// Colorize log levels
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				switch level {
				case slog.LevelDebug:
					return slog.Attr{
						Key:   slog.LevelKey,
						Value: slog.StringValue("\033[1;34mDEBUG\033[0m"), // Bold Blue
					}
				case slog.LevelInfo:
					return slog.Attr{
						Key:   slog.LevelKey,
						Value: slog.StringValue("\033[1;32mINFO\033[0m"), // Bold Green
					}
				case slog.LevelWarn:
					return slog.Attr{
						Key:   slog.LevelKey,
						Value: slog.StringValue("\033[1;33mWARN\033[0m"), // Bold Yellow
					}
				case slog.LevelError:
					return slog.Attr{
						Key:   slog.LevelKey,
						Value: slog.StringValue("\033[1;31mERROR\033[0m"), // Bold Red
					}
				}
			}

			// Colorize error values
			if err, ok := a.Value.Any().(error); ok {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue("\033[1;31m" + err.Error() + "\033[0m"), // Bold Red
				}
			}

			return a
		},
	}

	if os.Getenv("ENVIRONMENT") == "development" {
		opts.Level = slog.LevelDebug
	} else {
		if os.Getenv("ENVIRONMENT") == "" {
			opts.Level = slog.LevelDebug // Default to debug if no environment set
		} else {
			opts.Level = slog.LevelInfo // Production default
		}
	}

	// Set the default logger with tint handler
	slog.SetDefault(slog.New(tint.NewHandler(w, opts)))
}

// New creates a new logger instance
func New() *Logger {
	return &Logger{
		logger: slog.Default(),
	}
}

// WithContext creates a new logger with context fields
// This will be useful for OpenTelemetry integration
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// TODO: Add OpenTelemetry trace and span IDs when integrating
	// Example future implementation:
	// if span := trace.SpanFromContext(ctx); span != nil {
	//     traceID := span.SpanContext().TraceID().String()
	//     spanID := span.SpanContext().SpanID().String()
	//     return l.WithFields(Fields{
	//         "trace_id": traceID,
	//         "span_id":  spanID,
	//     })
	// }
	return l
}

// WithFields creates a new logger with additional fields
func (l *Logger) WithFields(fields Fields) *Logger {
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}
	return &Logger{logger: l.logger.With(attrs...)}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.logger.Debug(format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.logger.Info(format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.logger.Warn(format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.logger.Error(format, args...)
}

// Global logging functions
func Debug(format string, args ...interface{}) {
	slog.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	slog.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	slog.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	slog.Error(format, args...)
}

func WithContext(ctx context.Context) *Logger {
	return New().WithContext(ctx)
}

func WithFields(fields Fields) *Logger {
	return New().WithFields(fields)
}

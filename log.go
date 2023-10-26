package zeus_log

import (
	"context"
	"github.com/oyogames2023/zeus-log/internal/env"
	"os"
)

type (
	loggerKey struct{}
)

var traceEnabled = traceEnableFromEnv()

// traceEnableFromEnv checks whether trace is enabled by reading from environment.
// Close trace if empty or zero, open trace if not zero, default as closed.
func traceEnableFromEnv() bool {
	switch os.Getenv(env.EnabledTraceLog) {
	case "":
		fallthrough
	case "0":
		return false
	default:
		return true
	}
}

// SetLevel sets log level for different output which may be "trace", "debug", "info".
func SetLevel(output string, level Level) {
	GetDefaultLogger().SetLevel(output, level)
}

// GetLevel gets log level for different output.
func GetLevel(output string) Level {
	return GetDefaultLogger().GetLevel(output)
}

func With(args ...any) Logger {
	if ol, ok := GetDefaultLogger().(OptionLogger); ok {
		return ol.WithOptions(WithAdditionalCallerSkip(-1)).With(args...)
	}
	return GetDefaultLogger().With(args...)
}

func WithFields(fields ...Field) Logger {
	if ol, ok := GetDefaultLogger().(OptionLogger); ok {
		return ol.WithOptions(WithAdditionalCallerSkip(-1)).WithFields(fields...)
	}
	return GetDefaultLogger().WithFields(fields...)
}

func WithFieldsContext(ctx context.Context, fields ...Field) Logger {
	var logger Logger
	var ok bool
	logger, ok = ctx.Value(loggerKey{}).(Logger)
	if !ok {
		return WithFields(fields...)
	}
	var ol OptionLogger
	if ol, ok = logger.(OptionLogger); ok {
		return ol.WithOptions(WithAdditionalCallerSkip(-1)).WithFields(fields...)
	}
	return logger.WithFields(fields...)
}

func WithContext(ctx context.Context, args ...any) Logger {
	var logger Logger
	var ok bool
	logger, ok = ctx.Value(loggerKey{}).(Logger)
	if !ok {
		return With(args...)
	}
	var ol OptionLogger
	if ol, ok = logger.(OptionLogger); ok {
		return ol.WithOptions(WithAdditionalCallerSkip(-1)).With(args...)
	}
	return logger.With(args...)
}

// Trace logs to TRACE log. Arguments are handled in the manner of fmt.Print.
func Trace(args ...any) {
	if traceEnabled {
		GetDefaultLogger().Trace(args...)
	}
}

// Tracef logs to TRACE log. Arguments are handled in the manner of fmt.Printf.
func Tracef(format string, args ...any) {
	if traceEnabled {
		GetDefaultLogger().Tracef(format, args...)
	}
}

// Traceln logs to TRACE log. Arguments are handled in the manner of fmt.Printf.
func Traceln(args ...any) {
	if traceEnabled {
		GetDefaultLogger().Traceln(args...)
	}
}

// Debug logs to DEBUG log. Arguments are handled in the manner of fmt.Print.
func Debug(args ...interface{}) {
	GetDefaultLogger().Debug(args...)
}

// Debugf logs to DEBUG log. Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, args ...interface{}) {
	GetDefaultLogger().Debugf(format, args...)
}

// Debugf logs to DEBUG log. Arguments are handled in the manner of fmt.Println.
func Debugln(args ...any) {
	GetDefaultLogger().Debugln(args...)
}

// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
func Info(args ...interface{}) {
	GetDefaultLogger().Info(args...)
}

// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func Infof(format string, args ...interface{}) {
	GetDefaultLogger().Infof(format, args...)
}

// Warn logs to WARNING log. Arguments are handled in the manner of fmt.Print.
func Warn(args ...interface{}) {
	GetDefaultLogger().Warn(args...)
}

// Warnf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func Warnf(format string, args ...interface{}) {
	GetDefaultLogger().Warnf(format, args...)
}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
func Error(args ...interface{}) {
	GetDefaultLogger().Error(args...)
}

// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, args ...interface{}) {
	GetDefaultLogger().Errorf(format, args...)
}

// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
// All Fatal logs will exit by calling os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func Fatal(args ...interface{}) {
	GetDefaultLogger().Fatal(args...)
}

// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func Fatalf(format string, args ...interface{}) {
	GetDefaultLogger().Fatalf(format, args...)
}

// TraceContext logs to TRACE log. Arguments are handled in the manner of fmt.Print.
func TraceContext(ctx context.Context, args ...any) {
	if !traceEnabled {
		return
	}
	if l, ok := ctx.Value(loggerKey{}).(Logger); ok {
		l.Trace(args...)
		return
	}
	GetDefaultLogger().Trace(args...)
}

// TraceContextf logs to TRACE log. Arguments are handled in the manner of fmt.Printf.
func TraceContextf(ctx context.Context, format string, args ...any) {
	if !traceEnabled {
		return
	}
	if l, ok := ctx.Value(loggerKey{}).(Logger); ok {
		l.Tracef(format, args...)
		return
	}
	GetDefaultLogger().Tracef(format, args...)
}

// TraceContextln logs to TRACE log. Arguments are handled in the manner of fmt.Printf.
func TraceContextln(ctx context.Context, args ...any) {
	if !traceEnabled {
		return
	}
	if l, ok := ctx.Value(loggerKey{}).(Logger); ok {
		l.Traceln(args...)
		return
	}
	GetDefaultLogger().Traceln(args...)
}

package zeus_log

import "io"

type Level int

const (
	// LevelOff disables log output.
	LevelOff Level = iota
	LevelTrace
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

func (lv *Level) String() string {
	return LevelStrings[*lv]
}

var (
	// LevelStrings is the map from log level to its string representation.
	LevelStrings = map[Level]string{
		LevelTrace: "trace",
		LevelDebug: "debug",
		LevelInfo:  "info",
		LevelWarn:  "warn",
		LevelError: "error",
		LevelFatal: "fatal",
		LevelPanic: "panic",
	}
	// LevelNames is the map from string to log level.
	LevelNames = map[string]Level{
		"trace": LevelTrace,
		"debug": LevelDebug,
		"info":  LevelInfo,
		"warn":  LevelWarn,
		"error": LevelError,
		"fatal": LevelFatal,
		"panic": LevelPanic,
	}
)

// LoggerOptions is the log options.
type LoggerOptions struct {
	LogLevel Level
	Pattern  string
	Writer   io.Writer
}

// LoggerOption modifies the LoggerOptions.
type LoggerOption func(options *LoggerOptions)

// Field is the user defined log field.
type Field struct {
	Key   string
	Value any
}

// Logger provides an abstract definition for logging functionality.
type Logger interface {
	// Trace logs the provided arguments at [LevelTrace]. Spaces are added between
	// arguments when neither is a string.
	Trace(args ...any)
	// Tracef formats the message according to the format specifier and logs it at [LevelTrace].
	Tracef(args ...any)
	// Traceln logs a message at [LevelTrace]. Spaces are always added between arguments.
	Traceln(args ...any)
	// Debug logs the provided arguments at [LevelDebug]. Spaces are added between
	// arguments when neither is a string.
	Debug(args ...any)
	// Debugf formats the message according to the format specifier and logs it at [LevelDebug].
	Debugf(args ...any)
	// Debugln logs a message at [LevelDebug]. Spaces are always added between arguments.
	Debugln(args ...any)
	// Info logs the provided arguments at [LevelInfo]. Spaces are added between
	// arguments when neither is a string.
	Info(args ...any)
	// Infof formats the message according to the format specifier and logs it at [LevelInfo].
	Infof(args ...any)
	// Infoln logs a message at [LevelInfo]. Spaces are always added between arguments.
	Infoln(args ...any)
	// Warn logs the provided arguments at [LevelWarn]. Spaces are added between
	// arguments when neither is a string.
	Warn(args ...any)
	// Warnf formats the message according to the format specifier and logs it at [LevelWarn].
	Warnf(args ...any)
	// Warnln logs a message at [LevelWarn]. Spaces are always added between arguments.
	Warnln(args ...any)
	// Error logs the provided arguments at [LevelError]. Spaces are added between
	// arguments when neither is a string.
	Error(args ...any)
	// Errorf formats the message according to the format specifier and logs it at [LevelError].
	Errorf(args ...any)
	// Errorln logs a message at [LevelError]. Spaces are always added between arguments.
	Errorln(args ...any)
	// Fatal logs the provided arguments at [LevelFatal]. Spaces are added between
	// arguments when neither is a string.
	Fatal(args ...any)
	// Fatalf formats the message according to the format specifier and logs it at [LevelFatal].
	Fatalf(args ...any)
	// Fatalln logs a message at [LevelFatal]. Spaces are always added between arguments.
	Fatalln(args ...any)
	// Panic logs the provided arguments at [LevelPanic]. Spaces are added between
	// arguments when neither is a string.
	Panic(args ...any)
	// Panicf formats the message according to the format specifier and logs it at [LevelPanic].
	Panicf(args ...any)
	// Panicln logs a message at [LevelPanic]. Spaces are always added between arguments.
	Panicln(args ...any)

	// Sync calls the underlying Core's Sync method, flushing any buffer log entries.
	// Applications should take care to call Sync before exiting.
	Sync() error

	// SetLevel sets the output log level.
	SetLevel(output string, level Level)

	// GetLevel gets the output log level.
	GetLevel(output string) Level

	// With returns a new logger with key/value paris.
	With(args ...any) Logger

	// WithFields returns a new logger with `fields` set.
	WithFields(fields ...Field) Logger
}

type OptionLogger interface {
	WithOptions(opts ...Option) Logger
}

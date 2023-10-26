package zeus_log

import (
	"github.com/oyogames2023/zeus-log/plugin"
	"sync"
)

const (
	pluginType        = "log"
	defaultLoggerName = "default"
)

var (
	DefaultLogger  Logger
	DefaultFactory plugin.Factory

	mu      sync.RWMutex
	loggers = make(map[string]Logger)
)

// Register registers Logger. It supports multiple Logger implementation.
func Register(name string, logger Logger) {
	mu.Lock()
	defer mu.Unlock()
	if logger == nil {
		panic("log: Register logger is nil")
	}
	if _, dup := loggers[name]; dup && name != defaultLoggerName {
		panic("log: Register called twice for logger name " + name)
	}
	loggers[name] = logger
	if name == defaultLoggerName {
		DefaultLogger = logger
	}
}

// GetDefaultLogger gets the default Logger.
// To configure it, set key in configuration file to default.
// The console output is the default value.
func GetDefaultLogger() Logger {
	mu.RLock()
	l := DefaultLogger
	mu.RUnlock()
	return l
}

// SetLogger sets the default Logger.
func SetLogger(logger Logger) {
	mu.Lock()
	DefaultLogger = logger
	mu.Unlock()
}

// Get returns the Logger implementation by log name.
// log.Debug use DefaultLogger to print logs. You may also use log.Get("name").Debug.
func Get(name string) Logger {
	mu.RLock()
	l := loggers[name]
	mu.RUnlock()
	return l
}

// Sync syncs all registered loggers.
func Sync() {
	mu.RLock()
	defer mu.RUnlock()
	for _, logger := range loggers {
		_ = logger.Sync()
	}
}

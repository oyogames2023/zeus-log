package zeus_log

import (
	"github.com/oyogames2023/zeus-log/plugin"
)

var (
	// DefaultConsoleWriterFactory is the default console output implementation.
	DefaultConsoleWriterFactory plugin.Factory
	// DefaultFileWriterFactory is the default file output implementation.
	DefaultFileWriterFactory plugin.Factory

	writers = make(map[string]plugin.Factory)
)

// RegisterWriter registers log output writer. Writer may have multiple implementations.
func RegisterWriter(name string, writer plugin.Factory) {
	writers[name] = writer
}

// GetWriter gets log output writer, returns nil if not exist.
func GetWriter(name string) plugin.Factory {
	return writers[name]
}

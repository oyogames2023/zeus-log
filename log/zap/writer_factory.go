package zap

import (
	xlog "github.com/oyogames2023/zeus-log"
	ec "github.com/oyogames2023/zeus-log/errorcode"
	"github.com/oyogames2023/zeus-log/plugin"
	"path/filepath"
)

// ConsoleWriterFactory is the console writer instance.
type ConsoleWriterFactory struct {
}

// Type returns the log plugin type.
func (f *ConsoleWriterFactory) Type() string {
	return pluginType
}

// Setup starts, loads and registers console output writer.
func (f *ConsoleWriterFactory) Setup(name string, dec plugin.Decoder) error {
	if dec == nil {
		return ec.ErrInvalidWriterDecoderObject
	}
	decoder, ok := dec.(*Decoder)
	if !ok {
		return ec.ErrInvalidWriterDecoderType
	}
	cfg := &xlog.OutputConfig{}
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}
	decoder.Core, decoder.ZapLevel = newConsoleCore(cfg)
	return nil
}

// FileWriterFactory is the file writer instance Factory.
type FileWriterFactory struct {
}

// Type returns log file type.
func (f *FileWriterFactory) Type() string {
	return pluginType
}

// Setup starts, loads and register file output writer.
func (f *FileWriterFactory) Setup(name string, dec plugin.Decoder) error {
	if dec == nil {
		return ec.ErrInvalidWriterDecoderObject
	}
	decoder, ok := dec.(*Decoder)
	if !ok {
		return ec.ErrInvalidWriterDecoderType
	}
	if err := f.setupConfig(decoder); err != nil {
		return err
	}
	return nil
}

func (f *FileWriterFactory) setupConfig(decoder *Decoder) error {
	cfg := &xlog.OutputConfig{}
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}
	if cfg.WriterConfig.LogPath != "" {
		cfg.WriterConfig.FileName = filepath.Join(cfg.WriterConfig.LogPath,
			cfg.WriterConfig.FileName)
	}
	if cfg.WriterConfig.RollType == "" {
		cfg.WriterConfig.RollType = xlog.GetRollingType(xlog.RollingBySize)
	}

	core, level, err := newFileCore(cfg)
	if err != nil {
		return err
	}
	decoder.Core, decoder.ZapLevel = core, level
	return nil
}

package zap

import (
	"errors"
	"fmt"
	xlog "github.com/oyogames2023/zeus-log"
	"github.com/oyogames2023/zeus-log/plugin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	pluginType = "log"
)

// Decoder decodes the log.
type Decoder struct {
	OutputConfig *xlog.OutputConfig
	Core         zapcore.Core
	ZapLevel     zap.AtomicLevel
}

// Decode decodes writer configuration, copy one.
func (d *Decoder) Decode(cfg interface{}) error {
	output, ok := cfg.(**xlog.OutputConfig)
	if !ok {
		return fmt.Errorf("decoder config type:%T invalid, not **OutputConfig", cfg)
	}
	*output = d.OutputConfig
	return nil
}

// Factory is the log plugin factory.
// When server start, the configuration is feed to Factory to generate a log instance.
type Factory struct{}

// Type returns the log plugin type.
func (f *Factory) Type() string {
	return pluginType
}

// Setup starts, load and register logs.
func (f *Factory) Setup(name string, dec plugin.Decoder) error {
	if dec == nil {
		return errors.New("log config decoder empty")
	}
	cfg, callerSkip, err := f.setupConfig(dec)
	if err != nil {
		return err
	}
	logger := NewZapLogWithCallerSkip(cfg, callerSkip)
	if logger == nil {
		return errors.New("new zap logger fail")
	}
	xlog.Register(name, logger)
	return nil
}

func (f *Factory) setupConfig(configDec plugin.Decoder) (xlog.Config, int, error) {
	cfg := xlog.Config{}
	if err := configDec.Decode(&cfg); err != nil {
		return nil, 0, err
	}
	if len(cfg) == 0 {
		return nil, 0, errors.New("log config output empty")
	}

	// If caller skip is not configured, use 2 as default.
	callerSkip := 2
	for i := 0; i < len(cfg); i++ {
		if cfg[i].CallerSkip != 0 {
			callerSkip = cfg[i].CallerSkip
		}
	}
	return cfg, callerSkip, nil
}

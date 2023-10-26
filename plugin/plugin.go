package plugin

var (
	plugins = make(map[string]map[string]Factory)
)

type Factory interface {
	// Type returns type of the plugin, i.e. selector, log, config, tracing.
	Type() string
	// Setup loads plugin by configuration. The data structure of the configuration
	// of the plugin needs to be defined in advance.
	Setup(name string, decoder Decoder) error
}

type Decoder interface {
	// Decode the input parameter is the custom configuration of the plugin.
	Decode(cfg any) error
}

// Register registers a plugin factory. Name of the plugin should be specified.
// It is supported to register instances which are the same implementation of
// plugin Factory, but use different configuration.
func Register(name string, factory Factory) {
	factories, ok := plugins[factory.Type()]
	if !ok {
		factories = make(map[string]Factory)
		plugins[factory.Type()] = factories
	}
	factories[name] = factory
}

// Get returns a plugin Factory by its type and name.
func Get(factoryType string, name string) Factory {
	if factories, ok := plugins[factoryType]; ok {
		return factories[name]
	}
	return nil
}

package zeus_log

// Option modifies the options of OptionLogger.
type Option func(*Options)

type Options struct {
	Skip int
}

// WithAdditionalCallerSkip adds additional caller skip.
func WithAdditionalCallerSkip(skip int) Option {
	return func(o *Options) {
		o.Skip = skip
	}
}

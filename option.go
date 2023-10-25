package zeus_log

// Option modifies the options of OptionLogger.
type Option func(*options)

type options struct {
	skip int
}

// WithAdditionalCallerSkip adds additional caller skip.
func WithAdditionalCallerSkip(skip int) Option {
	return func(o *options) {
		o.skip = skip
	}
}

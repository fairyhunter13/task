package task

// OptionFunc specifies the optional function for this package.
type OptionFunc func(*Option)

// Option specifies the option to be used in this package.
type Option struct {
	// UsePanicHandler defaults to false
	UsePanicHandler bool
}

func (o *Option) Assign(opts ...OptionFunc) *Option {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// WithPanicHandler adds the option to toggle the panic handler on and off.
func WithPanicHandler(confirm bool) OptionFunc {
	return func(o *Option) {
		o.UsePanicHandler = confirm
	}
}

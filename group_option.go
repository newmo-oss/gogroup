package gogroup

import "github.com/newmo-oss/gogroup/internal"

// Option is an option to configure [(*Group).Run] behaviour.
type Option interface {
	apply(*internal.Options)
}

type optionFunc func(*internal.Options)

func (f optionFunc) apply(opts *internal.Options) {
	f(opts)
}

var _ Option = optionFunc(nil)

// WithLimit sets the maximum number of goroutines used to run the functions. Ignored if n < 1.
func WithLimit(n int) Option {
	return optionFunc(func(opts *internal.Options) {
		opts.Limit = n
	})
}

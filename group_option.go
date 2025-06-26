package gogroup

import "github.com/newmo-oss/gogroup/internal"

// GroupOption is an option to configure [(*Group).Run] behaviour.
type GroupOption interface {
	apply(*internal.GroupOptions)
}

// WithLimit sets the maximum number of goroutines used to run the functions. Ignored if n < 1.
func WithLimit(n int) GroupOption {
	return withLimitOption{limit: n}
}

type withLimitOption struct {
	limit int
}

func (wo withLimitOption) apply(o *internal.GroupOptions) {
	o.Limit = wo.limit
}

package internal

import (
	"context"

	"github.com/sourcegraph/conc/panics"
	"github.com/sourcegraph/conc/pool"
)

var Start = DefaultStart

// GroupOption is an option to configure [(*Group).Run] behaviour.
type Option func(*GroupOptions) error

type GroupOptions struct {
	Limit int
}

func DefaultStart(ctx context.Context, funcs []func(context.Context) error, opts ...Option) func() error {
	var o GroupOptions
	for _, opt := range opts {
		opt(&o)
	}

	p := pool.New().WithContext(ctx).WithCancelOnError()

	if o.Limit > 0 {
		p = p.WithMaxGoroutines(o.Limit)
	}

	for _, f := range funcs {
		p.Go(func(ctx context.Context) (rerr error) {
			if r := panics.Try(func() { rerr = f(ctx) }); r != nil {
				return r.AsError()
			}
			return rerr
		})
	}
	return p.Wait // return method value
}

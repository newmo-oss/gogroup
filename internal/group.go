package internal

import (
	"context"

	"github.com/sourcegraph/conc/panics"
	"github.com/sourcegraph/conc/pool"
)

var Start = DefaultStart

type GroupOptions struct {
	Limit int
}

func DefaultStart(ctx context.Context, funcs []func(context.Context) error, option GroupOptions) func() error {
	p := pool.New().WithContext(ctx).WithCancelOnError()

	if option.Limit > 0 {
		p = p.WithMaxGoroutines(option.Limit)
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

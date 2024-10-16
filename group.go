package gogroup

import (
	"context"
	"slices"
	"sync"

	"github.com/sourcegraph/conc/panics"
	"github.com/sourcegraph/conc/pool"
)

// Group is a group of goroutines used to run functions concurrently.
// Functions are registered with [(*Group).Add].
// [(*Group).Run] calls each functions in different goroutines and waits for the functions to finish.
// If a panic occurs in a function, Run recovers the panic and returns it as an error.
// If any function returns non nil error, the context is canceled.
type Group struct {
	mu    sync.Mutex
	funcs []func(context.Context) error
}

// Add adds a function to the group.
func (g *Group) Add(f func(context.Context) error) {
	g.mu.Lock()
	g.funcs = append(g.funcs, f)
	g.mu.Unlock()
}

func (g *Group) start(ctx context.Context) func() error {
	p := pool.New().WithContext(ctx).WithCancelOnError()

	g.mu.Lock()
	funcs := slices.Clone(g.funcs)
	g.mu.Unlock()

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

// Run calls all registered functions in different goroutines.
func (g *Group) Run(ctx context.Context) error {
	return g.start(ctx)()
}

// Start calls the function in new goroutine and returns a wait function.
// When the wait function is called, it waits for the goroutine and returns the returned value of the function.
// If a panic occurs in the function, the wait function recovers the panic and returns it as an error.
func Start(ctx context.Context, f func(context.Context) error) (wait func() error) {
	var g Group
	g.Add(f)
	return g.start(ctx)
}

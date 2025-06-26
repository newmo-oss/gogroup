package gogrouptest

import (
	"context"
	"sync"
	"testing"

	"github.com/newmo-oss/testid"
	"github.com/sourcegraph/conc/panics"
	"github.com/sourcegraph/conc/pool"

	"github.com/newmo-oss/gogroup/internal"
)

var withoutParallels sync.Map

func init() {
	if testing.Testing() {
		internal.Start = startForTest
	}
}

// WithoutParallel sets the parallel mode of [gogroup.Group] to off.
// It makes (*gogroup.Group).Run will execute synchronously.
// The parallel mode is set each test id which get from [testid.FromContext].
// If any test id cannot obtain from the context, the test will be fail with t.Fatal.
// The mode will be remove by t.Cleanup.
func WithoutParallel(t testing.TB, ctx context.Context) {
	tid, ok := testid.FromContext(ctx)
	if !ok {
		t.Fatal("failed to get test ID from the context")
	}

	t.Cleanup(func() {
		withoutParallels.Delete(tid)
	})

	withoutParallels.Store(tid, true)
}

func noParallel(ctx context.Context) bool {
	tid, ok := testid.FromContext(ctx)
	if !ok {
		return false
	}

	v, ok := withoutParallels.Load(tid)
	if !ok {
		return false
	}

	noparallel, ok := v.(bool)
	if !ok {
		return false
	}

	return noparallel
}

func startForTest(ctx context.Context, funcs []func(context.Context) error, _ ...internal.Option) func() error {
	if !noParallel(ctx) {
		return internal.DefaultStart(ctx, funcs)
	}

	p := pool.New().WithContext(ctx).WithCancelOnError()
	doneCh := make([]chan struct{}, len(funcs))

	for i, f := range funcs {
		doneCh[i] = make(chan struct{})
		p.Go(func(ctx context.Context) (rerr error) {
			defer func() {
				close(doneCh[i])
			}()

			// wait before function call
			if i > 0 {
				<-doneCh[i-1]
			}

			if r := panics.Try(func() { rerr = f(ctx) }); r != nil {
				return r.AsError()
			}

			return rerr
		})
	}

	return p.Wait // return method value
}

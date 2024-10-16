package gogroup_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/newmo-oss/gogroup"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestGroup(t *testing.T) {
	t.Parallel()

	var (
		funcOK     = func(ctx context.Context) error { return nil }
		funcErr    = func(ctx context.Context) error { return errors.New("error") }
		funcPanic  = func(ctx context.Context) error { panic("panic") }
		funcCancel = func(ctx context.Context) error { doCancel(ctx); return nil }
	)

	type funcs []func(ctx context.Context) error

	cases := map[string]struct {
		funcs            funcs
		wantCanceledFunc int
		wantErr          bool
	}{
		// name          funcs                          canceled  error
		"all-ok":        {funcs{funcOK, funcOK, funcOK}, -1, false},
		"first-error":   {funcs{funcErr, funcOK, funcOK}, 0, true},
		"second-error":  {funcs{funcOK, funcErr, funcOK}, 1, true},
		"first-panic":   {funcs{funcPanic, funcOK, funcOK}, 0, true},
		"second-panic":  {funcs{funcOK, funcPanic, funcOK}, 1, true},
		"first-cancel":  {funcs{funcCancel, funcOK, funcOK}, 0, false},
		"second-cancel": {funcs{funcOK, funcCancel, funcOK}, 1, false},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				doneCh   = make([]chan struct{}, len(tt.funcs))
				canceled atomic.Bool
			)
			var g gogroup.Group
			for i, f := range tt.funcs {
				doneCh[i] = make(chan struct{})
				g.Add(func(ctx context.Context) error {
					defer func() {
						close(doneCh[i])
					}()

					// wait before function call
					if i > 0 {
						<-doneCh[i-1]
					}

					if tt.wantCanceledFunc >= 0 && tt.wantCanceledFunc < i {
						select {
						case <-time.After(1 * time.Second):
						case <-ctx.Done():
							canceled.Store(true)
						}
					}

					return f(ctx)
				})
			}

			ctx := withCancel(context.Background())
			err := g.Run(ctx)

			switch {
			case tt.wantErr && err == nil:
				t.Error("expected error did not occur")
			case !tt.wantErr && err != nil:
				t.Error("unexpected error:", err)
			}

			switch {
			case tt.wantCanceledFunc >= 0 && !canceled.Load():
				t.Error("expected cancel did not occur")
			case tt.wantCanceledFunc < 0 && canceled.Load():
				t.Error("unexpected cancel")
			}
		})
	}
}

func TestStart(t *testing.T) {
	t.Parallel()

	var (
		funcOK     = func(ctx context.Context) error { return nil }
		funcErr    = func(ctx context.Context) error { return errors.New("error") }
		funcPanic  = func(ctx context.Context) error { panic("panic") }
		funcCancel = func(ctx context.Context) error { doCancel(ctx); return nil }
	)

	cases := map[string]struct {
		fun     func(context.Context) error
		wantErr bool
	}{
		// name    func    error
		"ok":     {funcOK, false},
		"error":  {funcErr, true},
		"cancel": {funcCancel, false},
		"panic":  {funcPanic, true},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := withCancel(context.Background())
			wait := gogroup.Start(ctx, tt.fun)

			err := wait()

			switch {
			case tt.wantErr && err == nil:
				t.Error("expected error did not occur")
			case !tt.wantErr && err != nil:
				t.Error("unexpected error:", err)
			}
		})
	}
}

type cancelCtxKey struct{}

func withCancel(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	return context.WithValue(ctx, cancelCtxKey{}, cancel)
}

func doCancel(ctx context.Context) {
	cancel, _ := ctx.Value(cancelCtxKey{}).(context.CancelFunc)
	if cancel != nil {
		cancel()
	}
}

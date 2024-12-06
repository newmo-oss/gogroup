package gogrouptest_test

import (
	"context"
	"runtime"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/newmo-oss/testid"
	"github.com/newmo-oss/gotestingmock"

	"github.com/newmo-oss/gogroup"
	"github.com/newmo-oss/gogroup/gogrouptest"
)

func TestWithoutParallel(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		noParallel bool
	}{
		"no parallel": {true},
		"parallel":    {false},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var cleanedup bool
			t.Cleanup(func() {
				if tt.noParallel && !cleanedup {
					t.Error("t.Cleanup must be called")
				}
			})

			var want, got []int
			var mu sync.Mutex
			var g gogroup.Group
			for i := range 100 {
				want = append(want, i)
				g.Add(func(ctx context.Context) error {
					if i%2 == 0 {
						runtime.Gosched()
					}
					mu.Lock()
					got = append(got, i)
					mu.Unlock()
					return nil
				})
			}

			tid := t.Name() + "/" + uuid.NewString()
			ctx := testid.WithValue(context.Background(), tid)

			if tt.noParallel {
				tb := &gotestingmock.TB{
					TB: t,
					CleanupFunc: func(f func()) {
						cleanedup = true
						t.Cleanup(f)
					},
				}
				gogrouptest.WithoutParallel(tb, ctx)
			}

			if err := g.Run(ctx); err != nil {
				t.Fatal("unexpected error:", err)
			}

			diff := cmp.Diff(got, want)
			switch {
			case tt.noParallel && diff != "":
				t.Error("executing order does not match:", diff)
			case !tt.noParallel && diff == "":
				t.Error("executing order may be randodm")
			}
		})
	}
}

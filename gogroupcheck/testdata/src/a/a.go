package a

import (
	"sync" // ok

	_ "github.com/sourcegraph/conc"      // want `github\.com/sourcegraph/conc is disallowed; use github\.com/newmo-oss/gogroup instead`
	_ "github.com/sourcegraph/conc/pool" // want `github\.com/sourcegraph/conc/pool is disallowed; use github\.com/newmo-oss/gogroup instead`
	_ "golang.org/x/sync/errgroup"       // want `golang\.org/x/sync/errgroup is disallowed; use github\.com/newmo-oss/gogroup instead`
	_ "golang.org/x/sync/singleflight"   // ok

	_ "github.com/newmo-oss/gogroup" // ok
)

func f() {
	var _ sync.WaitGroup // want `sync\.WaitGroup is disallowed; use github\.com/newmo-oss/gogroup\.Group instead`
}

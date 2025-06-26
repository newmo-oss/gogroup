package gogroup

import "github.com/newmo-oss/gogroup/internal"

// WithLimit sets the maximum number of goroutines used to run the functions. Ignored if n < 1.
func WithLimit(n int) internal.Option {
	return func(o *internal.GroupOptions) error {
		o.Limit = n
		return nil
	}
}

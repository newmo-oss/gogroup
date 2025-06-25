package gogroup

// GroupOption is an option to configure [(*Group).Run] behaviour.
type Option func(*groupOptions) error

type groupOptions struct {
	limit int
}

// WithLimit sets the maximum number of goroutines used to run the functions. Ignored if n < 1.
func WithLimit(n int) Option {
	return func(o *groupOptions) error {
		o.limit = n
		return nil
	}
}

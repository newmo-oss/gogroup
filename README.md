# gogroup [![Go Reference](https://pkg.go.dev/badge/github.com/newmo-oss/gogroupo.svg)](https://pkg.go.dev/github.com/newmo-oss/gogroup)[![Go Report Card](https://goreportcard.com/badge/github.com/newmo-oss/gogroup)](https://goreportcard.com/report/github.com/newmo-oss/gogroup)

gogroup provides a group of goroutines used to run functions concurrently.

## Usage

```go
var g gogroup.Group

g.Add(func(ctx context.Context) error {
	// do something
	return nil
})

g.Add(func(ctx context.Context) error {
	// convert panic as an error and cancel the context
	panic("panic")
})

if err := g.Run(ctx); err != nil {
	return err
}
```

## License
MIT

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

## gogroupcheck

**gogroupcheck** is an analyzer that reports any use of `sync.WaitGroup`, `golang.org/x/sync/errgroup.Group`,  
`github.com/sourcegraph/conc.WaitGroup`, and `github.com/sourcegraph/conc/pool` (and its subpackages), keeping your
codebase on a single, consistent concurrency library.

### Install

```sh
go install github.com/newmo-oss/gogroupcheck/cmd/gogroupcheck@latest
```

### Usage

```sh
go vet -vettool=$(which gogroupcheck) ./...
```

## License
MIT

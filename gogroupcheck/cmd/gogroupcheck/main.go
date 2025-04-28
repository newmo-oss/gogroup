package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"

	"github.com/newmo-oss/gogroup/gogroupcheck"
)

func main() {
	unitchecker.Main(gogroupcheck.Analyzer)
}

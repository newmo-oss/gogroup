package gogroupcheck

import (
	"slices"
	"strconv"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
)

const doc = "gogroupcheck is an analyzer that reports any use of sync.WaitGroup, golang.org/x/sync/errgroup.Group, github.com/sourcegraph/conc.WaitGroup, and github.com/sourcegraph/conc/pool to keep them out of your codebase."

var Analyzer = &analysis.Analyzer{
	Name: "gogroupcheck",
	Doc:  doc,
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	checkImportPath(pass)
	checkObject(pass)
	return nil, nil
}

func checkImportPath(pass *analysis.Pass) {
	checklist := []string{
		"golang.org/x/sync/errgroup",
		"github.com/sourcegraph/conc",
		"github.com/sourcegraph/conc/pool",
	}

	for _, files := range pass.Files {
		for _, spec := range files.Imports {
			path, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				// skip
				continue
			}

			if slices.Contains(checklist, path) {
				pass.Reportf(spec.Pos(), "%s is disallowed; use github.com/newmo-oss/gogroup instead", path)
			}
		}
	}
}

func checkObject(pass *analysis.Pass) {

	checklist := []struct {
		pkgpath, name string
	}{
		{"sync", "WaitGroup"},
	}

	for i := range checklist {
		obj := analysisutil.ObjectOf(pass, checklist[i].pkgpath, checklist[i].name)
		if obj == nil {
			continue
		}

		for id, use := range pass.TypesInfo.Uses {
			if use == obj {
				pass.Reportf(id.Pos(), "%s.%s is disallowed; use github.com/newmo-oss/gogroup.Group instead", checklist[i].pkgpath, checklist[i].name)
			}
		}
	}
}

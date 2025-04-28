package gogroupcheck_test

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"github.com/newmo-oss/gogroup/gogroupcheck"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	modfile := testutil.ModFile(t, ".", nil)
	testdata := testutil.WithModules(t, analysistest.TestData(), modfile)
	analysistest.Run(t, testdata, gogroupcheck.Analyzer, "a")
}

package main

import (
	"github.com/f044fs3t5w3f/metrics/pkg/analyzers/os_exit"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck/st1001"
)

func main() {
	analyzers := make([]*analysis.Analyzer, 0, len(staticcheck.Analyzers))
	for _, analyzer := range staticcheck.Analyzers {
		analyzers = append(analyzers, analyzer.Analyzer)
	}
	analyzers = append(analyzers, st1001.Analyzer)
	analyzers = append(analyzers, printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		os_exit.Analyzer,
		nilness.Analyzer,
		lostcancel.Analyzer,
		ifaceassert.Analyzer,
	)

	multichecker.Main(analyzers...)
}

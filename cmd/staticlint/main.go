// Package main implements a multichecker that runs a set of static analysis
// analyzers for Go source code.
//
// Usage:
//
//	staticlint ./...
//
// The multichecker includes the following analyzers:
//
//   - staticcheck.Analyzers reports a wide range of bugs, performance issues,
//     and style problems detected by the staticcheck.io project.
//
//   - st1001.Analyzer reports cases where defer is used with a method value
//     instead of a method call, which may lead to unintended behavior.
//
//   - printf.Analyzer checks consistency between format strings and arguments
//     in functions such as fmt.Printf, fmt.Sprintf, and related helpers.
//
//   - structtag.Analyzer validates the syntax and semantics of struct field tags
//     (for example, json:"field").
//
//   - os_exit.Analyzer reports calls to os.Exit in the main function of the
//     main package.
//
//   - nilness.Analyzer detects potential nil pointer dereferences.
//
//   - lostcancel.Analyzer reports contexts created by context.WithCancel,
//     context.WithTimeout, or context.WithDeadline whose cancellation function
//     is not called, which may lead to goroutine leaks.
//
//   - ifaceassert.Analyzer reports unsafe type assertions from interfaces
//     without checking the assertion result.
//     это всё написано с помощью LLM ^
//
//   - restrictpkg.RestrictPackageAnalyzer restrict using of certain packages.
//     in this checker it preset for restict reflect package
package main

import (
	"github.com/cybozu-go/golang-custom-analyzer/pkg/restrictpkg"
	"github.com/f044fs3t5w3f/metrics/pkg/analyzers/osexit"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
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
	f := restrictpkg.RestrictPackageAnalyzer.Flags.Lookup("packages")
	f.Value.Set("reflect")
	analyzers = append(analyzers, printf.Analyzer,
		structtag.Analyzer,
		osexit.Analyzer,
		nilness.Analyzer,
		lostcancel.Analyzer,
		ifaceassert.Analyzer,
		restrictpkg.RestrictPackageAnalyzer,
	)

	multichecker.Main(analyzers...)
}

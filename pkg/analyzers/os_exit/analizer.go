package osExit

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		fileName := pass.Fset.File(file.Pos()).Name()
		if strings.Contains(fileName, "go-build") {
			continue
		}

		if file.Name.Name != "main" {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			decl, ok := node.(*ast.FuncDecl)
			if !ok {
				return true
			}
			if decl.Name.Name != "main" {
				return false
			}
			hasOsExit := functionHasOsExit(decl)
			if hasOsExit {
				pass.Reportf(
					decl.Pos(),
					"os.Exit in main.main",
				)
			}
			return false
		})
	}
	return nil, nil
}

var Analyzer = &analysis.Analyzer{
	Name: "directExit",
	Doc:  "check for os.Exit in main func of every main package",
	Run:  run,
}

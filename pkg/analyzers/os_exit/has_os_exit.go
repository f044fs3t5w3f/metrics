package os_exit

import "go/ast"

func functionHasOsExit(function *ast.FuncDecl) bool {
	found := false

	ast.Inspect(function, func(node ast.Node) bool {
		if found {
			return false
		}

		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		pkg, ok := selector.X.(*ast.Ident)
		if !ok {
			return true
		}

		if pkg.Name == "os" && selector.Sel.Name == "Exit" {
			found = true
			return false
		}

		return true
	})

	return found
}

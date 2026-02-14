package osexit

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

type testCase struct {
	name      string
	src       string
	hasOsExit bool
}

func parseFunc(t *testing.T, src string) *ast.FuncDecl {

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			return fn
		}
	}

	t.Fatal("function not found")
	return nil
}

func TestFunctionHasOsExit(t *testing.T) {
	tests := []testCase{
		{
			name: "direct os.Exit",
			src: `
package main

import "os"

func main() {
	os.Exit(1)
}
`,
			hasOsExit: true,
		},
		{
			name: "no os.Exit",
			src: `
package main

func main() {
	println("hello")
}
`,
			hasOsExit: false,
		},
		{
			name: "nested os.Exit",
			src: `
package main

import "os"

func main() {
	if true {
		os.Exit(2)
	}
}
`,
			hasOsExit: true,
		},
		{
			name: "method Exit not from os package",
			src: `
package main

type notOs struct{}

func (notOs) Exit(int) {}

func main() {
	var o os
	notOs.Exit(1)
}
`,
			hasOsExit: false,
		},
		{
			name: "just Exit function",
			src: `
package main

func Exit(int) {
}

func main() {
	Exit(1)
}
`,
			hasOsExit: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fn := parseFunc(t, tc.src)

			got := functionHasOsExit(fn)
			if got != tc.hasOsExit {
				t.Errorf("got %v, want %v", got, tc.hasOsExit)
			}
		})
	}
}

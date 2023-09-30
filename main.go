package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "samples/comp.go", nil, 0)
	lastLine := 0
	ast.Inspect(f, func(n ast.Node) bool {
		var s string
		fmt.Printf("%T\n", n)
		switch x := n.(type) {
		case *ast.BasicLit:
			s = x.Value
		case *ast.Ident:
			s = x.Name
		}
		if s != "" {
			var pos = fset.Position(n.Pos())
			fmt.Printf("%s\t", s)
			if pos.Line != lastLine {
				fmt.Println()
				lastLine = pos.Line
			}
		}
		return true
	})
}

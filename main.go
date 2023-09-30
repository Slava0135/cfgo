package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"slava0135/cfgo/graph"
)

func main() {
	fset := token.NewFileSet()
	filename := os.Args[2] // go run main.go -- "samples/block.go" 
	source, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	f, _ := parser.ParseFile(fset, filename, nil, 0)
	var fd *ast.FuncDecl
	ast.Inspect(f, func(n ast.Node) bool {
		if fd != nil {
			return false;
		}
		if fun, ok := n.(*ast.FuncDecl); ok {
			fd = fun
			return false
		} 
		return true
	})
	graph := graph.BuildFuncGraph(source, fd)
	fmt.Println(graph.String())
}

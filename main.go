package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {
	fset := token.NewFileSet()
	filename := "samples/comp.go"
	data, _ := os.ReadFile(filename)
	f, _ := parser.ParseFile(fset, filename, nil, 0)
	var funcDeclarations []*ast.FuncDecl
	ast.Inspect(f, func(n ast.Node) bool {
		if fun, ok := n.(*ast.FuncDecl); ok {
			funcDeclarations = append(funcDeclarations, fun)
			return false
		} 
		return true
	})
	for _, fd := range funcDeclarations {
		fmt.Println("---")
		var graph Graph
		graph.Name = fd.Name.Name
		var start = fd.Body.Pos()
		var end = fd.Body.End()
		var statement = string(data[start:end])
		statement = strings.TrimPrefix(statement, "\n")
		statement = strings.TrimSuffix(statement, "}\n")
		fmt.Printf("%s", statement)
		fmt.Println("---")
	}
}

type Graph struct {
	Name string
	Root *Node
}

type Node struct {
	Text string
	Next []*Node
}

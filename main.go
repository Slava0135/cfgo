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
	if len(os.Args) < 2 {
		log.Fatalln("error: file name was not specified")
	}
	fset := token.NewFileSet()
	filename := os.Args[1]
	source, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	f, _ := parser.ParseFile(fset, filename, nil, 0)
	ast.Inspect(f, func(n ast.Node) bool {
		if fd, ok := n.(*ast.FuncDecl); ok {
			graph := graph.BuildFuncGraph(source, fd)
			fmt.Println(graph.String())
			fmt.Println()
			fmt.Println()
			fmt.Println()
			return false
		} 
		return true
	})
}

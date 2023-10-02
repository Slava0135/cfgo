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
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		log.Fatalln(err)
	}
	var graphsDot string
	ast.Inspect(f, func(n ast.Node) bool {
		if fd, ok := n.(*ast.FuncDecl); ok {
			graph := graph.BuildFuncGraph(source, fd)
			fmt.Println(graph.String())
			fmt.Println(graph.Dot())
			fmt.Println()
			fmt.Println()
			graphsDot += graph.Dot() + "\n"
			return false
		} 
		return true
	})
	outputFile, err := os.Create("cfg.dot")
	if err != nil {
		log.Fatalln(err)
	}
	var output []byte
	output = fmt.Appendf(output, "strict digraph {\n")
	output = fmt.Appendf(output, "\tlabel = \"%s\"\n", filename)
	output = fmt.Appendf(output, "%s}", graphsDot)
	outputFile.WriteString(string(output))
}

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
		graph := funcGraph(data, fd)
		fmt.Printf("%v", graph.Root.Text)
		fmt.Println()
		fmt.Println("---")
	}
}

func funcGraph(data []byte, fd *ast.FuncDecl) *Graph {
	var graph Graph
	graph.Name = fd.Name.Name
	graph.Root = blockStmtNode(data, fd.Body)
	return &graph
}

func blockStmtNode(data []byte, stmt *ast.BlockStmt) *Node {
	var node Node
	node.Next = make([]*Node, 0)
	var start = stmt.Pos()
	var end = stmt.End()
	loop:
	for _, stmt := range stmt.List {
		switch x := stmt.(type) {
		case *ast.IfStmt:
			end = x.Body.Lbrace
			node.Next = append(node.Next, ifStmtNode(data, x))
			break loop
		}
	}
	var text = string(data[start+1:end-1])
	text = levelOutIndent(text)
	node.Text = text
	return &node
}

func levelOutIndent(text string) string {
	var lines = strings.Split(text, "\n")
	const indent = "\t"
	for strings.HasPrefix(lines[0], indent) {
		for i, l := range lines {
			lines[i] = strings.TrimPrefix(l, indent)
		}
	}
	return strings.Join(lines, "\n")
}

func ifStmtNode(data []byte, stmt *ast.IfStmt) *Node {
	var node Node
	return &node
}

type Graph struct {
	Name string
	Root *Node
}

type Node struct {
	Text string
	Next []*Node
}

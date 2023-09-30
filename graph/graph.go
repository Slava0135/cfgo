package graph

import (
	"fmt"
	"go/ast"
)

type Graph struct {
	Name      string
	Source    []byte
	Root      *Node
	NodeCount uint
}

type Node struct {
	Text string
	Next []*Node
}

func BuildFuncGraph(source []byte, fd *ast.FuncDecl) *Graph {
	var graph Graph
	graph.Name = "function: '" + string(fd.Name.Name) + "'"
	graph.Source = source
	return &graph
}

func (g Graph) String() string {
	return fmt.Sprintln(g.Name)
}

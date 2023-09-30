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
	Index uint
	Text  string
}

func BuildFuncGraph(source []byte, fd *ast.FuncDecl) *Graph {
	var graph Graph
	graph.Name = "function: '" + string(fd.Name.Name) + "'"
	graph.Source = source
	graph.Root = graph.blockStmt(fd.Body)
	return &graph
}

func (g *Graph) blockStmt(blockStmt *ast.BlockStmt) *Node {
	var blockNode Node
	blockNode.Index = g.NodeCount
	g.NodeCount += 1
	blockNode.Text = string(g.Source[blockStmt.Lbrace+1:blockStmt.Rbrace-2])
	return &blockNode
}

func (g Graph) String() string {
	var res []byte
	res = fmt.Appendf(res, "%s\n", g.Name)
	res = fmt.Appendf(res, "#%d\n%s", g.Root.Index, g.Root.Text)
	return string(res)
}

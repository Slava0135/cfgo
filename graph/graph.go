package graph

import (
	"fmt"
	"go/ast"
)

type Graph struct {
	Name      string
	Source    []byte
	Root      *Node
	Exit      *Node
	NodeCount uint
}

type Node struct {
	Index uint
	Text  string
	Next  *Node
}

func BuildFuncGraph(source []byte, fd *ast.FuncDecl) *Graph {
	var graph Graph
	graph.Name = "function: '" + string(fd.Name.Name) + "'"
	graph.Source = source
	graph.Root = graph.blockStmt(fd.Body)
	var exitNode Node
	exitNode.Index = graph.NodeCount
	graph.NodeCount += 1
	exitNode.Text = "RETURN"
	graph.Exit = &exitNode
	graph.Root.Next = &exitNode
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
	res = fmt.Appendf(res, "%s", g.Name)
	res = fmt.Appendf(res, "\n[ %d -> %d ]\n%s", g.Root.Index, g.Root.Next.Index, g.Root.Text)
	res = fmt.Appendf(res, "\n[ %d ]\n%s", g.Exit.Index, g.Exit.Text)
	return string(res)
}

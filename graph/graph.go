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
	exitNode := graph.newNode()
	exitNode.Text = "RETURN"
	graph.Exit = exitNode
	graph.Root.Next = exitNode
	return &graph
}

func (g Graph) String() string {
	var res []byte
	res = fmt.Appendf(res, "%s", g.Name)
	res = fmt.Appendf(res, "\n[ %d -> %d ]\n%s", g.Root.Index, g.Root.Next.Index, g.Root.Text)
	res = fmt.Appendf(res, "\n[ %d ]\n%s", g.Exit.Index, g.Exit.Text)
	return string(res)
}

func (g *Graph) newNode() *Node {
	var node Node
	node.Index = g.NodeCount
	g.NodeCount += 1
	return &node
}

func (g *Graph) blockStmt(blockStmt *ast.BlockStmt) *Node {
	entryNode := g.newNode()
	var exitNode = entryNode
	var start = blockStmt.Lbrace+1
	for _, stmt := range blockStmt.List {
		switch s := stmt.(type) {
		case *ast.IfStmt:
			nextNode := g.ifStmt(s)
			exitNode.Text = string(g.Source[start:s.Cond.End()])
			exitNode.Next = nextNode
			exitNode = g.newNode()
			start = s.End()
		}
	}
	exitNode.Text = string(g.Source[start:blockStmt.End()])
	return entryNode
}

func (g *Graph) ifStmt(ifStmt *ast.IfStmt) *Node {
	return nil
}


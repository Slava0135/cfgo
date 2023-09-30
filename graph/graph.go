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
	blockEntryNode, blockExitNode := graph.blockStmt(fd.Body)
	graph.Root = blockEntryNode
	exitNode := graph.newNode()
	exitNode.Text = "RETURN"
	graph.Exit = exitNode
	blockExitNode.Next = exitNode
	return &graph
}

func (g Graph) String() string {
	var res []byte
	res = fmt.Appendf(res, "%s", g.Name)
	var next = g.Root
	for next != nil && next != g.Exit {
		res = fmt.Appendf(res, "\n[ %d -> %d ]\n%s", next.Index, next.Next.Index, next.Text)
		next = next.Next
	}
	res = fmt.Appendf(res, "\n[ %d ]\n%s", g.Exit.Index, g.Exit.Text)
	return string(res)
}

func (g *Graph) newNode() *Node {
	var node Node
	node.Index = g.NodeCount
	g.NodeCount += 1
	return &node
}

func (g *Graph) blockStmt(blockStmt *ast.BlockStmt) (entryNode, exitNode *Node) {
	entryNode = g.newNode()
	exitNode = entryNode
	var start = blockStmt.Lbrace+1
	for _, stmt := range blockStmt.List {
		switch s := stmt.(type) {
		case *ast.IfStmt:
			ifEntryNode, ifExitNode := g.ifStmt(s)
			exitNode.Next = ifEntryNode
			exitNode.Text = string(g.Source[start:s.Cond.End()])
			exitNode = g.newNode()
			ifExitNode.Next = exitNode
			start = s.End()
		}
	}
	exitNode.Text = string(g.Source[start:blockStmt.Rbrace])
	return
}

func (g *Graph) ifStmt(ifStmt *ast.IfStmt) (entryNode, exitNode *Node) {
	entryNode, exitNode = g.blockStmt(ifStmt.Body)
	return
}


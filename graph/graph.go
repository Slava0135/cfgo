package graph

import (
	"fmt"
	"go/ast"
	"strings"
)

type Graph struct {
	Name      string
	Source    []byte
	Root      *Node
	Exit      *Node
	NodeCount uint
	AllNodes  []*Node
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
	for _, n := range g.AllNodes {
		if n == g.Exit {
			break
		}
		res = fmt.Appendf(res, "\n[ %d -> %d ]\n%s", n.Index, n.Next.Index, n.Text)
	}
	res = fmt.Appendf(res, "\n[ %d ]\n%s", g.Exit.Index, g.Exit.Text)
	return string(res)
}

func (g *Graph) newNode() *Node {
	var node Node
	node.Index = g.NodeCount
	g.NodeCount += 1
	g.AllNodes = append(g.AllNodes, &node)
	return &node
}

func (g *Graph) blockStmt(blockStmt *ast.BlockStmt) (entryNode, exitNode *Node) {
	entryNode = g.newNode()
	exitNode = entryNode
	var start = blockStmt.Lbrace+1
	for _, stmt := range blockStmt.List {
		switch s := stmt.(type) {
		case *ast.IfStmt:
			exitNode.Text = string(g.Source[start:s.Cond.End()])
			ifEntryNode, ifExitNode := g.ifStmt(s)
			exitNode.Next = ifEntryNode
			if (ifEntryNode != ifExitNode) {
				exitNode = ifExitNode
			} else {
				exitNode = g.newNode()
				ifExitNode.Next = exitNode
			}
			start = s.End()
		}
	}
	text := string(g.Source[start:blockStmt.Rbrace])
	lines := strings.Split(text, "\n")
	exitNode.Text = strings.Join(lines[:len(lines)-1], "\n")
	return
}

func (g *Graph) ifStmt(ifStmt *ast.IfStmt) (entryNode, exitNode *Node) {
	entryNode, exitNode = g.blockStmt(ifStmt.Body)
	return
}


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
	AllNodes  []*Node
}

type Node struct {
	Index int
	Text  string
	Next  []*Node
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
	blockExitNode.Next = append(blockExitNode.Next, exitNode)
	return &graph
}

func (g Graph) String() string {
	var res []byte
	res = fmt.Appendf(res, "%s", g.Name)
	for _, n := range g.AllNodes {
		if n == g.Exit {
			break
		}
		if len(n.Next) == 0 {
			fmt.Appendf(res, "\n[ %d ]\n%s", n.Index, n.Text)
		} else {
			res = fmt.Appendf(res, "\n[ %d -> ", n.Index)
			for _, next := range n.Next {
				res = fmt.Appendf(res, "%d ", next.Index)
			}
			res = fmt.Appendf(res, "]\n%s", n.Text)
		}
	}
	res = fmt.Appendf(res, "\n[ %d ]\n%s\n", g.Exit.Index, g.Exit.Text)
	return string(res)
}

func (g *Graph) newNode() *Node {
	var node Node
	node.Index = len(g.AllNodes)
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
			exitNode.Next = append(exitNode.Next, ifExitNode)
			if (ifEntryNode != ifExitNode) {
				exitNode.Next = append(exitNode.Next, ifEntryNode)
				exitNode = ifExitNode
			} else {
				exitNode = g.newNode()
				ifExitNode.Next = append(ifExitNode.Next, exitNode)
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


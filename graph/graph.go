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

type Flow struct {
	Entries map[*Node] struct{}
	Exit *Node
}

func BuildFuncGraph(source []byte, fd *ast.FuncDecl) *Graph {
	var graph Graph
	graph.Name = "function: '" + string(fd.Name.Name) + "'"
	graph.Source = source
	flow := graph.blockStmt(fd.Body)
	graph.Root = graph.AllNodes[0]
	exitNode := graph.newNode()
	exitNode.Text = "RETURN"
	graph.Exit = exitNode
	flow.Exit.Next = append(flow.Exit.Next, exitNode)
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

func (g *Graph) blockStmt(blockStmt *ast.BlockStmt) Flow {
	var entryNode = g.newNode()
	var exitNode = entryNode
	var start = blockStmt.Lbrace+1
	for _, stmt := range blockStmt.List {
		switch s := stmt.(type) {
		case *ast.IfStmt:
			exitNode.Text = string(g.Source[start:s.Cond.End()])
			ifFlow := g.ifStmt(s)
			exitNode.Next = append(exitNode.Next, ifFlow.Exit)
			if _, ok := ifFlow.Entries[ifFlow.Exit]; ok {
				exitNode = g.newNode()
				ifFlow.Exit.Next = append(ifFlow.Exit.Next, exitNode)
			} else {
				for e := range ifFlow.Entries {
					exitNode.Next = append(exitNode.Next, e)
				}
				exitNode = ifFlow.Exit
			}
			start = s.End()
		}
	}
	text := string(g.Source[start:blockStmt.Rbrace])
	lines := strings.Split(text, "\n")
	exitNode.Text = strings.Join(lines[:len(lines)-1], "\n")
	var flow Flow
	flow.Entries = make(map[*Node]struct{})
	flow.Entries[entryNode] = struct{}{}
	flow.Exit = exitNode
	return flow
}

func (g *Graph) ifStmt(ifStmt *ast.IfStmt) Flow {
	blockFlow := g.blockStmt(ifStmt.Body)
	if ifStmt.Else != nil {
		var elseFlow Flow
		switch s := ifStmt.Else.(type) {
		case *ast.BlockStmt:
			elseFlow = g.blockStmt(s)
		case *ast.IfStmt:
			elseFlow = g.ifStmt(s)
		}
		var fullFlow Flow
		fullFlow.Entries = make(map[*Node]struct{})
		for n := range blockFlow.Entries {
			fullFlow.Entries[n] = struct{}{}
		}
		for n := range elseFlow.Entries {
			fullFlow.Entries[n] = struct{}{}
		}
		fullFlow.Exit = blockFlow.Exit
		return fullFlow
	}
	return blockFlow
}


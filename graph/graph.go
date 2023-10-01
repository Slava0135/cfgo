package graph

import (
	"fmt"
	"go/ast"
	"strings"
)

type Graph struct {
	Name     string
	Source   []byte
	Root     *Node
	Exit     *Node
	AllNodes []*Node
}

type Node struct {
	Index int
	Text  string
	Next  []*Node
}

type Flow struct {
	Entries map[*Node]struct{}
	Exits   map[*Node]struct{}
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
	for n := range flow.Exits {
		n.Next = append(n.Next, exitNode)
	}
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

func newFlow() *Flow {
	var flow Flow
	flow.Entries = make(map[*Node]struct{})
	flow.Exits = make(map[*Node]struct{})
	return &flow
}

func mergeFlows(into, from *Flow) {
	for n := range from.Entries {
		into.Entries[n] = struct{}{}
	}
	for n := range from.Exits {
		into.Exits[n] = struct{}{}
	}
}

func (g *Graph) blockStmt(blockStmt *ast.BlockStmt) *Flow {
	var entryNode = g.newNode()
	var exitNode = entryNode
	var start = blockStmt.Lbrace + 1
	for _, stmt := range blockStmt.List {
		switch s := stmt.(type) {
		case *ast.IfStmt:
			exitNode.Text = string(g.Source[start:s.Cond.End()])
			ifFlow := g.ifStmt(s)
			var nextNode = g.newNode()
			for n := range ifFlow.Exits {
				n.Next = append(n.Next, nextNode)
			}
			for n := range ifFlow.Entries {
				exitNode.Next = append(exitNode.Next, n)
			}
			exitNode.Next = append(exitNode.Next, nextNode)
			exitNode = nextNode
			start = s.End()
		}
	}
	text := string(g.Source[start:blockStmt.Rbrace])
	lines := strings.Split(text, "\n")
	exitNode.Text = strings.Join(lines[:len(lines)-1], "\n")
	flow := newFlow()
	flow.Entries[entryNode] = struct{}{}
	flow.Exits[exitNode] = struct{}{}
	return flow
}

func (g *Graph) ifStmt(ifStmt *ast.IfStmt) *Flow {
	blockFlow := g.blockStmt(ifStmt.Body)
	if ifStmt.Else != nil {
		var elseFlow *Flow
		switch s := ifStmt.Else.(type) {
		case *ast.BlockStmt:
			elseFlow = g.blockStmt(s)
		case *ast.IfStmt:
			elseFlow = g.ifStmt(s)
		}
		fullFlow := newFlow()
		mergeFlows(fullFlow, blockFlow)
		mergeFlows(fullFlow, elseFlow)
		return fullFlow
	}
	return blockFlow
}

package graph

import (
	"fmt"
	"go/ast"
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

func BuildFuncGraph(source []byte, fd *ast.FuncDecl) *Graph {
	var graph Graph
	graph.Name = "function: '" + string(fd.Name.Name) + "'"
	graph.Source = source
	var exit = graph.newNode()
	graph.Exit = exit
	graph.Root = graph.blockStmt(fd.Body, exit)
	exit.Text = "RETURN"
	return &graph
}

func (g Graph) String() string {
	var res []byte
	res = fmt.Appendf(res, "%s", g.Name)
	for _, n := range g.AllNodes {
		if n == g.Exit {
			continue
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

func (g *Graph) blockStmt(blockStmt *ast.BlockStmt, exit *Node) *Node {
	var entry = g.newNode()
	var next = entry
	var text = ""
	for i, stmt := range blockStmt.List {
		switch s := stmt.(type) {
		case *ast.IfStmt:
			var ifExit *Node
			if i < len(blockStmt.List) - 1 {
				ifExit = g.newNode()
			} else {
				ifExit = exit
			}
			var ifEntry = g.ifStmt(s, ifExit)
			next.Next = append(next.Next, ifEntry)
			next.Text = text
			text = ""
			next = ifExit
			continue
		}
		text += string(g.Source[stmt.Pos()-1:stmt.End()])
	}
	if next != exit {
		next.Text = text
		next.Next = append(next.Next, exit)
	}
	return entry
}

func (g *Graph) ifStmt(ifStmt *ast.IfStmt, exit *Node) *Node {
	var entry = g.newNode()
	var blockEntry = g.blockStmt(ifStmt.Body, exit)
	entry.Next = append(entry.Next, blockEntry)
	entry.Text = string(g.Source[ifStmt.If-1:ifStmt.Cond.End()])
	if ifStmt.Else != nil {
		switch s := ifStmt.Else.(type) {
		case *ast.BlockStmt:
			var elseEntry = g.blockStmt(s, exit)
			entry.Next = append(entry.Next, elseEntry) 
		case *ast.IfStmt:
			var elseIfEntry = g.ifStmt(s, exit)
			entry.Next = append(entry.Next, elseIfEntry)
		}
	}
	return entry
}

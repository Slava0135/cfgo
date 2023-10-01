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
	graph.createIndex(exit)
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
	node.Index = -1
	return &node
}

func (g *Graph) createIndex(node *Node) {
	if node.Index >= 0 {
		return
	}
	node.Index = len(g.AllNodes)
	g.AllNodes = append(g.AllNodes, node)
}

func (g *Graph) blockStmt(blockStmt *ast.BlockStmt, exit *Node) *Node {
	var first *Node
	var last *Node
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
			if first == nil {
				first = ifEntry
			} else {
				last.Next = append(last.Next, ifEntry)
				last.Text = text
				text = ""
			}
			last = ifExit
			g.createIndex(last)
			continue
		}
		if first == nil {
			first = g.newNode()
			g.createIndex(first)
			last = first
		}
		text += string(g.Source[stmt.Pos()-1:stmt.End()])
	}
	if last != exit {
		last.Text = text
		last.Next = append(last.Next, exit)
	}
	return first
}

func (g *Graph) ifStmt(ifStmt *ast.IfStmt, exit *Node) *Node {
	var entry = g.newNode()
	g.createIndex(entry)
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
	} else {
		entry.Next = append(entry.Next, exit)
	}
	if len(entry.Next) != 2 {
		panic("if block must have 2 branches")
	}
	return entry
}

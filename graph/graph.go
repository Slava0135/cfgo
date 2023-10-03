package graph

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type Graph struct {
	Name     string
	Source   []byte
	Root     *Node
	Exit     *Node
	AllNodes []*Node
	LoopEnd  *Node
	LoopPost *Node
}

type Node struct {
	Index int
	Text  string
	Next  []Link
	Kind  Kind
}

type Link struct {
	Node *Node
	Info string
}

type Kind int

const (
	SEQUENCE Kind = iota
	CONDITION
	BRANCH
)

func BuildFuncGraph(source []byte, fd *ast.FuncDecl) *Graph {
	var graph Graph
	graph.Name = fd.Name.Name
	graph.Source = source
	var exit = graph.newNode()
	exit.Kind = BRANCH
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
		var text = strings.TrimSuffix(n.Text, "\n")
		if len(n.Next) == 0 {
			fmt.Appendf(res, "\n-- %d --\n%s", n.Index, text)
		} else {
			res = fmt.Appendf(res, "\n-- %d >> ", n.Index)
			for _, next := range n.Next {
				res = fmt.Appendf(res, "%d ", next.Node.Index)
			}
			res = fmt.Appendf(res, "--\n%s\n", text)
		}
	}
	res = fmt.Appendf(res, "\n-- %d --\n%s\n", g.Exit.Index, g.Exit.Text)
	return string(res)
}

func (g Graph) Dot() string {
	var res []byte
	res = fmt.Appendf(res, "subgraph cluster_%s {\n", g.Name)
	res = fmt.Appendf(res, "\tlabel=\"%s\"\n", g.Name)
	for _, node := range g.AllNodes {
		text := strings.ReplaceAll(node.Text, "\n", "\\n")
		text = strings.ReplaceAll(text, "\"", "'")
		var shape = "box"
		switch node.Kind {
		case CONDITION:
			shape = "diamond"
		case BRANCH:
			shape = "cds"
		}
		res = fmt.Appendf(res, "\t%s_%d [shape=%s, label=\"%s\"]\n", g.Name, node.Index, shape, text)
	}
	for _, source := range g.AllNodes {
		for _, dest := range source.Next {
			res = fmt.Appendf(res, "\t%s_%d -> %s_%d [label=\"%s\"]\n", g.Name, source.Index, g.Name, dest.Node.Index, dest.Info)
		}
	}
	res = fmt.Appendf(res, "}")
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
	if len(blockStmt.List) == 0 {
		var node = g.newNode()
		g.createIndex(node)
		node.Text = "EMPTY BLOCK"
		node.Next = append(node.Next, Link{exit, ""})
		return node
	}
	var first *Node
	var last *Node
	var text = ""
	for i, stmt := range blockStmt.List {
		processInnerStmt := func(process func(innerExit *Node) *Node) {
			var innerExit *Node
			if i < len(blockStmt.List) - 1 {
				innerExit = g.newNode()
			} else {
				innerExit = exit
			}
			var innerEntry = process(innerExit)
			if first == nil {
				first = innerEntry
			} else {
				last.Next = append(last.Next, Link{innerEntry, ""})
				last.Text = text
				text = ""
			}
			last = innerExit
			g.createIndex(last)
		}
		switch s := stmt.(type) {
		case *ast.IfStmt:
			processInnerStmt(func(innerExit *Node) *Node {
				return g.ifStmt(s, innerExit)
			})
			continue
		case *ast.ForStmt:
			processInnerStmt(func(innerExit *Node) *Node {
				return g.forStmt(s, innerExit)
			})
			continue
		case *ast.RangeStmt:
			processInnerStmt(func(innerExit *Node) *Node {
				return g.rangeStmt(s, innerExit)
			})
			continue
		case *ast.ReturnStmt:
			if first == nil {
				first = g.newNode()
				g.createIndex(first)
				last = first
			}
			text += string(g.Source[s.Pos()-1:s.End()])
			last.Text = text
			last.Next = append(last.Next, Link{g.Exit, ""})
			last.Kind = BRANCH
			return first
		case *ast.BranchStmt:
			var gotoNode *Node
			switch s.Tok {
			case token.BREAK:
				gotoNode = g.LoopEnd
			case token.CONTINUE:
				gotoNode = g.LoopPost
			}
			if gotoNode == nil {
				continue
			}
			if first == nil {
				first = g.newNode()
				g.createIndex(first)
				last = first
			}
			text += string(g.Source[s.Pos()-1:s.End()])
			last.Text = text
			last.Next = append(last.Next, Link{gotoNode, ""})
			last.Kind = BRANCH
			return first
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
		last.Next = append(last.Next, Link{exit, ""})
	}
	return first
}

func (g *Graph) ifStmt(ifStmt *ast.IfStmt, exit *Node) *Node {
	var entry = g.newNode()
	g.createIndex(entry)
	entry.Kind = CONDITION
	var blockEntry = g.blockStmt(ifStmt.Body, exit)
	entry.Next = append(entry.Next, Link{blockEntry, "true"})
	entry.Text = string(g.Source[ifStmt.Cond.Pos()-1:ifStmt.Cond.End()])
	if ifStmt.Else != nil {
		switch s := ifStmt.Else.(type) {
		case *ast.BlockStmt:
			var elseEntry = g.blockStmt(s, exit)
			entry.Next = append(entry.Next, Link{elseEntry, "false"}) 
		case *ast.IfStmt:
			var elseIfEntry = g.ifStmt(s, exit)
			entry.Next = append(entry.Next, Link{elseIfEntry, "false"})
		}
	} else {
		entry.Next = append(entry.Next, Link{exit, "false"})
	}
	if len(entry.Next) != 2 {
		panic("if block must have 2 branches")
	}
	return entry
}

func (g *Graph) forStmt(forStmt *ast.ForStmt, exit *Node) *Node {
	var entry = g.newNode()
	var condition = entry 
	if forStmt.Init != nil {
		entry = g.newNode()
		g.createIndex(entry)
		entry.Next = append(entry.Next, Link{condition, ""})
		entry.Text = string(g.Source[forStmt.Init.Pos()-1:forStmt.Init.End()])
	}
	g.createIndex(condition)
	condition.Kind = CONDITION
	var post = condition
	if forStmt.Post != nil {
		post = g.newNode()
		g.createIndex(post)
		post.Next = append(post.Next, Link{condition, ""})
		post.Text = string(g.Source[forStmt.Post.Pos()-1:forStmt.Post.End()])
	}
	var prevLoopEnd = g.LoopEnd
	defer func() { g.LoopEnd = prevLoopEnd }()
	var prevLoopPost = g.LoopPost
	defer func() { g.LoopPost = prevLoopPost }()
	g.LoopEnd = exit
	g.LoopPost = post
	var blockEntry = g.blockStmt(forStmt.Body, post)
	condition.Next = append(condition.Next, Link{blockEntry, "true"})
	condition.Next = append(condition.Next, Link{exit, "false"})
	if forStmt.Cond != nil {
		condition.Text = string(g.Source[forStmt.Cond.Pos()-1:forStmt.Cond.End()])
	} else {
		condition.Text = "FOREVER"
	}
	return entry
}

func (g *Graph) rangeStmt(rangeStmt *ast.RangeStmt, exit *Node) *Node {
	var entry = g.newNode()
	g.createIndex(entry)
	entry.Kind = CONDITION
	entry.Text = string(g.Source[rangeStmt.Pos()-1:rangeStmt.Body.Lbrace-2])
	var prevLoopEnd = g.LoopEnd
	defer func() { g.LoopEnd = prevLoopEnd }()
	var prevLoopPost = g.LoopPost
	defer func() { g.LoopPost = prevLoopPost }()
	g.LoopEnd = exit
	g.LoopPost = entry
	var blockEntry = g.blockStmt(rangeStmt.Body, entry)
	entry.Next = append(entry.Next, Link{blockEntry, "not empty"})
	entry.Next = append(entry.Next, Link{exit, "empty"})
	return entry
}

// TODO: switch

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
	LoopEnd  *Node
	LoopPost *Node
	Returns  []*Node
}

type Node struct {
	Index int
	Text  string
	Next  map[*Node]string
	Kind  Kind
}

type Kind int

type Connection struct {
	Entry *Node
	Exits []*Node
}

const (
	SEQUENCE Kind = iota
	CONDITION
	BRANCH
)

func BuildFuncGraph(source []byte, fd *ast.FuncDecl) *Graph {
	var graph Graph
	graph.Name = fd.Name.Name
	graph.Source = source
	conn, _ := graph.listStmt(fd.Body.List, nil)
	graph.Root = conn.Entry
	var exit = graph.newNode()
	exit.Text = "EXIT"
	exit.Kind = BRANCH
	graph.Exit = exit
	for _, e := range conn.Exits {
		e.Next[exit] = ""
	}
	for _, r := range graph.Returns {
		r.Next[exit] = ""
	}
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
			for node := range n.Next {
				res = fmt.Appendf(res, "%d ", node.Index)
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
		for dest, info := range source.Next {
			res = fmt.Appendf(res, "\t%s_%d -> %s_%d [label=\"%s\"]\n", g.Name, source.Index, g.Name, dest.Index, info)
		}
	}
	res = fmt.Appendf(res, "}")
	return string(res)
}

func (g *Graph) newNode() *Node {
	var node Node
	node.Index = -1
	node.Next = make(map[*Node]string)
	node.Index = len(g.AllNodes)
	g.AllNodes = append(g.AllNodes, &node)
	return &node
}

func (g *Graph) listStmt(listStmt []ast.Stmt, prev *Node) (conn Connection, empty bool) {
	if len(listStmt) == 0 {
		empty = true
		return
	}
	text := ""
	var listConns []Connection
	pushText := func() *Node {
		if text != "" {
			node := g.newNode()
			node.Text = text
			text = ""
			var conn Connection
			conn.Entry = node
			conn.Exits = append(conn.Exits, node)
			listConns = append(listConns, conn)
			return node
		}
		return nil
	}
	connectAll := func() {
		conn.Entry = listConns[0].Entry
		for i := 0; i+1 < len(listConns); i += 1 {
			for _, e := range listConns[i].Exits {
				e.Next[listConns[i+1].Entry] = ""
			}
		}
	}
	for _, stmt := range listStmt {
		switch s := stmt.(type) {
		case *ast.IfStmt:
			pushText()
			listConns = append(listConns, g.ifStmt(s))
		case *ast.ForStmt:
			pushText()
			listConns = append(listConns, g.forStmt(s))
		case *ast.ReturnStmt:
			text += string(g.Source[stmt.Pos()-1 : stmt.End()])
			last := pushText()
			connectAll()
			last.Kind = BRANCH
			g.Returns = append(g.Returns, last)
			return
		case *ast.BranchStmt:
			pushText()
			connectAll()
			if len(listConns) == 0 {
				prev.Next[g.LoopPost] = ""
			} else {
				for _, e := range listConns[len(listConns)-1].Exits {
					e.Next[g.LoopPost] = ""
				}
			}
			return
		default:
			text += string(g.Source[stmt.Pos()-1 : stmt.End()])
		}
	}
	pushText()
	connectAll()
	conn.Exits = listConns[len(listConns)-1].Exits
	return
}

func (g *Graph) ifStmt(ifStmt *ast.IfStmt) (conn Connection) {
	condition := g.newNode()
	condition.Kind = CONDITION
	condition.Text = string(g.Source[ifStmt.Cond.Pos()-1 : ifStmt.Cond.End()])
	conn.Entry = condition
	bodyConn, empty := g.listStmt(ifStmt.Body.List, condition)
	if !empty {
		condition.Next[bodyConn.Entry] = "true"
		conn.Exits = append(conn.Exits, bodyConn.Exits...)
	}
	if ifStmt.Else == nil {
		conn.Exits = append(conn.Exits, condition)
	} else {
		var elseConn Connection
		switch s := ifStmt.Else.(type) {
		case *ast.BlockStmt:
			elseConn, empty = g.listStmt(s.List, condition)
		case *ast.IfStmt:
			elseConn = g.ifStmt(s)
		}
		condition.Next[elseConn.Entry] = "false"
		conn.Exits = append(conn.Exits, elseConn.Exits...)
	}
	return
}

func (g *Graph) forStmt(forStmt *ast.ForStmt) (conn Connection) {
	var init *Node
	var condition = g.newNode()
	condition.Kind = CONDITION
	if forStmt.Cond != nil {
		text := string(g.Source[forStmt.Cond.Pos()-1 : forStmt.Cond.End()])
		condition.Text = strings.TrimSuffix(text, ";")
	} else {
		condition.Text = ""
	}
	if forStmt.Init != nil {
		init = g.newNode()
		text := string(g.Source[forStmt.Init.Pos()-1 : forStmt.Init.End()])
		init.Text = strings.TrimSuffix(text, ";")
		init.Next[condition] = ""
	}
	if init == nil {
		init = condition
	}
	conn.Entry = init
	conn.Exits = append(conn.Exits, condition)
	var post *Node
	if forStmt.Post != nil {
		post = g.newNode()
		text := string(g.Source[forStmt.Post.Pos()-1 : forStmt.Post.End()])
		post.Text = strings.TrimSuffix(text, ";")
		post.Next[condition] = ""
	}
	if post == nil {
		post = condition
	}
	var prevLoopPost = g.LoopPost
	g.LoopPost = condition
	bodyConn, empty := g.listStmt(forStmt.Body.List, condition)
	g.LoopPost = prevLoopPost
	if !empty {
		condition.Next[bodyConn.Entry] = "true"
		for _, e := range bodyConn.Exits {
			e.Next[post] = ""
		}
	} else {
		condition.Next[post] = "true"
	}
	return
}

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
	Next  map[*Node]string
	Kind  Kind
}

type Kind int

type Connection struct {
	Entry *Node
	Exits []*Exit
}

type Exit struct {
	Node *Node
	Type ExitType
}

type ExitType int

const (
	SEQUENCE Kind = iota
	CONDITION
	BRANCH
)

const (
	NORMAL = iota
	RETURN
	CONTINUE
)

func BuildFuncGraph(source []byte, fd *ast.FuncDecl) *Graph {
	var graph Graph
	graph.Name = fd.Name.Name
	graph.Source = source
	conn, _ := graph.listStmt(fd.Body.List)
	graph.Root = conn.Entry
	var exit = graph.newNode()
	exit.Text = "EXIT"
	exit.Kind = BRANCH
	graph.Exit = exit
	for _, e := range conn.Exits {
		e.Node.Next[exit] = ""
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

func (g *Graph) listStmt(listStmt []ast.Stmt) (conn Connection, empty bool) {
	if len(listStmt) == 0 {
		empty = true
		return
	}
	text := ""
	var listConns []Connection
	pushText := func(exitType ExitType) *Node {
		if text != "" {
			node := g.newNode()
			node.Text = text
			text = ""
			var conn Connection
			conn.Entry = node
			conn.Exits = append(conn.Exits, &Exit{node, exitType})
			listConns = append(listConns, conn)
			return node
		}
		return nil
	}
	connectAll := func() {
		if len(listConns) == 0 {
			return
		}
		conn.Entry = listConns[0].Entry
		for i := 0; i+1 < len(listConns); i += 1 {
			for _, e := range listConns[i].Exits {
				if e.Type == NORMAL {
					e.Node.Next[listConns[i+1].Entry] = ""
				} else {
					conn.Exits = append(conn.Exits, e)
				}
			}
		}
		return
	}
	loop:
	for _, stmt := range listStmt {
		switch s := stmt.(type) {
		case *ast.IfStmt:
			pushText(NORMAL)
			listConns = append(listConns, g.ifStmt(s))
		case *ast.ForStmt:
			pushText(NORMAL)
			listConns = append(listConns, g.forStmt(s))
		case *ast.ReturnStmt:
			text += string(g.Source[stmt.Pos()-1 : stmt.End()])
			last := pushText(RETURN)
			last.Kind = BRANCH
			break loop
		case *ast.BranchStmt:
			pushText(CONTINUE)
			if len(listConns) == 0 {
				conn.Exits = append(conn.Exits, &Exit{nil, CONTINUE})
			} else {
				for _, e := range listConns[len(listConns)-1].Exits {
					if e.Type == NORMAL {
						conn.Exits = append(conn.Exits, &Exit{e.Node, CONTINUE})
					}
				}
			}
			break loop
		default:
			text += string(g.Source[stmt.Pos()-1 : stmt.End()])
		}
	}
	pushText(NORMAL)
	if len(listConns) == 0 {
		empty = true
		return
	}
	connectAll()
	conn.Exits = append(conn.Exits, listConns[len(listConns)-1].Exits...)
	return
}

func (g *Graph) ifStmt(ifStmt *ast.IfStmt) (conn Connection) {
	condition := g.newNode()
	condition.Kind = CONDITION
	condition.Text = string(g.Source[ifStmt.Cond.Pos()-1 : ifStmt.Cond.End()])
	conn.Entry = condition
	bodyConn, empty := g.listStmt(ifStmt.Body.List)
	if !empty {
		condition.Next[bodyConn.Entry] = "true"
		conn.Exits = append(conn.Exits, bodyConn.Exits...)
	}
	if ifStmt.Else == nil {
		conn.Exits = append(conn.Exits, &Exit{condition, NORMAL})
	} else {
		var elseConn Connection
		switch s := ifStmt.Else.(type) {
		case *ast.BlockStmt:
			elseConn, empty = g.listStmt(s.List)
			if !empty {
				condition.Next[elseConn.Entry] = "false"
				conn.Exits = append(conn.Exits, elseConn.Exits...)
			} else {
				conn.Exits = append(conn.Exits, &Exit{condition, NORMAL})
			}
		case *ast.IfStmt:
			elseConn = g.ifStmt(s)
			condition.Next[elseConn.Entry] = "false"
			conn.Exits = append(conn.Exits, elseConn.Exits...)
		}
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
	conn.Exits = append(conn.Exits, &Exit{condition, NORMAL})
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
	bodyConn, empty := g.listStmt(forStmt.Body.List)
	if !empty {
		condition.Next[bodyConn.Entry] = "true"
		for _, e := range bodyConn.Exits {
			if e.Type == NORMAL || e.Type == CONTINUE {
				e.Node.Next[post] = ""
			} else {
				conn.Exits = append(conn.Exits, e)
			}
		}
	} else {
		condition.Next[post] = "true"
	}
	return
}

package graph

import (
	"go/ast"
	"strings"
)

func FuncGraph(data []byte, fd *ast.FuncDecl) *Graph {
	var graph Graph
	graph.Name = fd.Name.Name
	graph.Root = blockStmtNode(data, fd.Body)
	return &graph
}

func blockStmtNode(data []byte, stmt *ast.BlockStmt) *Node {
	var firstNode Node
	var lastNode = &firstNode
	var start = stmt.Pos()+1
	for _, stmt := range stmt.List {
		switch x := stmt.(type) {
		case *ast.IfStmt:
			lastNode.Text = string(data[start:x.Body.Lbrace-1])
			var nextNode Node
			b, e := ifStmtNode(data, x, &nextNode) 
			lastNode.Next = append(lastNode.Next, b)
			if e != nil {
				lastNode.Next = append(lastNode.Next, e)
			}
			lastNode = &nextNode
			start = x.End()
		}
	}
	return &firstNode
}

func levelOutIndent(text string) string {
	var lines = strings.Split(text, "\n")
	const indent = "\t"
	for strings.HasPrefix(lines[0], indent) {
		for i, l := range lines {
			lines[i] = strings.TrimPrefix(l, indent)
		}
	}
	return strings.Join(lines, "\n")
}

func ifStmtNode(data []byte, stmt *ast.IfStmt, next *Node) (bodyNode *Node, elseNode *Node) {
	bodyNode = blockStmtNode(data, stmt.Body)
	bodyNode.Next = append(bodyNode.Next, next)
	return
}

type Graph struct {
	Name string
	Root *Node
}

type Node struct {
	Text string
	Next []*Node
}

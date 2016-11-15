package pkg_graph

import (
	"fmt"
	"go/ast"
	"go/types"
)

type PkgName string
type FuncName string
type CallStats map[PkgName]map[FuncName]int

func (c CallStats) inc(pkg string, fn string) {
	pelement, ok := c[PkgName(pkg)]
	if !ok {
		c[PkgName(pkg)] = make(map[FuncName]int)
		c[PkgName(pkg)][FuncName(fn)] = 1
	} else {
		_, ok := pelement[FuncName(fn)]
		if !ok {
			pelement[FuncName(fn)] = 1
		} else {
			pelement[FuncName(fn)]++
		}
	}
}

type PkgNode struct {
	Node      *types.Package
	Parents   []*PkgNode
	Children  []*PkgNode
	Files     []*ast.File
	CallStats CallStats
}

func NewPkgNode(root *types.Package, files []*ast.File) *PkgNode {
	top := &PkgNode{
		Node:      root,
		Parents:   make([]*PkgNode, 0),
		Children:  make([]*PkgNode, 0),
		Files:     files,
		CallStats: make(CallStats),
	}

	return top
}

func (n *PkgNode) FullName() string {
	return n.Node.Path()
}

func (n *PkgNode) ShortName() string {
	return n.Node.Name()
}

func (n *PkgNode) TotalFuncDecls() int {
	nFuncs := 0
	for _, file := range n.Files {
		for _, obj := range file.Scope.Objects {
			if obj.Kind == ast.Fun {
				nFuncs++
			}
		}
	}

	return nFuncs
}

func (n *PkgNode) CalcCallStats() {
	for _, file := range n.Files {
		counter := NewCallCounter(n.CallStats)
		ast.Walk(counter, file)
	}
}

func (n *PkgNode) CallStatsEdge(pkg PkgName) int {
	entry, ok := n.CallStats[pkg]
	if !ok {
		return 0
	}
	acc := 0
	for f, i := range entry {
		fmt.Printf("%v %v %v %v \n", n.FullName(), pkg, f, i)
		acc += i
	}

	return acc
}

type CallCounter struct {
	CallStats CallStats
}

func NewCallCounter(stats CallStats) *CallCounter {
	return &CallCounter{
		CallStats: stats,
	}
}

func (v *CallCounter) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		w = nil
		return
	}

	switch nodeObj := node.(type) {
	case *ast.SelectorExpr:
		switch xObj := nodeObj.X.(type) {
		case *ast.Ident:
			if xObj.Obj == nil {
				v.CallStats.inc(xObj.Name, nodeObj.Sel.Name)
			}
		}
	}

	w = v
	return
}

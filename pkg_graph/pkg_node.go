package pkg_graph

import (
	"fmt"
	"go/ast"
	"go/types"
	"reflect"
)

type PkgName string
type FuncName string
type CallStats map[PkgName]map[FuncName]int

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

	callExpr, ok := node.(*ast.CallExpr)
	if ok {
		fmt.Printf("%v of type %v\n", callExpr.Fun, reflect.TypeOf(callExpr.Fun))
	}

	w = v
	return
}

package pkg_graph

import (
	"go/ast"
	"go/types"
)

type FuncID string

type CallStats map[*PkgNode]map[FuncID]int

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

func (n *PkgNode) Visit(node ast.Node) (w ast.Visitor) {
	return nil
}

package pkg_graph

import (
	"go/types"
	"strings"
)

type Filter struct {
	includeStdLib bool
	permit        string
	deny          string
}

func NewFilter(includeStdLib bool, permit string, deny string) *Filter {
	return &Filter{
		includeStdLib: includeStdLib,
		permit:        permit,
		deny:          deny,
	}
}

type PkgNode struct {
	Node     *types.Package
	Parents  []*PkgNode
	Children []*PkgNode
}

func NewPkgNode(root *types.Package) *PkgNode {
	top := &PkgNode{
		Node:     root,
		Parents:  make([]*PkgNode, 0),
		Children: make([]*PkgNode, 0),
	}

	return top
}

type PkgGraph struct {
	Nodes  map[string]*PkgNode
	Filter *Filter
}

func NewPkgGraph(filter *Filter) *PkgGraph {
	return &PkgGraph{
		Nodes:  make(map[string]*PkgNode),
		Filter: filter,
	}
}

func (g *PkgGraph) Size() int {
	return len(g.Nodes)
}

func (p *PkgGraph) Populate(n *PkgNode) {
	_, ok := p.Nodes[n.Node.Path()]
	if !ok {
		p.Nodes[n.Node.Path()] = n

		for _, c := range n.Node.Imports() {
			cpath := c.Path()

			if !p.Filter.includeStdLib && isStandardImportPath(cpath) {
				continue
			}

			if p.Filter.permit != "" {
				if !strings.Contains(cpath, p.Filter.permit) {
					continue
				}
			}

			if p.Filter.deny != "" {
				if strings.Contains(cpath, p.Filter.deny) {
					continue
				}
			}

			cNode, ok := p.Nodes[cpath]
			if !ok {
				cNode = NewPkgNode(c)
				p.Populate(cNode)
			}
			cNode.Parents = append(cNode.Parents, n)
			n.Children = append(n.Children, cNode)
		}
	}
}

// Imported
// See go/src/cmd/go/pkg.go
//
// isStandardImportPath reports whether $GOROOT/src/path should be considered
// part of the standard distribution. For historical reasons we allow people to add
// their own code to $GOROOT instead of using $GOPATH, but we assume that
// code will start with a domain name (dot in the first element).
func isStandardImportPath(path string) bool {
	i := strings.Index(path, "/")
	if i < 0 {
		i = len(path)
	}
	elem := path[:i]
	return !strings.Contains(elem, ".")
}

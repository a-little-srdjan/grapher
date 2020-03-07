package model

import (
	"go/types"
	"regexp"
	"strings"

	"golang.org/x/tools/go/loader"
)

type Filter struct {
	includeStdLib bool
	permit        *regexp.Regexp
	deny          *regexp.Regexp
}

func NewFilter(includeStdLib bool, permit *regexp.Regexp, deny *regexp.Regexp) *Filter {
	return &Filter{
		includeStdLib: includeStdLib,
		permit:        permit,
		deny:          deny,
	}
}

type PkgGraph struct {
	PkgInfos map[*types.Package]*loader.PackageInfo
	Nodes    map[string]*PkgNode
	Filter   *Filter
}

func NewPkgGraph(filter *Filter, allPkgs map[*types.Package]*loader.PackageInfo) *PkgGraph {
	return &PkgGraph{
		Nodes:    make(map[string]*PkgNode),
		Filter:   filter,
		PkgInfos: allPkgs,
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

			if !p.Filter.permit.MatchString(cpath) {
				continue
			}

			if p.Filter.deny != nil {
				if p.Filter.deny.MatchString(cpath) {
					continue
				}
			}

			cNode, ok := p.Nodes[cpath]
			if !ok {
				cNode = NewPkgNode(c, p.PkgInfos[c].Files)
				p.Populate(cNode)
			}
			cNode.Parents = append(cNode.Parents, n)
			n.Children = append(n.Children, cNode)
		}
	}
}

func (g *PkgGraph) TotalFuncDecls() int {
	sum := 0
	for _, node := range g.Nodes {
		sum += node.TotalFuncDecls()
	}

	return sum
}

func (g *PkgGraph) CalcCallStats() {
	for _, node := range g.Nodes {
		node.CalcCallStats()
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

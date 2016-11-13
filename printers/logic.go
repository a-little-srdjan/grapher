package printers

import (
	"bytes"

	"github.com/a-little-srdjan/yagat/pkg_graph"
)

type PrologPrinter struct {
	graphPrinter
}

func NewPrologPrinter(graph *pkg_graph.PkgGraph) *PrologPrinter {
	p := &PrologPrinter{}
	p.graph = graph
	return p
}

func (p *PrologPrinter) WriteBuffer() *bytes.Buffer {
	p.buffer = new(bytes.Buffer)
	p.WriteParent()
	p.WriteDep()
	p.WriteAtoms()

	return p.buffer
}

func (p *PrologPrinter) WriteParent() {
	prologStmt(p.buffer, `parent(X, Y) :- direct_parent(X, Y), pkg(X), pkg(Y).`)
	prologStmt(p.buffer, `parent(X, Y) :- direct_parent(Z, Y), parent(X, Z).`)
}

func (p *PrologPrinter) WriteDep() {
	prologStmt(p.buffer, `dependency(X, Y) :- imports(X, Y), pkg(X), pkg(Y).`)
	prologStmt(p.buffer, `dependency(X, Y) :- imports(Z, Y), dependency(X, Z).`)
}

func (p *PrologPrinter) WriteAtoms() {
	for name, node := range p.graph.Nodes {
		prologStmt(p.buffer, atomStmt("pkg", name).String())
		for _, cnode := range node.Children {
			prologStmt(p.buffer, atomStmt("imports", name, cnode.Node.Path()).String())
		}
	}
}

func prologStmt(output *bytes.Buffer, stmt string) {
	output.WriteString(stmt)
	output.WriteString("\n")
}

func atomStmt(name string, params ...string) *bytes.Buffer {
	var b bytes.Buffer
	b.WriteString(name + "(")
	for _, p := range params {
		b.WriteString(stringConstant(p))
		b.WriteString(",")
	}
	b.Truncate(b.Len() - 1)
	b.WriteString(").")
	return &b
}

func stringConstant(constant string) string {
	return `"` + constant + `"`
}

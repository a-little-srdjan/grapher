package printers

import (
	"bytes"
	"strings"

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
	p.WriteSuppressWarnings()
	p.WriteNestedIDB()
	p.WriteDepIDB()
	p.WriteEDB()

	return p.buffer
}

func (p *PrologPrinter) WriteSuppressWarnings() {
	prologStmt(p.buffer, `:- discontiguous dir/1.`)
	prologStmt(p.buffer, `:- discontiguous direct_nested/2.`)
	prologStmt(p.buffer, `:- discontiguous pkg/1.`)
	prologStmt(p.buffer, `:- discontiguous imports/2.`)
}

func (p *PrologPrinter) WriteNestedIDB() {
	prologStmt(p.buffer, `nested(X, Y) :- direct_nested(X, Y), dir(X), dir(Y).`)
	prologStmt(p.buffer, `nested(X, Y) :- direct_nested(Z, Y), parent(X, Z).`)
	prologStmt(p.buffer, `pkg_dir(X) :- dir(X), pkg(X).`)
}

func (p *PrologPrinter) WriteDepIDB() {
	prologStmt(p.buffer, `dependency(X, Y) :- imports(X, Y), pkg(X), pkg(Y).`)
	prologStmt(p.buffer, `dependency(X, Y) :- imports(Z, Y), dependency(X, Z).`)
}

func (p *PrologPrinter) WriteEDB() {
	edbSet := make(map[string]struct{})

	for name, node := range p.graph.Nodes {
		prologStmt(p.buffer, atomStmt("pkg", name))

		for _, cnode := range node.Children {
			prologStmt(p.buffer, atomStmt("imports", name, cnode.Node.Path()))
		}

		nests := strings.Split(name, "/")
		p0 := nests[0]
		if !in(edbSet, p0) {
			edbSet[p0] = struct{}{}
			prologStmt(p.buffer, atomStmt("dir", p0))
		}

		for i := 1; i < len(nests); i++ {
			p1 := strings.Join([]string{p0, nests[i]}, "/")
			if !in(edbSet, p1) {
				edbSet[p1] = struct{}{}
				prologStmt(p.buffer, atomStmt("dir", p1))
				prologStmt(p.buffer, atomStmt("direct_nested", p0, p1))
			}
			p0 = p1
		}
	}
}

func in(set map[string]struct{}, e string) bool {
	_, ok := set[e]
	return ok
}

func prologStmt(output *bytes.Buffer, stmt string) {
	output.WriteString(stmt)
	output.WriteString("\n")
}

func atomStmt(name string, params ...string) string {
	var b bytes.Buffer
	b.WriteString(name + "(")
	for _, p := range params {
		b.WriteString(stringConstant(p))
		b.WriteString(",")
	}
	b.Truncate(b.Len() - 1)
	b.WriteString(").")
	return b.String()
}

func stringConstant(constant string) string {
	return `"` + constant + `"`
}

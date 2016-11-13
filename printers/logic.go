package printers

import (
	"bytes"
	"fmt"
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
	p.WriteLabelIDB()
	p.WriteConstraintsIDB()
	p.WriteEDB()

	return p.buffer
}

func (p *PrologPrinter) WriteSuppressWarnings() {
	prologStmt(p.buffer, `:- discontiguous dir/2.`)
	prologStmt(p.buffer, `:- discontiguous direct_nested/2.`)
	prologStmt(p.buffer, `:- discontiguous pkg/1.`)
	prologStmt(p.buffer, `:- discontiguous imports/2.`)
}

func (p *PrologPrinter) WriteNestedIDB() {
	prologStmt(p.buffer, `nested(X, Y) :- direct_nested(X, Y), dir(X, _), dir(Y, _).`)
	prologStmt(p.buffer, `nested(X, Y) :- direct_nested(Z, Y), nested(X, Z).`)
	prologStmt(p.buffer, `pkg_dir(X) :- dir(X), pkg(X).`)
}

func (p *PrologPrinter) WriteDepIDB() {
	prologStmt(p.buffer, `dependency(X, Y) :- imports(X, Y), pkg(X), pkg(Y).`)
	prologStmt(p.buffer, `dependency(X, Y) :- imports(Z, Y), dependency(X, Z).`)
}

func (p *PrologPrinter) WriteLabelIDB() {
	prologStmt(p.buffer, `p_label(M, D, Y) :- mark(M, Y), dir(Y, D).`)
	prologStmt(p.buffer, `p_label(M, D, Y) :- mark(M, Z), nested(Z, Y), dir(Z, D).`)
	prologStmt(p.buffer, `d_label(M, Y) :- p_label(M, D, Y), p_label(M2, D2, Y), M \== M2, D2 > D.`)
	prologStmt(p.buffer, `label(M, Y) :- p_label(M, _, Y), \+ d_label(M, Y).`)
}

func (p *PrologPrinter) WriteConstraintsIDB() {
	prologStmt(p.buffer, `violation(X, Y) :- dependency(X, Y), label(M, X), label(M2, Y), M \== M2, M < M2.`)
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
			prologStmt(p.buffer, atomStmt("dir", p0, 0))
		}

		for i := 1; i < len(nests); i++ {
			p1 := strings.Join([]string{p0, nests[i]}, "/")
			if !in(edbSet, p1) {
				edbSet[p1] = struct{}{}
				prologStmt(p.buffer, atomStmt("dir", p1, i))
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

func atomStmt(name string, params ...interface{}) string {
	var b bytes.Buffer
	b.WriteString(name + "(")
	for _, raw := range params {
		switch t := raw.(type) {
		default:
			b.WriteString("_")
		case string:
			b.WriteString(stringConstant(t))
		case int:
			b.WriteString(fmt.Sprintf("%d", t))
		}
		b.WriteString(",")
	}
	b.Truncate(b.Len() - 1)
	b.WriteString(").")
	return b.String()
}

func stringConstant(constant string) string {
	return `"` + constant + `"`
}

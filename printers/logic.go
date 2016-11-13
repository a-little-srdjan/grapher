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
	p.WriteNestedIDB()
	p.WriteDepIDB()
	p.WriteEDB()

	return p.buffer
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
		prologStmt(p.buffer, atomStmt("dir", p0))

		for i := 1; i < len(nests); i++ {
			p1 := strings.Join([]string{p0, nests[i]}, "/")
			_, ok := edbSet[p1]
			if !ok {
				edbSet[p1] = struct{}{}
				prologStmt(p.buffer, atomStmt("dir", p1))
				prologStmt(p.buffer, atomStmt("direct_nested", p0, p1))
			}
			p0 = p1
		}
	}
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

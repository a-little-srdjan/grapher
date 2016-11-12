package printers

import (
	"bytes"
	"log"
	"os"

	"github.com/a-little-srdjan/yagat/pkg_graph"
)

type PrologPrinter struct {
	graph *pkg_graph.PkgGraph
}

func NewPrologPrinter(graph *pkg_graph.PkgGraph) *PrologPrinter {
	return &PrologPrinter{
		graph: graph,
	}
}

func (p *PrologPrinter) Print(fileName string) {
	output, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer
	output.Write(p.GenerateProlog(&b).Bytes())
	output.Close()
}

func (p *PrologPrinter) GenerateProlog(output *bytes.Buffer) *bytes.Buffer {
	p.GenerateParent(output)
	p.GenerateDep(output)
	p.GenerateAtoms(output)
	return output
}

func (p *PrologPrinter) GenerateParent(output *bytes.Buffer) {
	prologStmt(output, `parent(X, Y) :- direct_parent(X, Y), pkg(X), pkg(Y).`)
	prologStmt(output, `parent(X, Y) :- direct_parent(Z, Y), parent(X, Z).`)
}

func (p *PrologPrinter) GenerateDep(output *bytes.Buffer) {
	prologStmt(output, `dependency(X, Y) :- imports(X, Y), pkg(X), pkg(Y).`)
	prologStmt(output, `dependency(X, Y) :- imports(Z, Y), dependency(X, Z).`)
}

func (p *PrologPrinter) GenerateAtoms(output *bytes.Buffer) {
	for name, node := range p.graph.Nodes {
		prologStmt(output, atomStmt("pkg", name).String())
		for _, cnode := range node.Children {
			prologStmt(output, atomStmt("imports", name, cnode.Node.Path()).String())
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

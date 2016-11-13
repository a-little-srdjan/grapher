package printers

import (
	"bytes"
	"log"
	"os"

	"github.com/a-little-srdjan/yagat/pkg_graph"
)

type graphPrinter struct {
	graph  *pkg_graph.PkgGraph
	buffer *bytes.Buffer
}

func Print(b *bytes.Buffer, fileName string) {
	output, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	output.Write(b.Bytes())
	output.Close()
}

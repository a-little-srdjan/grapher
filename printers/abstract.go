package printers

import (
	"bytes"
	"log"
	"os"

	"a-little-srdjan/grapher/model"
)

type graphPrinter struct {
	graph  *model.PkgGraph
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

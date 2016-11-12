package printers

import (
	"github.com/a-little-srdjan/yagat/pkg_graph"

	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
)

type GraphMLPrinter struct {
	nodeSize      float64
	nodeSizeBoost float64
	graph         *pkg_graph.PkgGraph
}

func NewGraphMLPrinter(graph *pkg_graph.PkgGraph, nodeSize, nodeSizeBoost float64) *GraphMLPrinter {
	return &GraphMLPrinter{
		nodeSize:      nodeSize,
		nodeSizeBoost: nodeSizeBoost,
		graph:         graph,
	}
}

func (p *GraphMLPrinter) Print(fileName string) {
	output, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer
	output.Write(p.GenerateGraphML(&b).Bytes())
	output.Close()
}

func (p *GraphMLPrinter) GenerateGraphML(output *bytes.Buffer) *bytes.Buffer {
	output.WriteString(graphMLPrefix)
	p.GenerateKeyElement(output)
	p.GenerateGraphElement(output)
	output.WriteString(graphMLSuffix)
	return output
}

func (p *GraphMLPrinter) GenerateKeyElement(output *bytes.Buffer) {
	output.WriteString(`<key for="node" id="d1" yfiles.type="nodegraphics"/>`)
}

func (p *GraphMLPrinter) GenerateGraphElement(output *bytes.Buffer) {
	output.WriteString(`<graph id="G" edgedefault="directed">`)

	graphTotalFuncDecls := p.graph.TotalFuncDecls()

	for name, node := range p.graph.Nodes {
		size := p.nodeSize + p.nodeSizeBoost*(float64(node.TotalFuncDecls())/float64(graphTotalFuncDecls))
		p.GenerateNodeElement(name, size, output)
	}

	id := 0
	for pname, pnode := range p.graph.Nodes {
		for _, cnode := range pnode.Children {
			p.GenerateEdgeElement(strconv.Itoa(id), pname, cnode.Node.Path(), output)
			id++
		}
	}

	output.WriteString(`</graph>`)
}

func (p *GraphMLPrinter) GenerateNodeElement(name string, size float64, output *bytes.Buffer) {
	output.WriteString(`<node id="` + name + `"><data key="d1"><y:ShapeNode>`)
	output.WriteString(`<y:Geometry height="` + fmt.Sprintf("%.2f", size) + `" width="` + fmt.Sprintf("%.2f", size) + `"/>`)
	output.WriteString(`<y:NodeLabel>` + name + `</y:NodeLabel>`)
	output.WriteString(`<y:Shape type="ellipse"/>`)
	output.WriteString(`</y:ShapeNode></data></node>`)
}

func (p *GraphMLPrinter) GenerateEdgeElement(id, source, target string, output *bytes.Buffer) {
	output.WriteString(`<edge id="` + id + `" source="` + source + `" target="` + target + `"/>`)
}

var graphMLPrefix = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<graphml
 xmlns="http://graphml.graphdrawing.org/xmlns"
 xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
 xmlns:y="http://www.yworks.com/xml/graphml"
 xmlns:yed="http://www.yworks.com/xml/yed/3"
 xsi:schemaLocation="http://graphml.graphdrawing.org/xmlns http://www.yworks.com/xml/schema/graphml/1.1/ygraphml.xsd">`

var graphMLSuffix = `</graphml>`

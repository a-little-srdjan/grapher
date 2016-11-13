package printers

import (
	"github.com/a-little-srdjan/yagat/pkg_graph"

	"bytes"
	"fmt"
	"strconv"
)

type GraphMLPrinter struct {
	graphPrinter
	nodeSize      float64
	nodeSizeBoost float64
}

func NewGraphMLPrinter(graph *pkg_graph.PkgGraph, nodeSize, nodeSizeBoost float64) *GraphMLPrinter {
	p := &GraphMLPrinter{
		nodeSize:      nodeSize,
		nodeSizeBoost: nodeSizeBoost,
	}
	p.graph = graph
	return p
}

func (p *GraphMLPrinter) WriteBuffer() *bytes.Buffer {
	p.buffer = new(bytes.Buffer)
	p.buffer.WriteString(graphMLPrefix)
	p.WriteKeyElement()
	p.WriteGraphElement()
	p.buffer.WriteString(graphMLSuffix)

	return p.buffer
}

func (p *GraphMLPrinter) WriteKeyElement() {
	p.buffer.WriteString(`<key for="node" id="d1" yfiles.type="nodegraphics"/>`)
}

func (p *GraphMLPrinter) WriteGraphElement() {
	p.buffer.WriteString(`<graph id="G" edgedefault="directed">`)

	graphTotalFuncDecls := p.graph.TotalFuncDecls()

	for name, node := range p.graph.Nodes {
		size := p.nodeSize + p.nodeSizeBoost*(float64(node.TotalFuncDecls())/float64(graphTotalFuncDecls))
		p.WriteNodeElement(name, size)
	}

	id := 0
	for pname, pnode := range p.graph.Nodes {
		for _, cnode := range pnode.Children {
			p.WriteEdgeElement(strconv.Itoa(id), pname, cnode.Node.Path())
			id++
		}
	}

	p.buffer.WriteString(`</graph>`)
}

func (p *GraphMLPrinter) WriteNodeElement(name string, size float64) {
	p.buffer.WriteString(`<node id="` + name + `"><data key="d1"><y:ShapeNode>`)
	p.buffer.WriteString(`<y:Geometry height="` + fmt.Sprintf("%.2f", size) + `" width="` + fmt.Sprintf("%.2f", size) + `"/>`)
	p.buffer.WriteString(`<y:NodeLabel>` + name + `</y:NodeLabel>`)
	p.buffer.WriteString(`<y:Shape type="ellipse"/>`)
	p.buffer.WriteString(`</y:ShapeNode></data></node>`)
}

func (p *GraphMLPrinter) WriteEdgeElement(id, source, target string) {
	p.buffer.WriteString(`<edge id="` + id + `" source="` + source + `" target="` + target + `"/>`)
}

var graphMLPrefix = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<graphml
 xmlns="http://graphml.graphdrawing.org/xmlns"
 xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
 xmlns:y="http://www.yworks.com/xml/graphml"
 xmlns:yed="http://www.yworks.com/xml/yed/3"
 xsi:schemaLocation="http://graphml.graphdrawing.org/xmlns http://www.yworks.com/xml/schema/graphml/1.1/ygraphml.xsd">`

var graphMLSuffix = `</graphml>`

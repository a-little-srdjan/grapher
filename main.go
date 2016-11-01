package main

import (
	"flag"
	"github.com/a-little-srdjan/yagat/pkg_graph"
	"golang.org/x/tools/go/loader"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var conf loader.Config

	pkgs := flag.String("pkgs", "fmt", "starting pkgs for the analysis")
	outputFile := flag.String("output", "a.graphml", "graphml output file")
	permit := flag.String("permit", "", "substring that has to be included in the pkg name")
	deny := flag.String("deny", "", "substraing that if contained removes the pkg from the graph")
	includeStdLib := flag.Bool("includeStdLib", false, "include std lib pkgs in the graph")

	flag.Parse()

	filter := pkg_graph.NewFilter(*includeStdLib, *permit, *deny)

	pkgList := strings.Split(*pkgs, ",")
	_, err := conf.FromArgs(pkgList, true)
	prog, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}

	g := pkg_graph.NewPkgGraph(filter)
	for _, v := range prog.Imported {
		n := pkg_graph.NewPkgNode(v.Pkg)
		g.Populate(n)
	}

	GenerateGraphML(g, *outputFile)
}

func GenerateGraphML(graph *pkg_graph.PkgGraph, fileName string) {
	output, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	output.WriteString(graphMLPrefix)
	output.WriteString(`<key for="node" id="d1" yfiles.type="nodegraphics"/>`)
	output.WriteString(`<graph id="G" edgedefault="directed">`)

	for name, _ := range graph.Nodes {
		output.WriteString(`<node id="` + name + `"><data key="d1"><y:ShapeNode>`)
		output.WriteString(`<y:NodeLabel>` + name + `</y:NodeLabel>`)
		output.WriteString(`</y:ShapeNode></data></node>`)
	}

	id := 0
	for _, node := range graph.Nodes {
		for _, child := range node.Children {
			output.WriteString(`<edge id="` + strconv.Itoa(id) + `" source="` + node.Node.Path() + `" target="` + child.Node.Path() + `"/>`)
			id++
		}
	}

	output.WriteString(`</graph>`)
	output.WriteString(graphMLSuffix)
	output.Close()

	log.Printf("Written %v nodes.\n", len(graph.Nodes))
	log.Printf("Written %v edges.\n", id)
}

var graphMLPrefix = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<graphml
 xmlns="http://graphml.graphdrawing.org/xmlns"
 xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
 xmlns:y="http://www.yworks.com/xml/graphml"
 xmlns:yed="http://www.yworks.com/xml/yed/3"
 xsi:schemaLocation="http://graphml.graphdrawing.org/xmlns http://www.yworks.com/xml/schema/graphml/1.1/ygraphml.xsd">`

var graphMLSuffix = `</graphml>`

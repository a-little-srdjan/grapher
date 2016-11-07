package main

import (
	"flag"
	"fmt"
	"github.com/a-little-srdjan/yagat/pkg_graph"
	"golang.org/x/tools/go/loader"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	var conf loader.Config
	var err error
	var pregexp *regexp.Regexp
	var dregexp *regexp.Regexp

	pkgs := flag.String("pkgs", "fmt", "root pkgs for the analysis")
	outputFile := flag.String("output", "a.graphml", "(yed) graphml output file")
	permit := flag.String("permit", "", "regex pattern that has to be included in the pkg name")
	deny := flag.String("deny", "", "regex pattern that must not be included in the pkg name")
	includeStdLib := flag.Bool("includeStdLib", false, "include std lib pkgs in the graph")

	flag.Parse()

	if *permit != "" {
		pregexp, err = regexp.Compile(*permit)
		if err != nil {
			log.Fatalf("Failed to compile permit regexp. %v", err)
		}
	}

	if *deny != "" {
		dregexp, err = regexp.Compile(*deny)
		if err != nil {
			log.Fatalf("Failed to compile deny regexp. %v", err)
		}
	}

	filter := pkg_graph.NewFilter(*includeStdLib, pregexp, dregexp)

	pkgList := strings.Split(*pkgs, ",")
	_, err = conf.FromArgs(pkgList, true)
	prog, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}

	g := pkg_graph.NewPkgGraph(filter, prog.AllPackages)
	for _, v := range prog.Imported {
		n := pkg_graph.NewPkgNode(v.Pkg, g.PkgInfos[v.Pkg].Files)
		g.Populate(n)
	}

	g.CalcCallStats()
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

	graphTotalFuncDecls := graph.TotalFuncDecls()
	for name, node := range graph.Nodes {
		size := 40.0 + 200*(float64(node.TotalFuncDecls())/float64(graphTotalFuncDecls))

		output.WriteString(`<node id="` + name + `"><data key="d1"><y:ShapeNode>`)
		output.WriteString(`<y:Geometry height="` + fmt.Sprintf("%.2f", size) + `" width="` + fmt.Sprintf("%.2f", size) + `"/>`)
		output.WriteString(`<y:NodeLabel>` + name + `</y:NodeLabel>`)
		output.WriteString(`<y:Shape type="ellipse"/>`)
		output.WriteString(`</y:ShapeNode></data></node>`)
	}

	id := 0
	for pname, pnode := range graph.Nodes {
		for _, cnode := range pnode.Children {
			output.WriteString(`<edge id="` + strconv.Itoa(id) + `" source="` + pname + `" target="` + cnode.Node.Path() + `"/>`)
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

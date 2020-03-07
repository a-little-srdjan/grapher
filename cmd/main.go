package main

import (
	"flag"
	"log"
	"regexp"
	"strings"

	"a-little-srdjan/grapher/model"
	"a-little-srdjan/grapher/printers"
	"golang.org/x/tools/go/loader"
)

func main() {
	var conf loader.Config
	var err error
	var pregexp *regexp.Regexp
	var dregexp *regexp.Regexp

	pkgs := flag.String("pkgs", "fmt", "root pkgs for the analysis")
	outputFile := flag.String("output", "a", "output file name for pl and grapml outputs")
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

	filter := model.NewFilter(*includeStdLib, pregexp, dregexp)

	pkgList := strings.Split(*pkgs, ",")
	_, err = conf.FromArgs(pkgList, true)
	prog, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}

	g := model.NewPkgGraph(filter, prog.AllPackages)
	for _, v := range prog.Imported {
		n := model.NewPkgNode(v.Pkg, g.PkgInfos[v.Pkg].Files)
		g.Populate(n)
	}

	g.CalcCallStats()
	printers.Print(printers.NewGraphMLPrinter(g, 40, 350).WriteBuffer(), *outputFile+".graphml")
	printers.Print(printers.NewPrologPrinter(g).WriteBuffer(), *outputFile+".pl")
}

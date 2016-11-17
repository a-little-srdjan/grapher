# Grapher
Yet Another Tool for Analyzing Go Packages

## Overview
Grapher is intended to aid code reviews and code analysis. Its focus is on
analyzing sizes of Go packages (in terms of functionality) and weights of
inter-package dependencies. Finally, all package links can be checked against constraints
using logic programming.

In detail, Grapher constructs the following graph:

1. Nodes are packages.
2. The size of a node represents the number of functions declared in that package, normalized accross all packages.
3. Directed edges represent package imports.
4. Edge weights represent the number of times the imported package has been used in variable definitions and function calls. 

The output consists of two declarative specifications:

1. [GraphML](http://graphml.graphdrawing.org/) specification
2. [Prolog](https://en.wikipedia.org/wiki/Prolog) program. 

These outputs form the basis for graph analysis and constraint checking over the package graph.

The GraphML spec can be examined by [yEd](http://www.yworks.com/products/yed). We can apply different grouping algorithms
to find package clusters and outliers. This can help confirm/refute different hypothesis that we may have about our code base.
For example, lack of distinct clusters indicate a code base with no structure and layering, a high number of outliers may indicate
a need to combine different packages, etc. Finally, we can also look at two centrality measures: _edge_ and _betweeness_, to 
find influential nodes and confirm whether the code base should in fact have super nodes.

The Prolog program (see [Swi-Prolog](http://www.swi-prolog.org/), an easy-to-install interpreter) contains the declarative specification
of the package dependencies, and the directory structure. We can then add constraints to check whether the logical separation of packages within
different layers (such as _endpoints_, _services_, _frameworks_, etc) is broken with the respect to the dependency structure. Clearly, this
analysis only applies if the code base has some logical groupings amongst packages, and the code is not alloed to flow from _high_ packages
(e.g. endpoints) to _low_ packages (e.g. crypto).

## Usage and Examples
Flags:
* _pkgs_ : root pkgs for the analysis
* _outputFile_ 
* _permit_ : regex pattern that has to be part of the pkg name to have the pkg included
* _deny_ : regex pattern that must not be part of the pkg name to have the pkg included
* _includeStdLib_ : include std lib pkgs in the graph

1. build the tool 
2. `grapher -deny=vendor -pkgs=code.wirelessregistry.com/data/readers/queries -output=depgraph`

## TODO
1. Increase the edge weights with method calls. That is, currently, expressions such as
`varName.Method()` are not taken into account for edge weights.
 
## Dependencies
1. [go loader tool](https://godoc.org/golang.org/x/tools/go/loader)

## Related Work
* [goviz](https://github.com/hirokidaichi/goviz)
* [Visualising dependencies](https://dave.cheney.net/2014/11/21/visualising-dependencies) by Dave Cheney
* [Building the simplest Go static analysis tool](https://blog.cloudflare.com/building-the-simplest-go-static-analysis-tool/) by Filippo Valsorda

## Feedback
Please send all comments and suggestion to _srdjan.marinovic@gmail.com_

# Gapher
Yet Another Go Analysis Tool (_feature set not complete_)

## Overview
Gapher takes as input a set of Go packages and produces a 
GraphML and a Prolog specification. These specs can then be inputs
to graph analysise:

1. GraphML spec can be examined by [yEd](http://www.yworks.com/products/yed).
The generated nodes' sizes indicate how much functionality is contained in each 
package (in terms of declared functions).  We can use yEd to look at two centrality 
measures: _edge_ and _betweeness_, to find influential packages.

2. Prolog spec can be run by [Swi-Prolog](http://www.swi-prolog.org/). 
Using Prolog, we can add constraints over our code base and check whether we have code flowing from
high-level packages to low-level packages.

## Usage
Flags:
* _pkgs_ : root pkgs for the analysis
* _outputFile_ 
* _permit_ : regex pattern that has to be part of the pkg name to have the pkg included
* _deny_ : regex pattern that must not be part of the pkg name to have the pkg included
* _includeStdLib_ : include std lib pkgs in the graph

1. build the tool 
2. `grapher -deny=vendor -pkgs=code.wirelessregistry.com/data/readers/queries -output=depgraph`

## Examples
...coming...

## Dependencies
1. [go loader tool](https://godoc.org/golang.org/x/tools/go/loader)

## Related Work
* [goviz](https://github.com/hirokidaichi/goviz)
* [Visualising dependencies](https://dave.cheney.net/2014/11/21/visualising-dependencies) by Dave Cheney
* [Building the simplest Go static analysis tool](https://blog.cloudflare.com/building-the-simplest-go-static-analysis-tool/) by Filippo Valsorda


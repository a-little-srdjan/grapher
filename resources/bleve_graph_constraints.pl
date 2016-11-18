
violation("bleve", Y) :- dependency(Y, "github.com/blevesearch/bleve").
violation("searching", Y) :- dependency("github.com/blevesearch/bleve/search", Y), nested("github.com/blevesearch/bleve/index", Y).
violation("indexing", Y) :- dependency("github.com/blevesearch/bleve/index", Y), nested("github.com/blevesearch/bleve/analysis", Y).

exception("searching", "github.com/blevesearch/bleve/index/store").

broken(X, Y) :- violation(X, Y), \+ exception(X, Y). 

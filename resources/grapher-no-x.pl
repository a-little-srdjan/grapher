:- discontiguous dir/2.
:- discontiguous direct_nested/2.
:- discontiguous pkg/1.
:- discontiguous imports/2.
nested(X, Y) :- direct_nested(X, Y), dir(X, _), dir(Y, _).
nested(X, Y) :- direct_nested(Z, Y), nested(X, Z).
pkg_dir(X) :- dir(X), pkg(X).
dependency(X, Y) :- imports(X, Y), pkg(X), pkg(Y).
dependency(X, Y) :- imports(Z, Y), dependency(X, Z).
p_label(M, D, Y) :- mark(M, Y), dir(Y, D).
p_label(M, D, Y) :- mark(M, Z), nested(Z, Y), dir(Z, D).
d_label(M, Y) :- p_label(M, D, Y), p_label(M2, D2, Y), M \== M2, D2 > D.
label(M, Y) :- p_label(M, _, Y), \+ d_label(M, Y).
violation(X, Y) :- dependency(X, Y), label(M, X), label(M2, Y), M \== M2, M < M2.
pkg("github.com/a-little-srdjan/grapher/printers").
imports("github.com/a-little-srdjan/grapher/printers","github.com/a-little-srdjan/grapher/pkg_graph").
dir("github.com",0).
dir("github.com/a-little-srdjan",1).
direct_nested("github.com","github.com/a-little-srdjan").
dir("github.com/a-little-srdjan/grapher",2).
direct_nested("github.com/a-little-srdjan","github.com/a-little-srdjan/grapher").
dir("github.com/a-little-srdjan/grapher/printers",3).
direct_nested("github.com/a-little-srdjan/grapher","github.com/a-little-srdjan/grapher/printers").
pkg("github.com/a-little-srdjan/grapher").
imports("github.com/a-little-srdjan/grapher","github.com/a-little-srdjan/grapher/pkg_graph").
imports("github.com/a-little-srdjan/grapher","github.com/a-little-srdjan/grapher/printers").
pkg("github.com/a-little-srdjan/grapher/pkg_graph").
dir("github.com/a-little-srdjan/grapher/pkg_graph",3).
direct_nested("github.com/a-little-srdjan/grapher","github.com/a-little-srdjan/grapher/pkg_graph").

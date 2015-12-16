package link8

import (
	"fmt"
)

type tracer struct {
	pkgs map[string]*Pkg
	hits map[string]map[string]bool
}

func newTracer(pkgs map[string]*Pkg) *tracer {
	ret := &tracer{
		pkgs: pkgs,
		hits: make(map[string]map[string]bool),
	}
	for path := range pkgs {
		ret.hits[path] = make(map[string]bool)
	}

	return ret
}

func (t *tracer) hit(sym *PkgSym) bool {
	p, found := t.hits[sym.Pkg]
	if !found {
		panic("path not found")
	}
	ret := p[sym.Sym]
	p[sym.Sym] = true
	return ret
}

// traceUsed traces symbols/objects that are used.
// only these objects need to be linked into the final result.
func traceUsed(pkgs map[string]*Pkg, roots []*PkgSym) []*PkgSym {
	t := newTracer(pkgs)
	cur := roots

	var next []*PkgSym
	var ret []*PkgSym

	addLink := func(link *link) {
		if pkgs[link.Pkg] == nil {
			panic(fmt.Errorf(
				"package %q missing", link.Pkg,
			))
		}

		if t.hit(link.PkgSym) {
			return
		}

		next = append(next, link.PkgSym)
	}

	// BFS traverse all the symbols used by the symbol
	for len(cur) > 0 {
		ret = append(ret, cur...)
		for _, ps := range cur {

			pkg := pkgs[ps.Pkg]
			s := pkg.SymbolByName(ps.Sym)
			switch s.Type {
			case SymFunc:
				f := pkg.Func(ps.Sym)
				for _, link := range f.links {
					addLink(link)
				}
			case SymVar:
				v := pkg.Var(ps.Sym)
				for _, link := range v.links {
					addLink(link)
				}
			}
		}

		cur = next
		next = nil
	}

	return ret
}

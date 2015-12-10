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

func (t *tracer) hit(pkg *Pkg, sym string) bool {
	p, found := t.hits[pkg.path]
	if !found {
		panic("path not found")
	}
	ret := p[sym]
	p[sym] = true
	return ret
}

// traceUsed traces symbols/objects that are used.
// only these objects need to be linked into the final result.
func traceUsed(pkgs map[string]*Pkg, path string, roots []string) []pkgSym {
	p := pkgs[path]
	t := newTracer(pkgs)

	var cur []pkgSym
	for _, name := range roots {
		cur = append(cur, pkgSym{p, name})
	}

	var next []pkgSym
	var ret []pkgSym

	addLink := func(ps pkgSym, link *link) {
		pkg := pkgs[link.pkg]
		if pkg == nil {
			panic(fmt.Errorf(
				"package %q missing in %q", link.pkg, ps.pkg.path,
			))
		}

		if t.hit(pkg, link.sym) {
			return
		}

		item := pkgSym{pkg, link.sym}
		next = append(next, item)
	}

	// BFS traverse all the symbols used by the symbol
	for len(cur) > 0 {
		for _, ps := range cur {
			ret = append(ret, ps)

			typ := ps.Type()
			switch typ {
			case SymFunc:
				f := ps.Func()
				for _, link := range f.links {
					addLink(ps, link)
				}
			case SymVar:
				v := ps.Var()
				for _, link := range v.links {
					addLink(ps, link)
				}
			}
		}

		cur = next
		next = nil
	}

	return ret
}

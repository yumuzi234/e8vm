package link8

import (
	"fmt"
)

type tracer struct {
	lnk  *linker
	hits map[string]map[string]bool
}

func newTracer(lnk *linker) *tracer {
	ret := new(tracer)
	ret.lnk = lnk
	ret.hits = make(map[string]map[string]bool)
	for path := range lnk.pkgs {
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
func traceUsed(lnk *linker, p *Pkg, roots []string) []pkgSym {
	t := newTracer(lnk)

	var cur []pkgSym
	for _, name := range roots {
		cur = append(cur, pkgSym{p, name})
	}

	var next []pkgSym
	var ret []pkgSym

	addLink := func(ps pkgSym, link *link) {
		pkg := lnk.pkg(link.pkg)
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

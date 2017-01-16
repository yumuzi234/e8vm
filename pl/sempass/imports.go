package sempass

import (
	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
	"shanhu.io/smlvm/syms"
)

func buildImports(
	b *builder, f *ast.File, imps map[string]*builds.Package,
) []*syms.Symbol {
	if f.Imports == nil {
		return nil
	}

	var ret []*syms.Symbol
	for _, d := range f.Imports.Decls {
		_, as, e := ast.ImportPathAs(d)
		if e != nil {
			b.Errorf(d.Path.Pos, "invalid import path")
			continue
		}

		p := imps[as]
		if p == nil {
			b.Errorf(d.Path.Pos, "package %s missing", as)
			continue
		}

		pos := ast.ImportPos(d)
		if !(p.Lang == "asm8" || p.Lang == "g8") {
			b.Errorf(pos, "cannot import %q pacakge %q",
				p.Lang, as,
			)
			continue
		}

		t := &types.Pkg{As: as, Lang: p.Lang, Syms: p.Symbols}
		sym := syms.Make(b.path, as, tast.SymImport, nil, t, pos)
		conflict := b.scope.Declare(sym)
		if conflict != nil {
			b.Errorf(pos, "%s already imported", as)
			continue
		}

		ret = append(ret, sym)
	}
	return ret
}

func checkUnusedImports(b *builder, imports []*syms.Symbol) {
	for _, s := range imports {
		name := s.Name()
		if name == "_" {
			continue
		}
		if !s.Used {
			b.Errorf(s.Pos, "unsed import: %q", name)
		}
	}
}

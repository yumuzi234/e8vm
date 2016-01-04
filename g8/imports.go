package g8

import (
	"io"

	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/lex8"
)

type importDecl struct {
	as   string
	path string
	pos  *lex8.Pos
}

func listImport(
	f string, rc io.ReadCloser, imp build8.Importer, golike bool,
) []*lex8.Error {
	fast, _, es := parse.File(f, rc, golike)
	if es != nil {
		return es
	}

	if fast.Imports == nil {
		return nil
	}

	m := make(map[string]*importDecl)
	log := lex8.NewErrorList()

	for _, d := range fast.Imports.Decls {
		p, as, e := ast.ImportPathAs(d)
		if e != nil {
			log.Errorf(d.Path.Pos, "invalid path string %s", d.Path.Lit)
			continue
		}

		pos := ast.ImportPos(d)
		if other, found := m[as]; found {
			log.Errorf(pos, "%s already imported", as)
			log.Errorf(other.pos, "  previously imported here")
			continue
		}

		m[as] = &importDecl{as: as, path: p, pos: pos}
	}

	if errs := log.Errs(); errs != nil {
		return errs
	}

	for as, d := range m {
		imp.Import(as, d.path, d.pos)
	}

	return nil
}

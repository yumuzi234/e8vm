package pl

import (
	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/parse"
)

type importDecl struct {
	as   string
	path string
	pos  *lexing.Pos
}

func listImport(
	f string, o builds.FileOpener, imp builds.Importer, golike bool,
) []*lexing.Error {
	rc, err := o.Open()
	if err != nil {
		return lexing.SingleErr(err)
	}
	fast, _, es := parse.File(f, rc, golike)
	if es != nil {
		return es
	}

	if fast.Imports == nil {
		return nil
	}

	m := make(map[string]*importDecl)
	log := lexing.NewErrorList()

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

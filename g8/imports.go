package g8

import (
	"io"
	"path"
	"strconv"

	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

type importDecl struct {
	as   string
	path string
	pos  *lex8.Pos
}

func importPos(d *ast.ImportDecl) *lex8.Pos {
	if d.As == nil {
		return d.Path.Pos
	}
	return d.As.Pos
}

func importPathAs(d *ast.ImportDecl) (p, as string, err error) {
	p, err = strconv.Unquote(d.Path.Lit)
	if err != nil {
		return "", "", err
	}

	if d.As == nil {
		return p, path.Base(p), nil
	}
	return p, d.As.Lit, nil
}

func listImport(
	f string, rc io.ReadCloser, imp build8.Importer, golike bool,
) []*lex8.Error {
	ast, es := parse.File(f, rc, golike)
	if es != nil {
		return es
	}

	if ast.Imports == nil {
		return nil
	}

	m := make(map[string]*importDecl)
	log := lex8.NewErrorList()

	for _, d := range ast.Imports.Decls {
		p, as, e := importPathAs(d)
		if e != nil {
			log.Errorf(d.Path.Pos, "invalid path string %s", d.Path.Lit)
			continue
		}

		pos := importPos(d)
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

func declareImports(b *builder, f *ast.File, pinfo *build8.PkgInfo) {
	if f.Imports == nil {
		return
	}

	for _, d := range f.Imports.Decls {
		_, as, e := importPathAs(d)
		if e != nil {
			b.Errorf(d.Path.Pos, "invalid path")
			continue
		}

		imported := pinfo.Import[as]
		if imported == nil {
			b.Errorf(d.Path.Pos, "package %s missing", as)
			continue
		}

		p := imported.Package
		if p.Lang != "g8" {
			// TODO(h8liu): issue #24, import assembly
			b.Errorf(d.Path.Pos, "not a G language package")
			continue
		}

		pos := importPos(d)
		ref := newRef(&types.Pkg{as, p.Symbols}, nil)
		obj := &objImport{ref}
		pre := b.scope.Declare(sym8.Make(b.symPkg, as, symImport, obj, pos))
		if pre != nil {
			b.Errorf(pos, "%s already declared", as)
			continue
		}
	}
}

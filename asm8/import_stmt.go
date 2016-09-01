package asm8

import (
	"path"
	"strconv"

	"e8vm.io/e8vm/asm8/ast"
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/lexing"
)

type importStmt struct {
	*ast.ImportStmt

	as   string
	path string

	pkg *build8.Package
}

func importPos(i *ast.ImportStmt) *lexing.Pos {
	if i.As == nil {
		return i.Path.Pos
	}
	return i.As.Pos
}

func (s *importStmt) pos() *lexing.Pos {
	return importPos(s.ImportStmt)
}

func resolveImportStmt(log lexing.Logger, imp *ast.ImportStmt) *importStmt {
	ret := new(importStmt)
	ret.ImportStmt = imp

	s, e := strconv.Unquote(imp.Path.Lit)
	if e != nil {
		log.Errorf(imp.Path.Pos, "invalid string %s", imp.Path.Lit)
		return nil
	}

	ret.path = s

	if imp.As != nil {
		ret.as = imp.As.Lit
	} else {
		ret.as = path.Base(ret.path)
	}

	return ret
}

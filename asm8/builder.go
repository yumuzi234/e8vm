package asm8

import (
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

// Builder manipulates an AST, checks its syntax, and builds the assembly
type builder struct {
	*lex8.ErrorList
	scope *sym8.Scope
	path  string

	curPkg *lib

	imports map[string]string
	pkgUsed map[string]struct{}
	inits   []*build8.PkgSym
}

func newBuilder(path string) *builder {
	return &builder{
		ErrorList: lex8.NewErrorList(),
		scope:     sym8.NewScope(),
		path:      path,
		imports:   make(map[string]string),
		pkgUsed:   make(map[string]struct{}),
	}
}

// Errs returns the building errors.
func (b *builder) Errs() []*lex8.Error {
	return b.ErrorList.Errs()
}

func (b *builder) importPkg(path, as string) {
	_, found := b.imports[as]
	if found {
		panic("redeclare")
	}

	b.imports[as] = path
}

func (b *builder) pkgPath(as string) string {
	ret, found := b.imports[as]
	if !found {
		panic("not found")
	}
	return ret
}

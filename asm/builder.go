package asm

import (
	"e8vm.io/e8vm/lexing"
	"e8vm.io/e8vm/syms"
)

// Builder manipulates an AST, checks its syntax, and builds the assembly
type builder struct {
	*lexing.ErrorList
	scope *syms.Scope
	path  string

	curPkg *lib

	imports map[string]string
	pkgUsed map[string]struct{}
}

func newBuilder(path string) *builder {
	return &builder{
		ErrorList: lexing.NewErrorList(),
		scope:     syms.NewScope(),
		path:      path,
		imports:   make(map[string]string),
		pkgUsed:   make(map[string]struct{}),
	}
}

// Errs returns the building errors.
func (b *builder) Errs() []*lexing.Error {
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

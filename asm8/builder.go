package asm8

import (
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

// Builder manipulates an AST, checks its syntax, and builds the assembly
type builder struct {
	*lex8.ErrorList
	scope  *sym8.Scope
	symPkg *sym8.Pkg

	curPkg *lib

	indices map[string]uint32
	pkgUsed map[string]struct{}
}

func newBuilder() *builder {
	ret := new(builder)
	ret.ErrorList = lex8.NewErrorList()
	ret.scope = sym8.NewScope()
	ret.symPkg = new(sym8.Pkg)
	ret.indices = make(map[string]uint32)
	ret.pkgUsed = make(map[string]struct{})

	return ret
}

// Errs returns the building errors.
func (b *builder) Errs() []*lex8.Error {
	return b.ErrorList.Errs()
}

func (b *builder) index(name string, index uint32) {
	_, found := b.indices[name]
	if found {
		panic("redeclare")
	}

	b.indices[name] = index
}

func (b *builder) getIndex(name string) uint32 {
	ret, found := b.indices[name]
	if !found {
		panic("not found")
	}
	return ret
}

package asm8

import (
	"e8vm.io/e8vm/link8"
	"e8vm.io/e8vm/sym8"
)

// Lib is the compiler output of a package
// it contains the package for linking,
// and also the symbols for importing
type lib struct {
	*link8.Pkg
	symbols map[string]*sym8.Symbol
}

// NewPkgObj creates a new package compile object
func newLib(p string) *lib {
	ret := new(lib)
	ret.Pkg = link8.NewPkg(p)
	ret.symbols = make(map[string]*sym8.Symbol)
	return ret
}

func (lib *lib) declare(s *sym8.Symbol) {
	_, found := lib.symbols[s.Name()]
	if found {
		panic("redeclare")
	}
	lib.symbols[s.Name()] = s

	switch s.Type {
	case SymConst:
		panic("todo")
	case SymFunc:
		lib.Pkg.DeclareFunc(s.Name())
	case SymVar:
		lib.Pkg.DeclareVar(s.Name())
	default:
		panic("declare with invalid sym type")
	}
}

// query returns the symbol declared by name and its symbol index
// if the symbol is a function or variable. It returns nil, 0 when
// the symbol of name is not found.
func (lib *lib) query(name string) *sym8.Symbol {
	ret, found := lib.symbols[name]
	if !found {
		return nil
	}

	switch ret.Type {
	case SymConst:
		return ret
	case SymFunc, SymVar:
		s := lib.Pkg.SymbolByName(name)
		if s == nil {
			panic("symbol missing")
		}
		return ret
	default:
		panic("bug")
	}
}

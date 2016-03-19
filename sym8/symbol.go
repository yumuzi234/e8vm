package sym8

import (
	"e8vm.io/e8vm/lex8"
)

// Symbol is a data structure for saving a symbol.
type Symbol struct {
	pkg  string
	name string

	Type    int
	Obj     interface{}
	ObjType interface{}
	Pos     *lex8.Pos
	Used    bool
}

// Name returns the symbol name.
// This name is immutable for its used for indexing in the tables.
func (s *Symbol) Name() string { return s.name }

// Pkg returns the package token of the symbol.
func (s *Symbol) Pkg() string { return s.pkg }

// Clone returns a clone of the symbol which item is the given one.
func (s *Symbol) Clone(item interface{}) *Symbol {
	ret := *s
	ret.Obj = item
	return &ret
}

// Make creates a new symbol
func Make(
	pkg, name string, t int,
	obj, objType interface{},
	pos *lex8.Pos,
) *Symbol {
	return &Symbol{pkg, name, t, obj, objType, pos, false}
}

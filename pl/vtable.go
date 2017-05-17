package pl

import (
	"shanhu.io/smlvm/pl/types"
	"shanhu.io/smlvm/syms"
)

// vTable is the virtual table to implement the interface
type vTable struct {
	funcs []string
	// will change it to *ref
	implementMap map[*types.Struct][]*syms.Symbol
}

func newTable(i *types.Interface) *vTable {
	size := len(i.Syms.List())
	m := make([]string, size)
	for n, sym := range i.Syms.List() {
		m[n] = sym.Name()
	}
	return &vTable{
		funcs:        m,
		implementMap: make(map[*types.Struct][]*syms.Symbol),
	}
}

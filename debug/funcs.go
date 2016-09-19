package debug

import (
	"shanhu.io/smlvm/lexing"
)

// Funcs saves all the debug symbols for all functions.
type Funcs struct {
	funcs map[string]*Func
}

// NewFuncs creates a new function symbol collection.
func NewFuncs() *Funcs {
	return &Funcs{make(map[string]*Func)}
}

func symKey(pkg, name string) string {
	return pkg + "." + name
}

// Add adds a function into the debug table.
func (fs *Funcs) Add(pkg, name string, pos *lexing.Pos, frameSize uint32) {
	key := symKey(pkg, name)
	if _, found := fs.funcs[key]; found {
		panic("bug")
	}

	fs.funcs[key] = &Func{Frame: frameSize, Pos: pos}
}

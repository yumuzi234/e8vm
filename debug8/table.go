package debug8

import (
	"encoding/json"

	"e8vm.io/e8vm/lex8"
)

// Table is a debug table that save symbol information.
type Table struct {
	Funcs map[string]*Func
}

// NewTable creates a new debug table.
func NewTable() *Table {
	return &Table{
		Funcs: make(map[string]*Func),
	}
}

// Marshal marshals the debug table out.
func (t *Table) Marshal() []byte {
	bs, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return bs
}

func symKey(pkg, name string) string {
	return pkg + "." + name
}

// Func returns the function entry of a particular symbol.
func (t *Table) Func(pkg, name string) *Func {
	key := symKey(pkg, name)
	f, found := t.Funcs[key]
	if !found {
		f = new(Func)
		t.Funcs[key] = f
	}
	return f
}

// LinkFunc saves the function linking debug information.
func (t *Table) LinkFunc(pkg, name string, addr, size uint32) {
	f := t.Func(pkg, name)
	f.Start = addr
	f.Size = size
}

// GenFunc saves the function generation debug information.
func (t *Table) GenFunc(pkg, name string, pos *lex8.Pos, frameSize uint32) {
	f := t.Func(pkg, name)
	f.Frame = frameSize
	f.Pos = pos
}

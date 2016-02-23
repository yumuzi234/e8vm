package debug8

import (
	"fmt"
	"io"

	"encoding/json"
)

// Table is a debug table that save symbol information.
type Table struct {
	Funcs map[string]*Func
}

// NewTable creates a new debug table.
func NewTable() *Table {
	return &Table{make(map[string]*Func)}
}

// UnmarshalTable unmarshals a debug table.
func UnmarshalTable(bs []byte) (*Table, error) {
	t := NewTable()
	err := json.Unmarshal(bs, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Marshal marshals the debug table out.
func (t *Table) Marshal() []byte {
	bs, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return bs
}

// LinkFunc saves the function linking debug information.
func (t *Table) LinkFunc(fs *Funcs, pkg, name string, addr, size uint32) {
	key := symKey(pkg, name)
	f, found := fs.funcs[key]
	if !found {
		f = new(Func)
	}

	t.Funcs[key] = f
	f.Start = addr
	f.Size = size
}

// PrintTo prints the table to an output stream.
func (t *Table) PrintTo(w io.Writer) error {
	for name, f := range t.Funcs {
		_, err := fmt.Fprintln(w, f.String(name))
		if err != nil {
			return err
		}
	}
	return nil
}

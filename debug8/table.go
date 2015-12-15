package debug8

import (
	"bytes"
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

func funcString(name string, f *Func) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%8x +%4d: ", f.Start, f.Size)
	if f.Frame > 0 {
		fmt.Fprintf(buf, "frame=%d - ", f.Frame)
	}
	fmt.Fprintf(buf, "%s", name)
	if f.Pos != nil {
		fmt.Fprintf(buf, " // %s", f.Pos)
	}
	return buf.String()
}

// PrintTo prints the table to an output stream.
func (t *Table) PrintTo(w io.Writer) error {
	for name, f := range t.Funcs {
		_, err := fmt.Fprintln(w, funcString(name, f))
		if err != nil {
			return err
		}
	}
	return nil
}

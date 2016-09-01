package types

import (
	"e8vm.io/e8vm/syms"
)

// Field is a named field in a struct
type Field struct {
	Name string
	T

	offset int32
}

// Offset returns the offset of a field in a struct.
// It will return the valid number only after the field
// has been added into a struct.
func (f *Field) Offset() int32 { return f.offset }

// Struct is the type of a structure
type Struct struct {
	Syms *syms.Table

	name         string
	size         int32
	regSizeAlign bool
}

// NewStruct constructs a new struct type.
func NewStruct(name string) *Struct {
	ret := new(Struct)
	ret.name = name
	ret.Syms = syms.NewTable()

	return ret
}

// AddField adds a field into the struct,
// assigns offset to the field.
func (t *Struct) AddField(f *Field) {
	fsize := f.T.Size()
	if f.T.RegSizeAlign() {
		t.size = RegSizeAlignUp(t.size)
		t.regSizeAlign = true
	}
	f.offset = t.size
	t.size += fsize
}

// Size returns the overall size of the structure type
func (t *Struct) Size() int32 { return t.size }

// String returns the name of the structure type
func (t *Struct) String() string { return t.name }

// RegSizeAlign returns true when at least one field in the struct
// is word aligned.
func (t *Struct) RegSizeAlign() bool { return t.regSizeAlign }

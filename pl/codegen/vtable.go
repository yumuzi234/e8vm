package codegen

import "fmt"

type vTable struct {
	pkg string

	// interface name: pkg.interface
	iName string
	// struct name:pkg.structName
	// sName can used for reflection
	sName string
	// implements will be sorted by func names
	implements []*FuncSym
}

func newVtable(i, s string) *vTable {
	return &vTable{
		iName:      i,
		sName:      s,
		implements: make([]*FuncSym, 0),
	}
}

func (t *vTable) String() string {
	return fmt.Sprintf("<vTable %s %s>", t.iName, t.sName)
}

func (t *vTable) RegSizeAlign() bool { return true }

func (t *vTable) Size() int32 {
	return regSize * 2
}

func (t *vTable) fill(fs []*FuncSym) {
	for _, f := range fs {
		t.implements = append(t.implements, f)
	}
}

type vTablePool struct {
	pkg     string
	vTables []*vTable
}

func newVTablePool(pkg string) *vTablePool {
	return &vTablePool{
		pkg:     pkg,
		vTables: make([]*vTable, 0),
	}
}

func (p *vTablePool) addTable(i, s string, funcs []*FuncSym) Ref {
	t := newVtable(i, s)
	t.fill(funcs)
	p.vTables = append(p.vTables, t)
	return t
}

// declare is not needed here

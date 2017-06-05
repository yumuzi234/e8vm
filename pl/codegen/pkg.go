package codegen

import (
	"fmt"

	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/link"
)

// Pkg is a package in its intermediate representation.
type Pkg struct {
	lib *link.Pkg

	path string

	funcs   []*Func
	vars    []*HeapSym
	tests   *testList
	strPool *strPool
	datPool *datPool

	// helper functions required for generating
	g *gener
}

// NewPkg creates a package with a particular path name.
func NewPkg(path string) *Pkg {
	return &Pkg{
		path:    path,
		lib:     link.NewPkg(path),
		strPool: newStrPool(path),
		datPool: newDatPool(path),
		g:       newGener(),
	}
}

// NewFunc creates a new function for the package.
func (p *Pkg) NewFunc(name string, pos *lexing.Pos, sig *FuncSig) *Func {
	ret := newFunc(p.path, name, pos, sig)
	p.lib.DeclareFunc(ret.name)
	p.funcs = append(p.funcs, ret)
	return ret
}

// NewGlobalVar creates a new global variable reference.
func (p *Pkg) NewGlobalVar(
	size int32, name string, u8, regSizeAlign bool,
) Ref {
	ret := NewHeapSym(p.path, name, size, u8, regSizeAlign)
	p.lib.DeclareVar(ret.name)
	p.vars = append(p.vars, ret)
	return ret
}

// NewTestList creates a global variable of a list of function symbols.
func (p *Pkg) NewTestList(name string, funcs []*Func) Ref {
	if len(funcs) > 1000000 {
		panic("too many test cases")
	}
	if p.tests != nil {
		panic("tests already built")
	}

	ret := newTestList(p.path, name, funcs)
	p.lib.DeclareVar(ret.name)
	p.tests = ret

	return ret
}

// HookBuiltin uses the builtin package that provides neccessary
// builtin functions for IR generation
func (p *Pkg) HookBuiltin(pkg *link.Pkg) error {
	var err error
	se := func(e error) {
		if err != nil {
			err = e
		}
	}

	o := func(f string) *FuncSym {
		sym := pkg.SymbolByName(f)
		if sym == nil {
			se(fmt.Errorf("%s missing in builtin", f))
		} else if sym.Type != link.SymFunc {
			se(fmt.Errorf("%s is not a function", f))
		}

		return &FuncSym{pkg: pkg.Path(), name: f}
	}

	p.g.memCopy = o("MemCopy")
	p.g.memClear = o("MemClear")

	return err
}

// NewString adds a new string constant and returns its reference.
func (p *Pkg) NewString(s string) Ref {
	return p.strPool.addString(s)
}

// NewHeapDat adds a byte array of heap static data.
func (p *Pkg) NewHeapDat(bs []byte, unit int32, regSizeAlign bool) Ref {
	return p.datPool.addBytes(bs, unit, regSizeAlign)
}

// NewVtable adds the vtable of the package.
func (p *Pkg) NewVtable(funcs []FuncSym) Ref {
	return p.datPool.addVtable(funcs)
}

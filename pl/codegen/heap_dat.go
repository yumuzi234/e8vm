package codegen

import (
	"fmt"

	"shanhu.io/smlvm/link"
)

type heapDat struct {
	pkg, name    string
	dat          interface{}
	regSizeAlign bool
	size         int32
	n            int32
}

func (s *heapDat) String() string {
	return fmt.Sprintf("<dat %s>", s.name)
}

func (s *heapDat) RegSizeAlign() bool { return s.regSizeAlign }

func (s *heapDat) Size() int32 { return s.size }

type datPool struct {
	pkg string
	dat []*heapDat
}

func newDatPool(pkg string) *datPool {
	return &datPool{
		pkg: pkg,
	}
}

func (p *datPool) addBytes(bs []byte, unit int32, regSizeAlign bool) *heapDat {
	//	what if n overflow?
	s := int32(len(bs))
	if s%unit != 0 {
		panic("dat not aligned to unit")
	}

	d := &heapDat{
		pkg:          p.pkg,
		dat:          bs,
		size:         s,
		n:            s / unit,
		regSizeAlign: regSizeAlign,
	}
	p.dat = append(p.dat, d)
	return d
}

func (p *datPool) addVtable(funcs []FuncSym) *heapDat {
	n := int32(len(funcs))
	d := &heapDat{
		dat:          funcs,
		pkg:          p.pkg,
		size:         n * regSize,
		n:            n,
		regSizeAlign: true,
	}
	p.dat = append(p.dat, d)
	return d
}

func (p *datPool) declare(lib *link.Pkg) {
	if lib.Path() != p.pkg {
		panic("package name mismatch")
	}

	if len(p.dat) == 0 {
		return
	}

	ndigit := countDigit(len(p.dat))
	nfmt := fmt.Sprintf(":dat_%%0%dd", ndigit)
	var v *link.Var
	for i, d := range p.dat {
		switch d.dat.(type) {
		case []byte:
			d.name = fmt.Sprintf(nfmt, i)
			align := uint32(0)
			if d.regSizeAlign {
				align = regSize
			}
			v = link.NewVar(align)
			v.Write(d.dat.([]byte))
		case []FuncSym:
			d.name = "_vtable"
			v = link.NewVar(regSize)
			fs := d.dat.([]FuncSym)
			for _, f := range fs {
				v.WriteLink(p.pkg, f.name)
			}
		}

		lib.DeclareVar(d.name)
		lib.DefineVar(d.name, v)
	}
}

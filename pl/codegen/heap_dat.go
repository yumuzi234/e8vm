package codegen

import (
	"fmt"

	"shanhu.io/smlvm/link"
)

type heapDat struct {
	pkg, name    string
	bs           []byte
	regSizeAlign bool
	unitSize     int32
	n            int32
}

func (s *heapDat) String() string {
	return fmt.Sprintf("<dat %dB>", len(s.bs))
}

func (s *heapDat) RegSizeAlign() bool { return s.regSizeAlign }

func (s *heapDat) Size() int32 { return int32(len(s.bs)) }

type datPool struct {
	pkg string
	dat []*heapDat
}

func newDatPool(pkg string) *datPool {
	return &datPool{
		pkg: pkg,
	}
}

func (p *datPool) addDat(bs []byte, unit int32, regSizeAlign bool) *heapDat {
	if int32(len(bs))%unit != 0 {
		panic("dat not aligned to unit")
	}

	d := &heapDat{
		pkg:          p.pkg,
		bs:           bs,
		unitSize:     unit,
		n:            int32(len(bs)) / unit,
		regSizeAlign: regSizeAlign,
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
	for i, d := range p.dat {
		d.name = fmt.Sprintf(nfmt, i)
		align := uint32(0)
		if d.regSizeAlign {
			align = regSize
		}
		v := link.NewVar(align)
		v.Write(d.bs)

		lib.DeclareVar(d.name)
		lib.DefineVar(d.name, v)
	}
}

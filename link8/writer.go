package link8

import (
	"io"

	"encoding/binary"
)

type writer struct {
	w io.Writer
	e error
}

func newWriter(w io.Writer) *writer {
	ret := new(writer)
	ret.w = w
	return ret
}

func (w *writer) Err() error {
	return w.e
}

func (w *writer) Write(buf []byte) (int, error) {
	if w.e != nil {
		return 0, w.e
	}

	n, e := w.w.Write(buf)
	if e != nil {
		w.e = e
	}
	return n, e
}

func (w *writer) writeU32(u uint32) {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], u)
	w.Write(b[:])
}

func (w *writer) writeBareFunc(f *Func) {
	if len(f.links) != 0 {
		panic("not a bare function")
	}

	for _, i := range f.insts {
		w.writeU32(i)
	}
}

func writeVar(w *writer, p *Pkg, v *Var) {
	if v.prePad > 0 {
		w.Write(make([]byte, v.prePad))
	}

	if v.IsZeros() {
		w.Write(make([]byte, v.Size()))
		return
	}

	bs := v.buf.Bytes()

	// fill the links
	for _, lnk := range v.links {
		s := bs[lnk.offset : lnk.offset+4]
		if binary.LittleEndian.Uint32(s) != 0 {
			panic("data to fill non zero")
		}

		v := symAddr(p, lnk)
		binary.LittleEndian.PutUint32(s, v)
	}

	w.Write(bs)
}

func symAddr(p *Pkg, lnk *link) uint32 {
	pkg := p.requires[lnk.pkg]
	s := pkg.symbols[lnk.sym]
	switch s.Type {
	case SymFunc:
		return pkg.Func(lnk.sym).addr
	case SymVar:
		return pkg.Var(lnk.sym).addr
	}
	panic("bug")
}

func funcAddr(p *Pkg, lnk *link) uint32 {
	return p.requires[lnk.pkg].Func(lnk.sym).addr
}

func writeFunc(w *writer, p *Pkg, f *Func) {
	cur := 0
	var curLink *link
	var curIndex int
	updateCur := func() {
		if cur < len(f.links) {
			curLink = f.links[cur]
			curIndex = int(curLink.offset >> 2)
		}
	}

	updateCur()
	for i, inst := range f.insts {
		if curLink != nil && i == curIndex {
			fill := curLink.offset & 0x3
			if fill == FillLink {
				if (inst >> 31) != 0x1 {
					panic("not a jump")
				}
				if (inst & 0x3fffffff) != 0 {
					panic("already filled")
				}

				pc := f.addr + uint32(i)*4 + 4
				target := funcAddr(p, curLink)
				inst |= (target - pc) >> 2
			} else if fill == FillHigh || fill == FillLow {
				if (inst & 0xffff) != 0 {
					panic("already filled")
				}

				v := symAddr(p, curLink)
				if fill == FillHigh {
					inst |= v >> 16
				} else { // fillLow
					inst |= v & 0xffff
				}
			} else {
				panic("invalid fill")
			}

			cur++
			updateCur()
		}

		w.writeU32(inst)
	}
}

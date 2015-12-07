package link8

import (
	"io"

	"encoding/binary"
)

type writer struct {
	lnk *linker
	w   io.Writer
	e   error
}

func newWriter(lnk *linker, w io.Writer) *writer {
	return &writer{lnk: lnk, w: w}
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

func (w *writer) writeVar(v *Var) {
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

		v := w.symAddr(lnk)
		binary.LittleEndian.PutUint32(s, v)
	}

	w.Write(bs)
}

func (w *writer) symAddr(lnk *link) uint32 {
	pkg := w.lnk.pkg(lnk.pkg)
	s := pkg.symbols[lnk.sym]
	switch s.Type {
	case SymFunc:
		return pkg.Func(lnk.sym).addr
	case SymVar:
		return pkg.Var(lnk.sym).addr
	}
	panic("bug")
}

func (w *writer) funcAddr(lnk *link) uint32 {
	return w.lnk.pkg(lnk.pkg).Func(lnk.sym).addr
}

func (w *writer) writeFunc(f *Func) {
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
				target := w.funcAddr(curLink)
				inst |= (target - pc) >> 2
			} else if fill == FillHigh || fill == FillLow {
				if (inst & 0xffff) != 0 {
					panic("already filled")
				}

				v := w.symAddr(curLink)
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

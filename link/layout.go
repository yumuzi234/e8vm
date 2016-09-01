package link

import (
	"fmt"
)

func layout(pkgs map[string]*Pkg, used []*PkgSym, main *Func, initPC uint32) (
	funcs, vars, zeros []*PkgSym, err error,
) {
	pt := initPC
	const codeMax uint32 = 0xffffffff

	putFunc := func(f *Func) error {
		f.addr = pt
		size := f.Size()
		if size > codeMax-pt {
			return fmt.Errorf("code section too large")
		}
		pt += size
		return nil
	}

	if err = putFunc(main); err != nil {
		return
	}

	for _, ps := range used {
		pkg := pkgs[ps.Pkg]
		s := pkg.SymbolByName(ps.Sym)
		switch s.Type {
		case SymFunc:
			funcs = append(funcs, ps)
			f := pkg.Func(ps.Sym)
			if err = putFunc(f); err != nil {
				return
			}
		case SymVar:
			v := pkg.Var(ps.Sym)
			if !v.IsZeros() {
				vars = append(vars, ps)
			} else {
				zeros = append(zeros, ps)
			}
		default:
			panic("bug")
		}
	}

	const dataMax uint32 = 0xffffffff

	putVar := func(ps *PkgSym) error {
		v := pkgVar(pkgs, ps)
		if v.align > 1 && pt%v.align != 0 {
			v.prePad = v.align - pt%v.align
			pt += v.prePad
		} else {
			v.prePad = 0
		}
		if v.align > 1 && pt%v.align != 0 {
			panic("bug")
		}

		v.addr = pt
		size := v.Size()
		if size > dataMax-pt {
			return fmt.Errorf("binary too large")
		}

		pt += size
		return nil
	}

	for _, ps := range vars {
		if err = putVar(ps); err != nil {
			return
		}
	}

	for _, ps := range zeros {
		if err = putVar(ps); err != nil {
			return
		}
	}

	const totalMax = 1024 * 1024 // 1MB
	if pt-initPC > totalMax {
		return nil, nil, nil, fmt.Errorf("binary too large")
	}

	return
}

package link8

import (
	"fmt"
)

func layout(used []pkgSym, initPC uint32) (
	funcs, vars, zeros []pkgSym, e error,
) {
	pt := initPC
	const codeMax uint32 = 0xffffffff

	for _, ps := range used {
		typ := ps.Type()
		switch typ {
		case SymFunc:
			funcs = append(funcs, ps)

			f := ps.Func()
			f.addr = pt
			size := f.Size()
			if size > codeMax-pt {
				return nil, nil, nil, fmt.Errorf("code section too large")
			}
			pt += size
		case SymVar:
			v := ps.Var()
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

	putVar := func(v *Var) error {
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
		err := putVar(ps.Var())
		if err != nil {
			return nil, nil, nil, err
		}
	}

	for _, ps := range zeros {
		err := putVar(ps.Var())
		if err != nil {
			return nil, nil, nil, err
		}
	}

	const totalMax = 1024 * 1024 // 1MB
	if pt-initPC > totalMax {
		return nil, nil, nil, fmt.Errorf("binary too large")
	}

	return
}

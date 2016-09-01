package codegen

import (
	"e8vm.io/e8vm/lexing"
	"e8vm.io/e8vm/link"
)

// BuildPkg builds a package and returns the built lib
func BuildPkg(p *Pkg) (*link.Pkg, []*lexing.Error) {
	p.strPool.declare(p.lib)
	p.datPool.declare(p.lib)

	for _, v := range p.vars {
		var align uint32 = regSize
		if v.size <= 1 {
			align = 1
		}
		obj := link.NewVar(align)
		obj.Zeros(uint32(v.size))
		p.lib.DefineVar(v.name, obj)
	}

	for _, f := range p.funcs {
		genFunc(p.g, f)
	}
	if errs := p.g.Errs(); errs != nil {
		return nil, errs
	}

	for _, f := range p.funcs {
		writeFunc(p, f)
	}

	if p.tests != nil {
		v := link.NewVar(regSize)
		for _, f := range p.tests.funcs {
			if err := v.WriteLink(p.path, f.name); err != nil {
				panic(err)
			}
		}
		p.lib.DefineVar(p.tests.name, v)
	}

	return p.lib, nil
}

package asm8

func buildFile(b *builder, f *file) {
	pkg := b.curPkg

	for _, fn := range f.funcs {
		if obj := buildFunc(b, fn); obj != nil {
			pkg.DefineFunc(fn.Name.Lit, obj)
		}
	}

	for _, v := range f.vars {
		if obj := buildVar(b, v); obj != nil {
			pkg.DefineVar(v.Name.Lit, obj)
		}
	}
}

package g8

import (
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/sym8"
)

type pkg struct {
	files map[string]*ast.File

	consts  []*ast.ConstDecls
	funcs   []*ast.Func
	methods []*ast.Func
	structs []*ast.Struct
	vars    []*ast.VarDecls

	structMap   map[string]*structInfo
	structOrder []*structInfo
	funcObjs    []*objFunc

	tops *sym8.Table

	testNames []string
	testList  ir.Ref
}

func newPkg(asts map[string]*ast.File) *pkg {
	ret := new(pkg)
	ret.files = asts
	ret.structMap = make(map[string]*structInfo)

	return ret
}

func (p *pkg) declareConsts(b *builder) {
	for _, decls := range p.consts {
		for _, d := range decls.Decls {
			buildConstDecl(b, d)
		}
	}
}

func (p *pkg) declareStructs(b *builder) {
	for _, d := range p.structs {
		info := declareStruct(b, d)
		if info == nil {
			continue
		}

		name := info.Name()
		if p.structMap[name] != nil {
			panic("struct with same name")
		}
		p.structMap[name] = info // also save it in map for dep analysis
	}
}

func (p *pkg) pushStruct(info *structInfo, b *builder) {
	info.queuing = true
	name := info.Name()

	for _, dep := range info.deps {
		if dep == name {
			b.Errorf(info.name.Pos, "struct %s depends on itself", name)
			continue
		}

		depInfo := p.structMap[dep]
		if depInfo == nil || depInfo.queued {
			continue
		}

		if depInfo.queuing {
			b.Errorf(info.name.Pos,
				"struct %s circular depends on struct %s",
				info.name.Lit, dep,
			)
			continue
		}

		p.pushStruct(depInfo, b)
	}

	info.queuing = false
	info.queued = true
	p.structOrder = append(p.structOrder, info)
}

func (p *pkg) sortStructs(b *builder) {
	p.structOrder = nil

	for _, info := range p.structMap {
		if !info.queued {
			p.pushStruct(info, b)
		}
	}
}

func (p *pkg) defineStructs(b *builder) {
	for _, info := range p.structOrder {
		defineStructFields(b, info)
		info.defined = true
	}
}

func (p *pkg) declareFuncs(b *builder) {
	for _, f := range p.funcs {
		ret := declareFunc(b, f)
		if ret != nil {
			p.funcObjs = append(p.funcObjs, ret)
		}
	}

	for _, f := range p.methods {
		sname := f.Recv.StructName
		info := p.structMap[sname.Lit]
		if info == nil {
			b.Errorf(sname.Pos, "struct %s not defined", sname.Lit)
			continue
		}

		obj := declareMethod(b, info, f)
		if obj != nil {
			info.methodObjs = append(info.methodObjs, obj)
		}
	}

	for _, s := range p.structOrder {
		declareStructMethods(b, s)
	}
}

func (p *pkg) declareVars(b *builder) {
	for _, decls := range p.vars {
		for _, d := range decls.Decls {
			buildGlobalVarDecl(b, d)
		}
	}
}

func (p *pkg) buildFuncs(b *builder) {
	b.this = nil
	for _, f := range p.funcObjs {
		buildFunc(b, f)
	}
	for _, s := range p.structOrder {
		buildMethods(b, s)
	}
}

func (p *pkg) collectSymbols(b *builder) {
	for _, f := range p.files {
		decls := f.Decls
		for _, d := range decls {
			switch d := d.(type) {
			case *ast.Func:
				if d.Recv == nil {
					p.funcs = append(p.funcs, d)
				} else {
					p.methods = append(p.methods, d)
				}
			case *ast.VarDecls:
				p.vars = append(p.vars, d)
			case *ast.Struct:
				p.structs = append(p.structs, d)
			case *ast.ConstDecls:
				p.consts = append(p.consts, d)
			default:
				b.Errorf(nil, "invalid top declare: %T", d)
			}
		}
	}
}

func (p *pkg) onlyFile() *ast.File {
	if len(p.files) != 1 {
		return nil
	}
	for _, f := range p.files {
		return f
	}
	panic("unrechable")
}

func (p *pkg) declareImports(b *builder, pinfo *build8.PkgInfo) {
	if f := p.onlyFile(); f != nil {
		declareImports(b, f, pinfo)
		return
	}

	for name, f := range p.files {
		if name == "import.g" {
			if len(f.Decls) > 0 {
				first := f.Decls[0]
				b.Errorf(ast.DeclPos(first),
					`"import.g" in multi-file package only allows import`,
				)
			} else {
				declareImports(b, f, pinfo)
			}

			continue
		}

		if f.Imports != nil {
			b.Errorf(f.Imports.Kw.Pos,
				`import only allowed in "import.g" for multi-file package`,
			)
		}
	}
}

func (p *pkg) buildTests(b *builder) {
	tests := listTests(p.tops)
	n := len(tests)

	if n > 100000 {
		b.Errorf(nil, "too many tests in the package")
		return
	}

	perm := b.rand.Perm(n)

	var irs []*ir.Func
	var names []string
	for _, index := range perm {
		t := tests[index]
		irs = append(irs, t.ref.IR().(*ir.Func))
		names = append(names, t.name)
	}
	if n > 0 {
		p.testList = b.p.NewTestList(":tests", irs)
		p.testNames = names
	}
}

func (p *pkg) build(b *builder, pinfo *build8.PkgInfo) {
	p.tops = sym8.NewTable()
	b.scope.PushTable(p.tops) // package scope
	defer b.scope.Pop()

	o := func(f func(b *builder)) {
		if b.Errs() != nil {
			return
		}
		f(b)
	}

	p.declareImports(b, pinfo)

	o(p.collectSymbols)
	o(p.declareConsts)
	o(p.declareStructs)
	o(p.sortStructs)
	o(p.defineStructs)
	o(p.declareFuncs)
	o(p.declareVars)
	o(p.buildFuncs)
	o(p.buildTests)

	if b.Errs() != nil {
		return
	}

	addStart(b)
	if p.testList != nil {
		addTestStart(b, p.testList, len(p.testNames))
	}
}

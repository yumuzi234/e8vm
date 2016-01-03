package sempass

import (
	"fmt"

	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

// NewBuilder creates a new builder with a specific path.
func NewBuilder(path string, scope *sym8.Scope) *Builder {
	ret := newBuilder(path)
	ret.exprFunc = buildExpr
	ret.constFunc = buildConstExpr
	ret.typeFunc = buildType
	ret.stmtFunc = buildStmt

	ret.scope = scope // TODO: remove this

	return ret
}

func makeBuilder(path string) *Builder {
	scope := sym8.NewScope()
	return NewBuilder(path, scope)
}

// Pkg is a package that contains information for a sementics pass.
type Pkg struct {
	Path    string
	Files   map[string]*ast.File
	Imports map[string]*build8.Package
}

type symbols struct {
	consts  []*ast.ConstDecls
	funcs   []*ast.Func
	methods []*ast.Func
	structs []*ast.Struct
	vars    []*ast.VarDecls
}

func (p *Pkg) symbols() *symbols {
	ret := new(symbols)
	for _, f := range p.Files {
		decls := f.Decls
		for _, d := range decls {
			switch d := d.(type) {
			case *ast.Func:
				if d.Recv == nil {
					ret.funcs = append(ret.funcs, d)
				} else {
					ret.methods = append(ret.methods, d)
				}
			case *ast.VarDecls:
				ret.vars = append(ret.vars, d)
			case *ast.Struct:
				ret.structs = append(ret.structs, d)
			case *ast.ConstDecls:
				ret.consts = append(ret.consts, d)
			default:
				panic(fmt.Errorf("invalid top declare: %T", d))
			}
		}
	}

	return ret
}

func structSyms(pkgStructs []*pkgStruct) []*sym8.Symbol {
	ret := make([]*sym8.Symbol, 0, len(pkgStructs))
	for _, ps := range pkgStructs {
		ret = append(ret, ps.sym)
	}
	return ret
}

func declareFuncs(b *Builder, funcs []*ast.Func) (
	[]*pkgFunc, []*tast.FuncAlias,
) {
	var ret []*pkgFunc
	var aliases []*tast.FuncAlias

	for _, f := range funcs {
		if f.Alias != nil {
			a := buildFuncAlias(b, f)
			if a != nil {
				aliases = append(aliases, a)
			}
			continue
		}

		r := declareFunc(b, f)
		if r != nil {
			ret = append(ret, r)
		}
	}

	return ret, aliases
}

func buildFuncs(b *Builder, funcs []*pkgFunc) []*tast.Func {
	b.this = nil
	b.thisType = nil

	ret := make([]*tast.Func, 0, len(funcs))
	for _, f := range funcs {
		res := buildFunc(b, f)
		if res != nil {
			ret = append(ret, res)
		}
	}

	return ret
}

func declareMethods(
	b *Builder, methods []*ast.Func, pkgStructs []*pkgStruct,
) []*pkgFunc {
	m := make(map[string]*pkgStruct)
	for _, ps := range pkgStructs {
		m[ps.name.Lit] = ps
	}

	var ret []*pkgFunc

	// inlined ones
	for _, ps := range pkgStructs {
		for _, f := range ps.ast.Methods {
			pf := declareMethod(b, ps, f)
			if pf != nil {
				ret = append(ret, pf)
			}
		}
	}

	// go-like ones
	for _, f := range methods {
		recv := f.Recv.StructName
		ps := m[recv.Lit]
		if ps != nil {
			b.Errorf(recv.Pos, "struct %s not defined", recv.Lit)
			continue
		}

		pf := declareMethod(b, ps, f)
		if pf != nil {
			ret = append(ret, pf)
		}
	}

	return ret
}

func buildMethods(b *Builder, funcs []*pkgFunc) []*tast.Func {
	var ret []*tast.Func
	for _, f := range funcs {
		r := buildMethod(b, f)
		if r != nil {
			ret = append(ret, r)
		}
	}

	return ret
}

// Build builds a package from an set of file AST's to a typed-AST.
func (p *Pkg) Build() (*tast.Pkg, []*lex8.Error) {
	syms := p.symbols()
	b := makeBuilder(p.Path)

	tops := sym8.NewTable()
	b.scope.PushTable(tops)
	defer b.scope.Pop()

	// TODO: imports

	consts := buildPkgConsts(b, syms.consts)
	if errs := b.Errs(); errs != nil {
		return nil, errs
	}

	pkgStructs := buildStructs(b, syms.structs)
	if errs := b.Errs(); errs != nil {
		return nil, errs
	}

	pkgFuncs, aliases := declareFuncs(b, syms.funcs)
	if errs := b.Errs(); errs != nil {
		return nil, errs
	}

	pkgMethods := declareMethods(b, syms.methods, pkgStructs)
	if errs := b.Errs(); errs != nil {
		return nil, errs
	}

	vars := buildPkgVars(b, syms.vars)
	if errs := b.Errs(); errs != nil {
		return nil, errs
	}

	funcs := buildFuncs(b, pkgFuncs)
	if errs := b.Errs(); errs != nil {
		return nil, errs
	}

	methods := buildMethods(b, pkgMethods)
	if errs := b.Errs(); errs != nil {
		return nil, errs
	}

	structs := structSyms(pkgStructs)

	return &tast.Pkg{
		Consts:      consts,
		Structs:     structs,
		Vars:        vars,
		Funcs:       funcs,
		Methods:     methods,
		FuncAliases: aliases,
	}, nil
}

// BuildPkgConsts is a temp function for building package consts.
var BuildPkgConsts = buildPkgConsts

// BuildPkgVars is a temp function for building package vars.
var BuildPkgVars = buildPkgVars

package sempass

import (
	"fmt"

	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/dagvis"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

func makeBuilder(path string, scope *sym8.Scope) *builder {
	ret := newBuilder(path, scope)
	ret.exprFunc = buildExpr
	ret.constFunc = buildConstExpr
	ret.typeFunc = buildType
	ret.stmtFunc = buildStmt

	return ret
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

func (p *Pkg) onlyFile() *ast.File {
	if len(p.Files) == 1 {
		for _, f := range p.Files {
			return f
		}
	}
	return nil
}

func (p *Pkg) buildImports(
	b *builder, imps map[string]*build8.Package,
) []*sym8.Symbol {
	if f := p.onlyFile(); f != nil {
		return buildImports(b, f, imps)
	}

	var ret []*sym8.Symbol
	for name, f := range p.Files {
		if name == "import.g" {
			if len(f.Decls) > 0 {
				first := f.Decls[0]
				b.Errorf(ast.DeclPos(first),
					`"import.g" in multi-file packages only allows import`,
				)
			} else {
				ret = buildImports(b, f, imps)
			}
			continue
		}

		if f.Imports != nil {
			b.Errorf(f.Imports.Kw.Pos,
				`import only allowed in "import.g" for multi-file package`,
			)
		}
	}

	return ret
}

// Build builds a package from an set of file AST's to a typed-AST.
func (p *Pkg) Build(scope *sym8.Scope) (
	*tast.Pkg, *dagvis.Graph, []*lex8.Error,
) {
	b := makeBuilder(p.Path, scope)
	b.initDeps(p.Files)

	imports := p.buildImports(b, p.Imports)
	if errs := b.Errs(); errs != nil {
		return nil, nil, errs
	}

	syms := p.symbols()

	consts := buildPkgConsts(b, syms.consts)
	if errs := b.Errs(); errs != nil {
		return nil, nil, errs
	}

	pkgStructs := buildStructs(b, syms.structs)
	if errs := b.Errs(); errs != nil {
		return nil, nil, errs
	}

	pkgFuncs, aliases := declareFuncs(b, syms.funcs)
	if errs := b.Errs(); errs != nil {
		return nil, nil, errs
	}

	pkgMethods := declareMethods(b, syms.methods, pkgStructs)
	if errs := b.Errs(); errs != nil {
		return nil, nil, errs
	}

	vars := buildPkgVars(b, syms.vars)
	if errs := b.Errs(); errs != nil {
		return nil, nil, errs
	}

	funcs := buildFuncs(b, pkgFuncs)
	if errs := b.Errs(); errs != nil {
		return nil, nil, errs
	}

	methods := buildMethods(b, pkgMethods)
	if errs := b.Errs(); errs != nil {
		return nil, nil, errs
	}

	depGraph := b.depGraph()
	structs := structSyms(pkgStructs)

	return &tast.Pkg{
		Imports:     imports,
		Consts:      consts,
		Structs:     structs,
		Vars:        vars,
		Funcs:       funcs,
		Methods:     methods,
		FuncAliases: aliases,
	}, depGraph, nil
}

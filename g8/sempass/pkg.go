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

// Build builds a package from an set of file AST's to a typed-AST.
func (p *Pkg) Build() (*tast.Pkg, []*lex8.Error) {
	syms := p.symbols()
	b := makeBuilder(p.Path)
	consts := buildPkgConsts(b, syms.consts)
	structs := buildStructs(b, syms.structs)
	funcs, aliases := declareFuncs(b, syms.funcs)

	errs := b.Errs()
	if errs != nil {
		return nil, errs
	}

	return &tast.Pkg{
		Consts:      consts,
		Structs:     structs,
		Funcs:       funcs,
		FuncAliases: aliases,
	}, nil
}

// BuildPkgConsts is a temp function for building package consts.
var BuildPkgConsts = buildPkgConsts

// BuildStructs is a temp function for building struct types.
var BuildStructs = buildStructs

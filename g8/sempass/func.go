package sempass

import (
	"e8vm.io/e8vm/asm8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/sym8"
)

type pkgFunc struct {
	sym *sym8.Symbol
	f   *ast.Func

	recv *pkgStruct
}

func declareFuncSym(b *builder, f *ast.Func, t types.T) *sym8.Symbol {
	name := f.Name.Lit
	s := sym8.Make(b.path, name, tast.SymFunc, nil, t, f.Name.Pos)
	conflict := b.scope.Declare(s)
	if conflict != nil {
		b.Errorf(f.Name.Pos, "%q already defined as a %s",
			name, tast.SymStr(conflict.Type),
		)
		b.Errorf(conflict.Pos, "previously defined here")
		return nil
	}
	return s
}

func buildFuncAlias(b *builder, f *ast.Func) *tast.FuncAlias {
	t := buildFuncType(b, nil, f.FuncSig)
	if t == nil {
		return nil
	}

	alias := f.Alias
	pkg := buildPkgRef(b, alias.Pkg)
	if pkg == nil {
		return nil
	}
	if pkg.Lang != "asm8" {
		b.Errorf(alias.Pkg.Pos, "can only alias functions in asm packages")
		return nil
	}

	sym := findPackageSym(b, alias.Name, pkg)
	if sym == nil {
		return nil
	}

	if sym.Type != asm8.SymFunc {
		b.Errorf(alias.Name.Pos, "the symbol is not a function")
		return nil
	}

	ret := declareFuncSym(b, f, t)
	if ret == nil {
		return nil
	}

	return &tast.FuncAlias{Sym: ret, Of: sym}
}

func declareFunc(b *builder, f *ast.Func) *pkgFunc {
	if f.Alias != nil {
		panic("bug")
	}

	t := buildFuncType(b, nil, f.FuncSig)
	if t == nil {
		return nil
	}

	s := declareFuncSym(b, f, t)
	if s == nil {
		return nil
	}

	return &pkgFunc{sym: s, f: f}
}

func declareParas(
	b *builder, lst *ast.ParaList, ts []*types.Arg,
) []*sym8.Symbol {
	var ret []*sym8.Symbol
	paras := lst.Paras

	for i, t := range ts {
		if t.Name == thisName {
			panic("trying to declare <this>")
		}

		var s *sym8.Symbol
		if t.Name != "" {
			s = declareVar(b, paras[i].Ident, t.T)
		}
		ret = append(ret, s)
	}
	return ret
}

func buildFunc(b *builder, f *pkgFunc) *tast.Func {
	b.scope.Push()
	defer b.scope.Pop()

	t := f.sym.ObjType.(*types.Func)
	b.retNamed = f.f.NamedRet()
	b.retType = t.RetTypes

	ret := new(tast.Func)
	ret.Sym = f.sym

	if b.this != nil {
		if f.f.Recv != nil {
			if recvTok := f.f.Recv.Recv; recvTok != nil {
				recvSym := declareVar(b, recvTok, b.thisType)
				if recvSym == nil {
					return nil
				}
				ret.Receiver = recvSym
			}
		} else {
			// marking keyword <this> if it is an inlined method
			ret.This = b.thisType
		}
	}

	if b.this != nil {
		ret.Args = declareParas(b, f.f.Args, t.Args[1:])
	} else {
		ret.Args = declareParas(b, f.f.Args, t.Args)
	}

	if b.retNamed {
		ret.NamedRets = declareParas(b, f.f.Rets, t.Rets)
	}

	ret.Body = buildStmts(b, f.f.Body.Stmts)

	if len(b.retType) > 0 && !isBlockTerminal(f.f.Body) {
		b.Errorf(f.f.Body.Rbrace.Pos, "missing return at the end of function")
	}

	// clear for safety
	b.retType = nil
	b.retNamed = false

	return ret
}

func buildMethod(b *builder, f *pkgFunc) *tast.Func {
	this := f.recv.pt
	b.thisType = this
	if f.f.Recv != nil { // go-like, explicit receiver
		b.this = tast.NewAddressableRef(this)
	} else { // inlined
		b.this = tast.NewRef(this)
		b.scope.PushTable(f.recv.t.Syms)
		defer b.scope.Pop()
	}

	return buildFunc(b, f)
}

func declareMethod(b *builder, ps *pkgStruct, f *ast.Func) *pkgFunc {
	if f.Alias != nil {
		b.Errorf(f.Alias.Eq.Pos, "cannot alias a function for a method")
		return nil
	}

	t := buildFuncType(b, ps.pt, f.FuncSig)
	if t == nil {
		return nil
	}

	name := f.Name.Lit
	sym := sym8.Make(b.path, name, tast.SymFunc, nil, t, f.Name.Pos)
	conflict := ps.t.Syms.Declare(sym)
	if conflict != nil {
		b.Errorf(f.Name.Pos, "member %s already defined", name)
		b.Errorf(conflict.Pos, "previously defined here")
		return nil
	}

	return &pkgFunc{sym: sym, f: f, recv: ps}
}

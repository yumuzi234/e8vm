package g8

import (
	"e8vm.io/e8vm/asm8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

func declareFuncSym(b *builder, f *ast.Func, t types.T) *objFunc {
	// NewFunc() will create the variables required for the sigs
	name := f.Name.Lit
	ret := new(objFunc)
	ret.name = name
	ret.f = f

	// add this item to the top scope
	s := sym8.Make(b.path, name, tast.SymFunc, ret, t, f.Name.Pos)
	conflict := b.scope.Declare(s) // lets declare the function
	if conflict != nil {
		b.Errorf(f.Name.Pos, "%q already declared as a %s",
			name, tast.SymStr(conflict.Type),
		)
		b.Errorf(conflict.Pos, "previously declared here")
		return nil
	}
	return ret
}

func buildPkgRef(b *builder, ident *lex8.Token) *types.Pkg {
	s := b.scope.Query(ident.Lit)
	if s == nil {
		b.Errorf(ident.Pos, "undefined identifier %s", ident.Lit)
		return nil
	}

	b.spass.RefSym(s, ident.Pos)
	if s.Type != tast.SymImport {
		b.Errorf(ident.Pos, "%s is not an imported package", ident.Lit)
		return nil
	}

	return s.Obj.(*objImport).ref.Type().(*types.Pkg)
}

func findPackageSym(
	b *builder, sub *lex8.Token, pkg *types.Pkg,
) *sym8.Symbol {
	sym := pkg.Syms.Query(sub.Lit)
	if sym == nil {
		b.Errorf(sub.Pos, "%s has no symbol named %s",
			pkg, sub.Lit,
		)
		return nil
	}
	name := sym.Name()
	if !sym8.IsPublic(name) && sym.Pkg() != b.path {
		b.Errorf(sub.Pos, "symbol %s is not public", name)
		return nil
	}

	return sym
}

func declareFuncAlias(b *builder, f *ast.Func, t *types.Func) *objFunc {
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

	obj := declareFuncSym(b, f, t)
	if obj == nil {
		return nil
	}

	sig := makeFuncSig(t)
	fsym := ir.NewFuncSym(sym.Pkg(), alias.Name.Lit, sig)
	obj.ref = newRef(t, fsym)
	obj.isAlias = true

	return obj
}

func declareFunc(b *builder, f *ast.Func) *objFunc {
	t := buildFuncType(b, nil, f.FuncSig)
	if t == nil {
		return nil
	}

	if f.Alias != nil {
		return declareFuncAlias(b, f, t)
	}

	ret := declareFuncSym(b, f, t)
	if ret == nil {
		return nil
	}

	sig := makeFuncSig(t)
	irFunc := b.p.NewFunc(b.anonyName(f.Name.Lit), f.Name.Pos, sig)
	ret.ref = newRef(t, irFunc)

	return ret
}

func declareParas(
	b *builder, lst *ast.ParaList, ts []*types.Arg, irs []ir.Ref,
) {
	if len(ts) != len(irs) {
		panic("bug")
	}
	paras := lst.Paras

	for i, t := range ts {
		if t.Name == thisName {
			panic("trying to declare this")
		}
		if t.Name == "" { // unamed parameter
			continue
		}

		r := newAddressableRef(t.T, irs[i])
		declareVarRef(b, paras[i].Ident, r)
	}
}

func makeRetRef(ts []*types.Arg, irs []ir.Ref) *ref {
	if len(ts) != len(irs) {
		panic("bug")
	}
	if len(ts) == 0 {
		return nil
	}

	ret := new(ref)
	for i, t := range ts {
		ret = appendRef(ret, newAddressableRef(t.T, irs[i]))
	}
	return ret
}

func buildFunc(b *builder, f *objFunc) {
	b.scope.Push() // func body scope
	defer b.scope.Pop()

	t := f.ref.Type().(*types.Func) // the signature of the function
	irFunc := f.ref.IR().(*ir.Func)
	b.f = irFunc
	b.fretNamed = f.f.NamedRet()

	if f.f.Recv != nil {
		if recvTok := f.f.Recv.Recv; recvTok != nil {
			declareVarRef(b, recvTok, b.this)
		}
	}

	// build context for resolving symbols
	if b.this != nil {
		declareParas(b, f.f.Args, t.Args[1:], irFunc.ArgRefs()[1:])
	} else {
		declareParas(b, f.f.Args, t.Args, irFunc.ArgRefs())
	}

	retIRRefs := irFunc.RetRefs()
	if b.fretNamed {
		declareParas(b, f.f.Rets, t.Rets, retIRRefs)
	}

	b.fretRef = makeRetRef(t.Rets, retIRRefs)
	if b.fretRef == nil {
		b.spass.SetRetType(nil, b.fretNamed)
	} else {
		b.spass.SetRetType(t.RetTypes, b.fretNamed)
	}

	b.b = b.f.NewBlock(nil)
	b.buildStmts(f.f.Body.Stmts)

	b.fretRef = nil
	b.spass.SetRetType(nil, false)
}

func buildMethodFunc(b *builder, this *types.Pointer, f *objFunc) {
	t := f.ref.Type().(*types.Func)
	ir := f.ref.IR().(*ir.Func)

	if len(t.Args) == 0 {
		panic("this pointer missing")
	}

	if !b.golike {
		b.this = newRef(this, ir.ThisRef())
		b.spass.SetThis(tast.NewRef(this))
	} else {
		b.this = newAddressableRef(this, ir.ThisRef())
		b.spass.SetThis(tast.NewAddressableRef(this))
	}
	buildFunc(b, f)
}

func genFunc(b *builder, f *tast.Func, ref *ref) {
	irFunc := ref.IR().(*ir.Func)
	b.f = irFunc

	if f.Receiver != nil {
		// bind the receiver
		t := f.Receiver.ObjType.(types.T)
		f.Receiver.Obj = newAddressableRef(t, irFunc.ThisRef())
	} else if f.This != nil {
		// bind this pointer
		b.this = newRef(f.This, irFunc.ThisRef())
	}

	// bind arg symbols
	args := irFunc.ArgRefs()
	if f.IsMethod() {
		args = args[1:] // skip <this>
	}
	for i, s := range f.Args {
		if s != nil {
			s.Obj = args[i]
		}
	}

	if f.NamedRets != nil {
		// bind named return symbols
		rets := irFunc.RetRefs()
		for i, s := range f.NamedRets {
			if s != nil {
				s.Obj = rets[i]
			}
		}
	}

	for _, stmt := range f.Body {
		b.buildStmt2(stmt)
	}
}

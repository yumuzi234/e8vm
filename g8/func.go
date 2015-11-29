package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

func declareFunc(b *builder, f *ast.Func) *objFunc {
	t := buildFuncType(b, nil, f.FuncSig)
	if t == nil {
		return nil
	}

	// NewFunc() will create the variables required for the sigs
	name := f.Name.Lit
	ret := new(objFunc)
	ret.name = name
	ret.f = f

	// add this item to the top scope
	s := sym8.Make(b.symPkg, name, symFunc, ret, f.Name.Pos)
	conflict := b.scope.Declare(s) // lets declare the function
	if conflict != nil {
		b.Errorf(f.Name.Pos, "%q already declared as a %s",
			name, symStr(conflict.Type),
		)
		b.Errorf(conflict.Pos, "previously declared here")
		return nil
	}

	irFunc := b.p.NewFunc(b.anonyName(name), t.Sig)
	ret.ref = newRef(t, irFunc)

	return ret
}

func paraIdent(b *builder,
	paras []*ast.Para, i int, withThis bool,
) *lex8.Token {
	if !withThis {
		return paras[i].Ident
	}
	return paras[i-1].Ident
}

func declareParas(b *builder,
	lst *ast.ParaList, ts []*types.Arg, irs []ir.Ref, withThis bool,
) {
	if len(ts) != len(irs) {
		panic("bug")
	}

	for i, t := range ts {
		if t.Name == "" || t.Name == thisName {
			continue
		}

		r := newAddressableRef(t.T, irs[i])
		declareVarRef(b, paraIdent(b, lst.Paras, i, withThis), r)
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

	if b.golike && f.f.Recv != nil {
		if recvTok := f.f.Recv.Recv; recvTok != nil {
			declareVarRef(b, recvTok, b.this)
		}
	}

	// build context for resolving symbols
	retIRRefs := irFunc.RetRefs()
	declareParas(b, f.f.Args, t.Args, irFunc.ArgRefs(), b.this != nil)
	declareParas(b, f.f.Rets, t.Rets, retIRRefs, false)
	b.fretRef = makeRetRef(t.Rets, retIRRefs)

	b.b = b.f.NewBlock(nil)
	b.buildStmts(f.f.Body.Stmts)
}

func buildMethodFunc(b *builder, s *structInfo, f *objFunc) {
	t := f.ref.Type().(*types.Func)
	ir := f.ref.IR().(*ir.Func)

	if len(t.Args) == 0 {
		panic("this pointer missing")
	}

	if !b.golike {
		b.this = newRef(s.pt, ir.ThisRef())
	} else {
		b.this = newAddressableRef(s.pt, ir.ThisRef())
	}
	buildFunc(b, f)
}

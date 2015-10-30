package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

func allocVars(b *builder, toks []*lex8.Token, ts []types.T) *ref {
	ret := new(ref)
	ret.typ = append(ret.typ, ts...)

	for i, tok := range toks {
		// the name here is just for debugging
		// it is not a var declare
		t := ret.typ[i]
		if types.IsNil(t) {
			b.Errorf(tok.Pos, "cannot infer type from nil for %q", tok.Lit)
			return nil
		}
		if types.IsConst(t) {
			t = types.Int
			ret.typ[i] = t
		}
		if !types.IsAllocable(t) {
			b.Errorf(tok.Pos, "cannot allocate for %s", t)
			return nil
		}

		v := b.newLocal(t, tok.Lit)
		ret.ir = append(ret.ir, v)
		ret.addressable = append(ret.addressable, true)
	}
	return ret
}

func declareVar(b *builder, tok *lex8.Token) *objVar {
	name := tok.Lit
	v := &objVar{name: name}
	s := sym8.Make(b.symPkg, name, symVar, v, tok.Pos)
	conflict := b.scope.Declare(s)
	if conflict != nil {
		b.Errorf(tok.Pos, "%q already declared as a %s",
			name, symStr(conflict.Type),
		)
		return nil
	}
	return v
}

func declareVarRef(b *builder, tok *lex8.Token, r *ref) {
	obj := declareVar(b, tok)
	if obj != nil {
		obj.ref = r
	}
}

func declareVars(b *builder, toks []*lex8.Token, r *ref) {
	for i, t := range r.typ {
		declareVarRef(b, toks[i], newAddressableRef(t, r.ir[i]))
	}
}

func define(b *builder, idents []*lex8.Token, expr *ref, eq *lex8.Token) {
	// check count matching
	nleft := len(idents)
	nright := expr.Len()
	if nleft != nright {
		b.Errorf(eq.Pos,
			"defined %d identifers with %d expressions",
			nleft, nright,
		)
		return
	}

	left := allocVars(b, idents, expr.typ)
	if left == nil {
		return
	}

	if assign(b, left, expr, eq) {
		declareVars(b, idents, left)
	}
}

func buildDefineStmt(b *builder, stmt *ast.DefineStmt) {
	right := buildExprList(b, stmt.Right)
	if right == nil { // an error occured on the expression list
		return
	}

	idents, err := buildIdentExprList(b, stmt.Left)
	if err != nil {
		b.Errorf(ast.ExprPos(err), "left side of := must be identifer")
		return
	}

	define(b, idents, right, stmt.Define)
}

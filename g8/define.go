package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

func allocVars(b *builder, toks []*lex8.Token, ts []types.T) *ref {
	ret := new(ref)

	for i, tok := range toks {
		t := ts[i]
		if types.IsNil(t) {
			b.Errorf(tok.Pos, "cannot infer type from nil for %q", tok.Lit)
			return nil
		}
		if types.IsConst(t) {
			t = types.Int
		}
		if !types.IsAllocable(t) {
			b.Errorf(tok.Pos, "cannot allocate for %s", t)
			return nil
		}

		v := b.newLocal(t, tok.Lit)
		ret = appendRef(ret, newAddressableRef(t, v))
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
	n := r.Len()
	for i := 0; i < n; i++ {
		ref := r.At(i)
		if !ref.Addressable() {
			panic("ref not addressable")
		}
		declareVarRef(b, toks[i], ref)
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

	left := allocVars(b, idents, expr.TypeList())
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

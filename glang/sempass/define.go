package sempass

import (
	"e8vm.io/e8vm/glang/ast"
	"e8vm.io/e8vm/glang/parse"
	"e8vm.io/e8vm/glang/tast"
	"e8vm.io/e8vm/glang/types"
	"e8vm.io/e8vm/lexing"
	"e8vm.io/e8vm/syms"
)

// allocPrepare checks if the provided types are all allocable, and insert
// implicit type casts if needed. Only literay expression list needs alloc
// prepare.
func allocPrepare(
	b *builder, toks []*lexing.Token, lst *tast.ExprList,
) *tast.ExprList {
	ret := tast.NewExprList()
	for i, tok := range toks {
		e := lst.Exprs[i]
		t := e.Type()
		if types.IsNil(t) {
			b.Errorf(tok.Pos, "cannot infer type from nil for %q", tok.Lit)
			return nil
		}
		if v, ok := types.NumConst(t); ok {
			e = constCastInt(b, tok.Pos, v, e)
			if e == nil {
				return nil
			}
		}
		if !types.IsAllocable(t) {
			b.Errorf(tok.Pos, "cannot allocate for %s", t)
			return nil
		}
		ret.Append(e)
	}
	return ret
}

func define(
	b *builder, ids []*lexing.Token, expr tast.Expr, eq *lexing.Token,
) *tast.Define {
	// check count matching
	r := expr.R()
	nleft := len(ids)
	nright := r.Len()
	if nleft != nright {
		b.Errorf(eq.Pos,
			"defined %d identifers with %d expressions",
			nleft, nright,
		)
		return nil
	}

	if exprList, ok := tast.MakeExprList(expr); ok {
		exprList = allocPrepare(b, ids, exprList)
		if exprList == nil {
			return nil
		}
		expr = exprList
	}

	var ret []*syms.Symbol
	ts := expr.R().TypeList()
	for i, tok := range ids {
		s := declareVar(b, tok, ts[i], false)
		if s == nil {
			return nil
		}
		ret = append(ret, s)
	}

	return &tast.Define{ret, expr}
}

func buildIdentExprList(b *builder, list *ast.ExprList) (
	idents []*lexing.Token, firstError ast.Expr,
) {
	ret := make([]*lexing.Token, 0, list.Len())
	for _, expr := range list.Exprs {
		op, ok := expr.(*ast.Operand)
		if !ok {
			return nil, expr
		}
		if op.Token.Type != parse.Ident {
			return nil, expr
		}
		ret = append(ret, op.Token)
	}

	return ret, nil
}

func buildDefineStmt(b *builder, stmt *ast.DefineStmt) tast.Stmt {
	right := b.buildExpr(stmt.Right)
	if right == nil {
		return nil
	}

	idents, err := buildIdentExprList(b, stmt.Left)
	if err != nil {
		b.Errorf(ast.ExprPos(err), "left side of := must be identifier")
		return nil
	}
	ret := define(b, idents, right, stmt.Define)
	if ret == nil {
		return nil
	}
	return ret
}

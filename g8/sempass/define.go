package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

// allocPrepare checks if the provided types are all allocable, and insert
// implicit type casts if needed. Only literay expression list needs alloc
// prepare.
func allocPrepare(
	b *Builder, toks []*lex8.Token, lst *tast.ExprList,
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

func declareVar(b *Builder, tok *lex8.Token, t types.T) *sym8.Symbol {
	name := tok.Lit
	s := sym8.Make(b.path, name, tast.SymVar, nil, t, tok.Pos)
	conflict := b.scope.Declare(s)
	if conflict != nil {
		b.Errorf(tok.Pos, "%q already declared as a %s",
			name, tast.SymStr(conflict.Type),
		)
		return nil
	}
	return s
}

func define(
	b *Builder, ids []*lex8.Token, expr tast.Expr, eq *lex8.Token,
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

	var syms []*sym8.Symbol
	ts := expr.R().TypeList()
	for i, tok := range ids {
		s := declareVar(b, tok, ts[i])
		if s == nil {
			return nil
		}
		syms = append(syms, s)
	}

	return &tast.Define{syms, expr}
}

func buildIdentExprList(b *Builder, list *ast.ExprList) (
	idents []*lex8.Token, firstError ast.Expr,
) {
	ret := make([]*lex8.Token, 0, list.Len())
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

func buildDefineStmt(b *Builder, stmt *ast.DefineStmt) tast.Stmt {
	right := b.BuildExpr(stmt.Right)
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

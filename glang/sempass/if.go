package sempass

import (
	"e8vm.io/e8vm/glang/ast"
	"e8vm.io/e8vm/glang/tast"
)

func buildIfStmt(b *builder, stmt *ast.IfStmt) tast.Stmt {
	return buildIf(b, stmt.Expr, stmt.Body, stmt.Else)
}

func buildIf(
	b *builder, cond ast.Expr, ifs ast.Stmt, elses *ast.ElseStmt,
) tast.Stmt {
	c := b.buildExpr(cond)
	if c == nil {
		return nil
	}

	condRef := c.R()
	if !condRef.IsBool() {
		pos := ast.ExprPos(cond)
		b.Errorf(pos, "expect boolean expression, got %s", condRef)
		return nil
	}

	if elses == nil {
		ret := &tast.IfStmt{Expr: c}
		var body tast.Stmt
		switch ifs := ifs.(type) {
		case *ast.Block:
			body = buildBlock(b, ifs)
		case *ast.ReturnStmt:
			body = buildReturnStmt(b, ifs)
		case *ast.BreakStmt:
			body = buildBreakStmt(b, ifs)
		case *ast.ContinueStmt:
			body = buildContinueStmt(b, ifs)
		default:
			pos := ast.ExprPos(cond)
			b.Errorf(pos, "if only takes block, return, break and continue")
		}
		if body == nil {
			return nil
		}
		ret.Body = body
		return ret
	}

	body := buildBlock(b, ifs.(*ast.Block))
	next := buildElseStmt(b, elses)
	return &tast.IfStmt{c, body, next}
}

func buildElseStmt(b *builder, stmt *ast.ElseStmt) tast.Stmt {
	if stmt.If == nil {
		if stmt.Expr != nil {
			panic("invalid expression in else")
		}
		return buildBlock(b, stmt.Body)
	}
	return buildIf(b, stmt.Expr, stmt.Body, stmt.Next)
}

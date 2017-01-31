package pl

import (
	"fmt"

	"shanhu.io/smlvm/pl/tast"
)

func buildStmt(b *builder, s tast.Stmt) {
	switch stmt := s.(type) {
	case nil:
		return // empty statement
	case *tast.ContinueStmt:
		buildContinueStmt(b)
	case *tast.BreakStmt:
		buildBreakStmt(b)
	case *tast.IncStmt:
		buildIncStmt(b, stmt)
	case *tast.ExprStmt:
		b.buildExpr(stmt.Expr)
	case *tast.Define:
		buildDefine(b, stmt)
	case *tast.VarDecls:
		for _, d := range stmt.Decls {
			buildDefine(b, d)
		}
	case *tast.ConstDecls:
		for _, d := range stmt.Decls {
			buildConstDefine(b, d)
		}
	case *tast.AssignStmt:
		buildAssignStmt(b, stmt)
	case *tast.ReturnStmt:
		buildReturnStmt(b, stmt)
	case *tast.Block:
		buildBlock(b, stmt)
	case *tast.ForStmt:
		buildForStmt(b, stmt)
	case *tast.IfStmt:
		buildIfStmt(b, stmt)
	default:
		panic(fmt.Errorf("unimplemented: %T", stmt))
	}
}

func buildBlock(b *builder, block *tast.Block) {
	for _, s := range block.Stmts {
		buildStmt(b, s)
	}
}

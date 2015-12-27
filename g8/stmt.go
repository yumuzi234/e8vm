package g8

import (
	"fmt"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

func buildStmt2(b *builder, stmt ast.Stmt) {
	s := b.spass.BuildStmt(stmt)
	if s == nil {
		return
	}

	switch stmt := s.(type) {
	case nil:
		return // empty statement
	case *tast.ContinueStmt:
		genContinueStmt(b)
	case *tast.BreakStmt:
		genBreakStmt(b)
	case *tast.IncStmt:
		buildIncStmt(b, stmt)
	case *tast.ExprStmt:
		b.buildExpr2(stmt.Expr)
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
		genAssignStmt(b, stmt)
	default:
		panic(fmt.Errorf("unimplemented: %T", stmt))
	}
}

func buildStmt(b *builder, stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.EmptyStmt:
		buildStmt2(b, stmt)
	case *ast.ExprStmt:
		buildStmt2(b, stmt)
	case *ast.IncStmt:
		buildStmt2(b, stmt)

	case *ast.DefineStmt:
		buildStmt2(b, stmt)
	case *ast.VarDecls:
		buildStmt2(b, stmt)
	case *ast.ConstDecls:
		buildStmt2(b, stmt)

	case *ast.ContinueStmt:
		buildContinueStmt(b, stmt)
	case *ast.BreakStmt:
		buildBreakStmt(b, stmt)

	case *ast.AssignStmt:
		buildAssignStmt(b, stmt)
	case *ast.ReturnStmt:
		buildReturnStmt(b, stmt)

	case *ast.IfStmt:
		buildIfStmt(b, stmt)
	case *ast.ForStmt:
		buildForStmt(b, stmt)
	case *ast.BlockStmt:
		buildBlock(b, stmt.Block)
	default:
		b.Errorf(nil, "invalid or not implemented: %T", stmt)
	}
}

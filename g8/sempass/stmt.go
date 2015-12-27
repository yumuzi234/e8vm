package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

func buildStmt(b *Builder, stmt ast.Stmt) tast.Stmt {
	switch stmt := stmt.(type) {
	case *ast.EmptyStmt:
		return nil
	case *ast.ExprStmt:
		return buildExprStmt(b, stmt.Expr)
	case *ast.IncStmt:
		return buildIncStmt(b, stmt)
	case *ast.DefineStmt:
		return buildDefineStmt(b, stmt)
	case *ast.VarDecls:
		return buildVarDecls(b, stmt)
	case *ast.AssignStmt:
		panic("todo")
	case *ast.ConstDecls:
		panic("todo")
	case *ast.IfStmt:
		panic("todo")
	case *ast.ForStmt:
		panic("todo")
	case *ast.ReturnStmt:
		panic("todo")
	case *ast.ContinueStmt:
		return buildContinueStmt(b, stmt)
	case *ast.BreakStmt:
		return buildBreakStmt(b, stmt)
	}

	b.Errorf(nil, "invalid or not implemented: %T", stmt)
	return nil
}

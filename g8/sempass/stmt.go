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
		return buildExprStmt(b, stmt)
	case *ast.IncStmt:
		panic("todo")
	case *ast.DefineStmt:
		panic("todo")
	case *ast.AssignStmt:
		panic("todo")
	case *ast.VarDecls:
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
		// TODO: check if in a loop
		return &tast.ContinueStmt{}
	case *ast.BreakStmt:
		// TODO: check if in a loop
		return &tast.BreakStmt{}
	}

	b.Errorf(nil, "invalid or not implemented: %T", stmt)
	return nil
}

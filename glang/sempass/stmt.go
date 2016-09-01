package sempass

import (
	"e8vm.io/e8vm/glang/ast"
	"e8vm.io/e8vm/glang/tast"
)

func buildStmt(b *builder, stmt ast.Stmt) tast.Stmt {
	switch stmt := stmt.(type) {
	case *ast.EmptyStmt:
		return nil
	case *ast.ExprStmt:
		return buildExprStmt(b, stmt.Expr)
	case *ast.IncStmt:
		return buildIncStmt(b, stmt)

	case *ast.ContinueStmt:
		return buildContinueStmt(b, stmt)
	case *ast.BreakStmt:
		return buildBreakStmt(b, stmt)

	case *ast.DefineStmt:
		return buildDefineStmt(b, stmt)
	case *ast.VarDecls:
		return buildVarDecls(b, stmt)
	case *ast.ConstDecls:
		return buildConstDecls(b, stmt)
	case *ast.AssignStmt:
		return buildAssignStmt(b, stmt)
	case *ast.ReturnStmt:
		return buildReturnStmt(b, stmt)
	case *ast.BlockStmt:
		return buildBlockStmt(b, stmt)
	case *ast.IfStmt:
		return buildIfStmt(b, stmt)
	case *ast.ForStmt:
		return buildForStmt(b, stmt)
	}

	b.Errorf(nil, "invalid or not implemented: %T", stmt)
	return nil
}

func buildStmts(b *builder, stmts []ast.Stmt) []tast.Stmt {
	var ret []tast.Stmt
	for _, stmt := range stmts {
		s := buildStmt(b, stmt)
		if s != nil {
			ret = append(ret, s)
		}
	}
	return ret
}

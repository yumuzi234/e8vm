package gfmt

import (
	"e8vm.io/e8vm/g8/ast"
)

func printStmt(f *formatter, stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.EmptyStmt:
		// empty, print nothing
	case *ast.Block:
		if len(stmt.Stmts) > 0 {
			f.printToken(stmt.Lbrace)
			f.printEndl()
			f.Tab()
			for _, s := range stmt.Stmts {
				printStmt(f, s)
				f.printEndlPlus(true, false)
			}
			f.cueTo(stmt.Rbrace)
			f.ShiftTab()
			f.printToken(stmt.Rbrace)
		} else {
			f.printToken(stmt.Lbrace)
			f.printToken(stmt.Rbrace)
		}
	case *ast.BlockStmt:
		printStmt(f, stmt.Block)
	case *ast.IfStmt:
		f.printExprs(stmt.If, " ", stmt.Expr, " ")
		printStmt(f, stmt.Body)
		if stmt.Else != nil {
			printStmt(f, stmt.Else)
		}
	case *ast.ElseStmt:
		f.printExprs(" ", stmt.Else, " ")
		if stmt.If == nil {
			printStmt(f, stmt.Body)
		} else {
			f.printExprs(stmt.If, " ", stmt.Expr, " ")
			printStmt(f, stmt.Body)
			if stmt.Next != nil {
				printStmt(f, stmt.Next)
			}
		}
	case *ast.ForStmt:
		f.printExprs(stmt.Kw, " ")
		if stmt.ThreeFold {
			if stmt.Init != nil {
				printStmt(f, stmt.Init)
			}
			f.printStr("; ")
			if stmt.Cond != nil {
				printExpr(f, stmt.Cond)
			}
			f.printStr("; ")
			printStmt(f, stmt.Iter)
			f.printSpace()
		} else if stmt.Cond != nil {
			f.printExprs(stmt.Cond, " ")
		}
		printStmt(f, stmt.Body)
	case *ast.AssignStmt:
		f.printExprs(stmt.Left, " ", stmt.Assign, " ", stmt.Right)
	case *ast.DefineStmt:
		f.printExprs(stmt.Left, " ", stmt.Define, " ", stmt.Right)
	case *ast.ExprStmt:
		printExpr(f, stmt.Expr)
	case *ast.IncStmt:
		f.printExprs(stmt.Expr, stmt.Op)
	case *ast.ReturnStmt:
		f.printToken(stmt.Kw)
		if stmt.Exprs != nil {
			f.printExprs(" ", stmt.Exprs)
		}
	case *ast.ContinueStmt:
		f.printToken(stmt.Kw)
		if stmt.Label != nil {
			f.printExprs(" ", stmt.Label)
		}
	case *ast.BreakStmt:
		f.printToken(stmt.Kw)
		if stmt.Label != nil {
			f.printExprs(" ", stmt.Label)
		}
	case *ast.VarDecls:
		printVarDecls(f, stmt)
	case *ast.ConstDecls:
		printConstDecls(f, stmt)
	default:
		f.errorf(nil, "invalid statement type: %T", stmt)
	}
}

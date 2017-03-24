package gfmt

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
)

func sameLine(t1, t2 *lexing.Token) bool {
	return t1.Pos.Line == t2.Pos.Line
}

func printStmt(f *formatter, stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.EmptyStmt:
		// empty, print nothing
	case *ast.Block:
		if !sameLine(stmt.Lbrace, stmt.Rbrace) || len(stmt.Stmts) > 0 {
			f.printToken(stmt.Lbrace)
			f.printEndl()
			f.Tab()
			for _, s := range stmt.Stmts {
				printStmt(f, s)
				f.printEndPara()
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
		body, ok := stmt.Body.(*ast.Block)
		if ok && stmt.Else != nil && len(body.Stmts) == 0 {
			f.printToken(body.Lbrace)
			f.printEndl()
			f.cueTo(body.Rbrace)
			f.printToken(body.Rbrace)
		} else {
			printStmt(f, stmt.Body)
		}
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
			if stmt.Iter != nil {
				printStmt(f, stmt.Iter)
			}
			f.printSpace()
		} else if stmt.Cond != nil {
			f.printExprs(stmt.Cond, " ")
		}
		printStmt(f, stmt.Body)
	case *ast.SwitchStmt:
		f.printExprs(stmt.Kw, " ")
		if stmt.Expr != nil {
			f.printExprs(stmt.Expr, " ")
		}
		if !sameLine(stmt.Lbrace, stmt.Rbrace) || len(stmt.Cases) > 0 {
			f.printToken(stmt.Lbrace)
			f.printEndl()
			for _, c := range stmt.Cases {
				printCase(f, c)
			}
			f.cueTo(stmt.Rbrace)
			f.printToken(stmt.Rbrace)
		} else {
			f.printToken(stmt.Lbrace)
			f.printToken(stmt.Rbrace)
		}
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
	case *ast.FallthroughStmt:
		f.printToken(stmt.Kw)
	case *ast.VarDecls:
		printVarDecls(f, stmt)
	case *ast.ConstDecls:
		printConstDecls(f, stmt)
	default:
		f.errorf(nil, "invalid statement type: %T", stmt)
	}
}

func printCase(f *formatter, c *ast.Case) {
	f.printExprs(c.Kw)
	if c.Expr != nil {
		f.printExprs(" ", c.Expr)
	}
	f.printToken(c.Colon)
	f.printEndl()
	f.Tab()
	for _, s := range c.Stmts {
		printStmt(f, s)
		f.printEndPara()
	}
	if c.Fallthrough != nil {
		printStmt(f, c.Fallthrough)
		f.printEndPara()
	}
	f.ShiftTab()
}

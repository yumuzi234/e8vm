package gfmt

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/g8/ast"
)

func printStmt(f *formatter, stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.EmptyStmt:
		f.printStr("; // empty")
	case *ast.Block:
		if len(stmt.Stmts) > 0 {
			f.printToken(stmt.Lbrace)
			f.printEndl()
			f.Tab()
			for _, s := range stmt.Stmts {
				printStmt(f, s)
				f.printEndlPlus(true, 0)
			}
			f.ShiftTab()
			f.printToken(stmt.Rbrace)
		} else {
			printExprs(f, stmt.Lbrace, stmt.Rbrace)
		}
	case *ast.BlockStmt:
		printStmt(f, stmt.Block)
	case *ast.IfStmt:
		printExprs(f, stmt.If, " ", stmt.Expr, " ")
		printStmt(f, stmt.Body)
		if stmt.Else != nil {
			printStmt(f, stmt.Else)
		}
	case *ast.ElseStmt:
		printExprs(f, " ", stmt.Else, " ")
		if stmt.If == nil {
			printStmt(f, stmt.Body)
		} else {
			printExprs(f, stmt.If, " ", stmt.Expr, " ")
			printStmt(f, stmt.Body)
			if stmt.Next != nil {
				printStmt(f, stmt.Next)
			}
		}
	case *ast.ForStmt:
		printExprs(f, stmt.Kw, " ")
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
			printExprs(f, stmt.Cond, " ")
		}
		printStmt(f, stmt.Body)
	case *ast.AssignStmt:
		printExprs(f, stmt.Left, " ", stmt.Assign, " ", stmt.Right)
	case *ast.DefineStmt:
		printExprs(f, stmt.Left, " ", stmt.Define, " ", stmt.Right)
	case *ast.ExprStmt:
		printExpr(f, stmt.Expr)
	case *ast.IncStmt:
		printExprs(f, stmt.Expr, stmt.Op)
	case *ast.ReturnStmt:
		f.printToken(stmt.Kw)
		if stmt.Exprs != nil {
			printExprs(f, " ", stmt.Exprs)
		}
	case *ast.ContinueStmt:
		f.printToken(stmt.Kw)
		if stmt.Label != nil {
			printExprs(f, " ", stmt.Label)
		}
	case *ast.BreakStmt:
		f.printToken(stmt.Kw)
		if stmt.Label != nil {
			printExprs(f, " ", stmt.Label)
		}
	//case *FallthroughStmt:
	// fmt.Fprint(p, "fallthrough")
	case *ast.VarDecls:
		printVarDecls(f, stmt)
	case *ast.ConstDecls:
		printConstDecls(f, stmt)
	default:
		panic(fmt.Errorf("invalid statement type: %T", stmt))
	}
}

// FprintStmts prints the statements out to a writer
func FprintStmts(out io.Writer, stmts []ast.Stmt) {
	f := newFormatter(out, nil)
	printStmt(f, stmts)
}

package gfmt

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
)

func printStmt(p *fmt8.Printer, m *matcher, stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.EmptyStmt:
		fmt.Fprint(p, "; // empty")
	case *ast.Block:
		if len(stmt.Stmts) > 0 {
			printToken(p, m, stmt.Lbrace)
			fmt.Fprintln(p)
			p.Tab()
			printStmt(p, m, stmt.Stmts)
			p.ShiftTab()
			printToken(p, m, stmt.Rbrace)
		} else {
			printExprs(p, m, stmt.Lbrace, stmt.Rbrace)
		}
	case *ast.BlockStmt:
		printStmt(p, m, stmt.Block)
	case []ast.Stmt:
		for _, s := range stmt {
			printStmt(p, m, s)
			fmt.Fprintln(p)
		}
	case *ast.IfStmt:
		printExprs(p, m, stmt.If, " ", stmt.Expr, " ")
		printStmt(p, m, stmt.Body)
		if stmt.Else != nil {
			printStmt(p, m, stmt.Else)
		}
	case *ast.ElseStmt:
		printExprs(p, m, " ", stmt.Else, " ")
		if stmt.If == nil {
			printStmt(p, m, stmt.Body)
		} else {
			printExprs(p, m, stmt.If, " ", stmt.Expr, " ")
			printStmt(p, m, stmt.Body)
			if stmt.Next != nil {
				printStmt(p, m, stmt.Next)
			}
		}
	case *ast.ForStmt:
		printExprs(p, m, stmt.Kw, " ")
		if stmt.ThreeFold {
			if stmt.Init != nil {
				printStmt(p, m, stmt.Init)
			}
			fmt.Fprint(p, "; ")
			if stmt.Cond != nil {
				printExpr(p, m, stmt.Cond)
			}
			fmt.Fprint(p, "; ")
			printStmt(p, m, stmt.Iter)
			fmt.Fprint(p, " ")
		} else if stmt.Cond != nil {
			printExprs(p, m, stmt.Cond, " ")
		}
		printStmt(p, m, stmt.Body)
	case *ast.AssignStmt:
		printExprs(p, m, stmt.Left, " ", stmt.Assign, " ", stmt.Right)
	case *ast.DefineStmt:
		printExprs(p, m, stmt.Left, " ", stmt.Define, " ", stmt.Right)
	case *ast.ExprStmt:
		printExpr(p, m, stmt.Expr)
	case *ast.IncStmt:
		printExprs(p, m, stmt.Expr, stmt.Op)
	case *ast.ReturnStmt:
		printToken(p, m, stmt.Kw)
		if stmt.Exprs != nil {
			printExprs(p, m, " ", stmt.Exprs)
		}
	case *ast.ContinueStmt:
		printToken(p, m, stmt.Kw)
		if stmt.Label != nil {
			printExprs(p, m, " ", stmt.Label)
		}
	case *ast.BreakStmt:
		printToken(p, m, stmt.Kw)
		if stmt.Label != nil {
			printExprs(p, m, " ", stmt.Label)
		}
	//case *FallthroughStmt:
	// fmt.Fprint(p, "fallthrough")
	case *ast.VarDecls:
		printVarDecls(p, m, stmt)
	case *ast.ConstDecls:
		printConstDecls(p, m, stmt)
	default:
		panic(fmt.Errorf("invalid statement type: %T", stmt))
	}
}

// FprintStmts prints the statements out to a writer
func FprintStmts(out io.Writer, stmts []ast.Stmt) {
	p := fmt8.NewPrinter(out)
	printStmt(p, nil, stmts)
}

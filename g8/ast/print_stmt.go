package ast

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/fmt8"
)

func printStmt(p *fmt8.Printer, stmt Stmt) {
	switch stmt := stmt.(type) {
	case *EmptyStmt:
		fmt.Fprint(p, "; // empty")
	case *Block:
		if len(stmt.Stmts) > 0 {
			fmt.Fprintln(p, "{")
			p.Tab()
			printStmt(p, stmt.Stmts)
			p.ShiftTab()
			fmt.Fprint(p, "}")
		} else {
			fmt.Fprint(p, "{}")
		}
	case *BlockStmt:
		printStmt(p, stmt.Block)
	case []Stmt:
		for _, s := range stmt {
			printStmt(p, s)
			fmt.Fprintln(p)
		}
	case *IfStmt:
		printExprs(p, "if ", stmt.Expr, " ")
		printStmt(p, stmt.Body)
		if stmt.Else != nil {
			printStmt(p, stmt.Else)
		}
	case *ElseStmt:
		if stmt.If == nil {
			printExprs(p, " else ")
			printStmt(p, stmt.Body)
		} else {
			printExprs(p, " else if ", stmt.Expr, " ")
			printStmt(p, stmt.Body)
			if stmt.Next != nil {
				printStmt(p, stmt.Next)
			}
		}
	case *ForStmt:
		fmt.Fprint(p, "for ")
		if stmt.ThreeFold {
			if stmt.Init != nil {
				printStmt(p, stmt.Init)
			}
			fmt.Fprint(p, "; ")
			if stmt.Cond != nil {
				printExprs(p, stmt.Cond)
			}
			fmt.Fprint(p, "; ")
			printStmt(p, stmt.Iter)
			fmt.Fprint(p, " ")
		} else if stmt.Cond != nil {
			printExprs(p, stmt.Cond, " ")
		}
		printStmt(p, stmt.Body)
	case *AssignStmt:
		printExprs(p, stmt.Left, " = ", stmt.Right)
	case *DefineStmt:
		printExprs(p, stmt.Left, " := ", stmt.Right)
	case *ExprStmt:
		printExprs(p, stmt.Expr)
	case *IncStmt:
		printExprs(p, stmt.Expr, stmt.Op.Lit)
	case *ReturnStmt:
		if stmt.Exprs != nil {
			printExprs(p, "return ", stmt.Exprs)
		} else {
			printExprs(p, "return")
		}
	case *ContinueStmt:
		if stmt.Label == nil {
			printExprs(p, "continue")
		} else {
			printExprs(p, "continue ", stmt.Label.Lit)
		}
	case *BreakStmt:
		if stmt.Label == nil {
			printExprs(p, "break")
		} else {
			printExprs(p, "break ", stmt.Label.Lit)
		}
	//case *FallthroughStmt:
	// fmt.Fprint(p, "fallthrough")
	case *VarDecls:
		printVarDecls(p, stmt)
	case *ConstDecls:
		printConstDecls(p, stmt)
	default:
		panic(fmt.Errorf("invalid statement type: %T", stmt))
	}
}

// FprintStmts prints the statements out to a writer
func FprintStmts(out io.Writer, stmts []Stmt) {
	p := fmt8.NewPrinter(out)
	printStmt(p, stmts)
}

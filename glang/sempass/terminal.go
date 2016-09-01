package sempass

import (
	"e8vm.io/e8vm/glang/ast"
)

func isBlockTerminal(block *ast.Block) bool {
	stmts := block.Stmts
	nstmt := len(stmts)
	if nstmt == 0 {
		return false
	}
	return isTerminal(stmts[nstmt-1])
}

func blockHasBreak(b *ast.Block) bool {
	for _, s := range b.Stmts {
		if hasBreak(s) {
			return true
		}
	}
	return false
}

func hasBreak(stmt ast.Stmt) bool {
	switch stmt := stmt.(type) {
	case *ast.BlockStmt:
		return blockHasBreak(stmt.Block)
	case *ast.Block:
		return blockHasBreak(stmt)
	case *ast.BreakStmt:
		return true
	case *ast.ForStmt:
		// TODO(h8liu): need to change this if labeled break is added
		return false
	case *ast.IfStmt:
		if hasBreak(stmt.Body) {
			return true
		}
		selse := stmt.Else
		for selse != nil {
			for _, s := range selse.Body.Stmts {
				if hasBreak(s) {
					return true
				}
			}
			selse = selse.Next
		}
		return false
	default:
		return false
	}
}

func isTerminal(stmt ast.Stmt) bool {
	switch stmt := stmt.(type) {
	case *ast.BlockStmt:
		return isBlockTerminal(stmt.Block)
	case *ast.Block:
		return isBlockTerminal(stmt)
	case *ast.IfStmt:
		if stmt.Else == nil {
			return false
		}
		if !isTerminal(stmt.Body) {
			return false
		}
		selse := stmt.Else
		for selse != nil {
			if selse.Else == nil && selse.If != nil {
				// else if with no further else
				return false
			}
			if !isBlockTerminal(selse.Body) {
				return false
			}
			selse = selse.Next
		}
		return true
	case *ast.ForStmt:
		if stmt.Cond != nil {
			return false
		}
		for _, s := range stmt.Body.Stmts {
			if hasBreak(s) {
				return false
			}
		}
		return true
	case *ast.ReturnStmt:
		return true
	default:
		return false
	}
}

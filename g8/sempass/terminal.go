package sempass

import (
	"e8vm.io/e8vm/g8/ast"
)

func isBlockTerminal(block *ast.Block) bool {
	stmts := block.Stmts
	nstmt := len(stmts)
	if nstmt == 0 {
		return false
	}
	return isTerminal(stmts[nstmt-1])
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
	case *ast.ReturnStmt:
		return true
	default:
		return false
	}
}

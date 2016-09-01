package asm

import (
	"e8vm.io/e8vm/asm/ast"
	"e8vm.io/e8vm/lexing"
)

type funcStmt struct {
	*ast.FuncStmt

	*inst
	label  string
	offset uint32
}

func resolveFuncStmt(log lexing.Logger, s *ast.FuncStmt) *funcStmt {
	ops := s.Ops
	op0 := ops[0]
	lead := op0.Lit
	if lead == "" {
		panic("empty operand")
	}

	if checkLabel(log, op0) {
		if len(ops) > 1 {
			log.Errorf(op0.Pos, "label should take the entire line")
			return nil
		}
		return &funcStmt{label: lead, FuncStmt: s}
	}

	return &funcStmt{inst: resolveInst(log, ops), FuncStmt: s}
}

func (s *funcStmt) isLabel() bool {
	return s.inst == nil && s.label != ""
}

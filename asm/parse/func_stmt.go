package parse

import (
	"shanhu.io/smlvm/asm/ast"
	"shanhu.io/smlvm/lexing"
)

func parseOps(p *parser) (ops []*lexing.Token) {
	for !p.Accept(Semi) {
		t := p.Expect(Operand)
		if t == nil {
			p.skipErrStmt()
			return nil
		}

		ops = append(ops, t)
	}

	return ops
}

func parseFuncStmt(p *parser) *ast.FuncStmt {
	ops := parseOps(p)
	if len(ops) == 0 {
		return nil
	}

	return &ast.FuncStmt{Ops: ops}
}

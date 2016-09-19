package parse

import (
	"shanhu.io/smlvm/asm/ast"
	"shanhu.io/smlvm/lexing"
)

func parseArgs(p *parser) (typ *lexing.Token, args []*lexing.Token) {
	typ = p.Expect(Operand)
	if typ == nil {
		p.skipErrStmt()
		return nil, nil
	}

	for !p.Accept(Semi) {
		if !p.InError() {
			t := p.Token()
			if t.Type == Operand || t.Type == String {
				args = append(args, t)
			} else {
				p.Errorf(t.Pos, "expect operand or string, got %s",
					p.TypeStr(t.Type),
				)
			}
		}
		if p.See(lexing.EOF) {
			break
		}
		p.Next()
	}

	p.BailOut()

	return typ, args
}

func parseVarStmt(p *parser) *ast.VarStmt {
	typ, args := parseArgs(p)
	if typ == nil {
		return nil
	}

	ret := new(ast.VarStmt)
	ret.Type = typ
	ret.Args = args

	return ret
}

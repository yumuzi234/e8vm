package parse

import (
	"io"

	"e8vm.io/e8vm/asm8/ast"
	"e8vm.io/e8vm/lexing"
)

func parseFuncStmts(p *parser, f *ast.Func) {
	for !(p.See(Rbrace) || p.See(lexing.EOF)) {
		stmt := parseFuncStmt(p)
		if stmt != nil {
			f.Stmts = append(f.Stmts, stmt)
		}
	}
}

func parseBareFunc(p *parser) *ast.Func {
	ret := new(ast.Func)
	ret.Name = &lexing.Token{
		Type: Operand,
		Lit:  "_",
		Pos:  nil,
	}
	parseFuncStmts(p, ret)
	return ret
}

// BareFunc parses a file as a bare function.
func BareFunc(f string, rc io.ReadCloser) (*ast.Func, []*lexing.Error) {
	p, _ := newParser(f, rc)
	fn := parseBareFunc(p)
	if es := p.Errs(); es != nil {
		return nil, es
	}

	return fn, nil
}

func parseFunc(p *parser) *ast.Func {
	ret := new(ast.Func)

	ret.Kw = p.ExpectKeyword("func")
	ret.Name = p.Expect(Operand)

	if ret.Name != nil {
		name := ret.Name.Lit
		if !IsIdent(name) {
			p.Errorf(ret.Name.Pos, "invalid func name %q", name)
		}
	}

	ret.Lbrace = p.Expect(Lbrace)
	if p.skipErrStmt() { // header broken
		return ret
	}

	parseFuncStmts(p, ret)

	ret.Rbrace = p.Expect(Rbrace)
	ret.Semi = p.Expect(Semi)
	p.skipErrStmt()

	return ret
}

package parse

import (
	"io"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/lex8"
)

func parseStmt(p *parser) ast.Stmt {
	first := p.Token()
	if first.Type == Keyword {
		switch first.Lit {
		case "const":
			return parseConstDecls(p)
		case "var":
			return parseVarDecls(p)
		case "if":
			return parseIfStmt(p)
		case "for":
			return parseForStmt(p)
		//case "switch":
		//	return parseSwitchStmt(p)
		case "return":
			return parseReturnStmt(p, true)
		case "break":
			return parseBreakStmt(p, true)
		case "continue":
			return parseContinueStmt(p, true)
		}
	}

	if p.SeeOp("{") {
		ret := new(ast.BlockStmt)
		ret.Block = parseBlock(p)
		ret.Semi = p.ExpectSemi()
		return ret
	}

	return parseSimpleStmt(p)
}

func makeParser(f string, r io.Reader, golike bool) (
	*parser, *lex8.Recorder,
) {
	p, rec := newParser(f, r, golike)
	p.exprFunc = parseExpr
	p.stmtFunc = parseStmt
	p.typeFunc = parseType
	p.seeTypeFunc = seeType
	return p, rec
}

// Stmts parses a file input stream as a list of statements,
// like a bare function body.
func Stmts(f string, r io.Reader) ([]ast.Stmt, []*lex8.Error) {
	p, _ := makeParser(f, r, false)

	var ret []ast.Stmt
	for !p.See(lex8.EOF) {
		if stmt := p.parseStmt(); stmt != nil {
			ret = append(ret, stmt)
		}
		p.skipErrStmt()
	}

	if es := p.Errs(); es != nil {
		return nil, es
	}

	return ret, nil
}

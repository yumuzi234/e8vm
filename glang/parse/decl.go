package parse

import (
	"e8vm.io/e8vm/glang/ast"
	"e8vm.io/e8vm/lexing"
)

func parseIdentList(p *parser) *ast.IdentList {
	if !p.See(Ident) {
		p.Expect(Ident)
		return nil
	}

	ret := new(ast.IdentList)
	for p.See(Ident) {
		ret.Idents = append(ret.Idents, p.Shift())
		if !p.SeeOp(",") {
			break
		}
		ret.Commas = append(ret.Commas, p.Shift())
	}
	return ret
}

func parseConstDecls(p *parser) *ast.ConstDecls {
	if !p.SeeKeyword("const") {
		panic("expect keyword")
	}

	ret := new(ast.ConstDecls)
	ret.Kw = p.Shift()

	if p.SeeOp("(") {
		ret.Lparen = p.Shift()
		for !p.See(lexing.EOF) && !p.SeeOp(")", "}") {
			if !p.See(Ident) {
				p.Expect(Ident)
				p.skipErrStmt()
				continue
			}
			d := parseConstDecl(p)
			if d != nil {
				ret.Decls = append(ret.Decls, d)
			} else {
				p.skipErrStmt()
			}
		}
		ret.Rparen = p.ExpectOp(")")
		ret.Semi = p.ExpectSemi()
		return ret
	}

	d := parseConstDecl(p)
	ret.Decls = []*ast.ConstDecl{d}
	return ret
}

func parseConstDecl(p *parser) *ast.ConstDecl {
	ret := new(ast.ConstDecl)
	ret.Idents = parseIdentList(p)
	if p.InError() {
		return nil
	}

	if !p.SeeOp("=") {
		ret.Type = p.parseType()
	}

	ret.Eq = p.ExpectOp("=")
	if p.InError() {
		return nil
	}

	ret.Exprs = parseExprList(p)
	if p.InError() {
		return nil
	}

	ret.Semi = p.ExpectSemi()
	if p.InError() {
		return nil
	}
	return ret
}

func parseVarDecl(p *parser) *ast.VarDecl {
	ret := new(ast.VarDecl)
	ret.Idents = parseIdentList(p)
	if p.InError() {
		return nil
	}

	if !p.See(Semi) && !p.SeeOp("=", ")", "}") {
		ret.Type = p.parseType() // it has a type
	}

	if p.SeeOp("=") {
		ret.Eq = p.Shift()
		ret.Exprs = parseExprList(p)
	} else if ret.Type == nil {
		p.ErrorfHere("expect type")
	}

	ret.Semi = p.ExpectSemi()
	if p.InError() {
		return nil
	}

	return ret
}

func parseVarDecls(p *parser) *ast.VarDecls {
	if !p.SeeKeyword("var") {
		panic("expect keyword")
	}

	ret := new(ast.VarDecls)
	ret.Kw = p.Shift()

	if p.SeeOp("(") {
		ret.Lparen = p.Shift()
		for !p.See(lexing.EOF) && !p.SeeOp(")", "}") {
			if !p.See(Ident) {
				p.Expect(Ident)
				p.skipErrStmt()
				continue
			}

			d := parseVarDecl(p)
			if d != nil {
				ret.Decls = append(ret.Decls, d)
			} else {
				p.skipErrStmt()
			}
		}
		ret.Rparen = p.ExpectOp(")")
		ret.Semi = p.ExpectSemi()

		return ret
	}

	d := parseVarDecl(p)
	ret.Decls = []*ast.VarDecl{d}
	return ret // no semi means it is not a group
}

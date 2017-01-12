package parse

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
)

func parsePara(p *parser) *ast.Para {
	ret := new(ast.Para)
	if p.See(Ident) {
		ident := p.Shift()
		if p.SeeOp(".") {
			ret.Type = parseMemberExpr(p, &ast.Operand{ident})
		} else {
			ret.Ident = ident
			if !(p.SeeOp(",") || p.SeeOp(")")) {
				ret.Type = p.parseType()
			}
		}
	} else {
		ret.Type = p.parseType()
	}
	return ret
}

func parseParaList(p *parser) *ast.ParaList {
	ret := new(ast.ParaList)
	ret.Lparen = p.ExpectOp("(")
	if p.InError() {
		return nil
	}

	if p.SeeOp(")") {
		// empty parameter list
		ret.Rparen = p.Shift()
		return ret
	}

	for !p.See(lexing.EOF) {
		para := parsePara(p)
		if p.InError() {
			return nil
		}

		ret.Paras = append(ret.Paras, para)
		if p.SeeOp(",") {
			ret.Commas = append(ret.Commas, p.Shift())
		} else if !p.SeeOp(")") {
			p.ExpectOp(",")
			return nil
		}

		if p.SeeOp(")") {
			break
		}
	}

	ret.Rparen = p.ExpectOp(")")
	return ret
}

func parseFuncSig(p *parser) *ast.FuncSig {
	ret := new(ast.FuncSig)
	ret.Args = parseParaList(p)
	if p.InError() {
		return nil
	}

	if p.SeeOp("(") {
		ret.Rets = parseParaList(p)
		if p.InError() {
			return nil
		}
		if len(ret.Rets.Paras) == 0 {
			p.CodeErrorf(ret.Rets.Rparen.Pos, "pl.expectReturnList",
				"expect return list in \"()\" after the function")
		}
	} else if p.SeeType() {
		ret.RetType = p.parseType()
	}
	return ret
}

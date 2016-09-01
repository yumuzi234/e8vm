package parse

import (
	"e8vm.io/e8vm/glang/ast"
	"e8vm.io/e8vm/lexing"
)

func parsePara(p *parser) *ast.Para {
	ret := new(ast.Para)
	if p.See(Ident) {
		ret.Ident = p.Shift()
		if !(p.SeeOp(",") || p.SeeOp(")")) {
			ret.Type = p.parseType()
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
			p.Errorf(ret.Rets.Rparen.Pos, "expect return list")
		}
	} else if p.SeeType() {
		ret.RetType = p.parseType()
	}
	return ret
}

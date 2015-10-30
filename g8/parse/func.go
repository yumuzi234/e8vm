package parse

import (
	"e8vm.io/e8vm/g8/ast"
)

func parseFunc(p *parser) *ast.Func {
	if !p.SeeKeyword("func") {
		panic("expect keyword")
	}

	ret := new(ast.Func)
	ret.Kw = p.Shift()

	// function receiver
	if p.golike && p.SeeOp("(") {
		recv := new(ast.FuncRecv)
		recv.Lparen = p.Shift()
		if !p.SeeOp("*") {
			recv.Recv = p.Expect(Ident)
		}
		recv.Star = p.ExpectOp("*")
		recv.StructName = p.Expect(Ident)
		recv.Rparen = p.ExpectOp(")")

		ret.Recv = recv

		if p.InError() {
			return nil
		}
	}

	ret.Name = p.Expect(Ident)
	if p.InError() {
		return nil
	}

	ret.FuncSig = parseFuncSig(p)
	if p.InError() {
		return nil
	}

	ret.Body = parseBlock(p)
	ret.Semi = p.ExpectSemi()
	return ret
}

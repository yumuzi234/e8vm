package parse

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
)

func parseInterface(p *parser) *ast.Interface {
	var ret *ast.Interface
	if p.SeeKeyword("interface") {
		ret = &ast.Interface{
			Kw:     p.ExpectKeyword("interface"),
			Name:   p.Expect(Ident),
			Lbrace: p.ExpectOp("{"),
		}
	} else {
		panic("expect keyword")
	}
	for !p.SeeOp("}") && !p.See(lexing.EOF) {

		name := p.Expect(Ident)
		if p.InError() {
			return nil
		}
		ret.Funcs = append(ret.Funcs, name)
		f := parseFuncSig(p)
		if p.InError() {
			return nil
		}
		ret.FuncSigs = append(ret.FuncSigs, f)
		p.ExpectSemi()
	}
	ret.Rbrace = p.ExpectOp("}")
	ret.Semi = p.ExpectSemi()
	return ret
}

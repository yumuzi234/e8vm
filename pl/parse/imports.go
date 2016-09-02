package parse

import (
	"e8vm.io/e8vm/lexing"
	"e8vm.io/e8vm/pl/ast"
)

func parseImports(p *parser) *ast.ImportDecls {
	if !p.SeeKeyword("import") {
		panic("expect keyword")
	}

	ret := &ast.ImportDecls{
		Kw:     p.ExpectKeyword("import"),
		Lparen: p.ExpectOp("("),
	}

	for !p.SeeOp(")", "}") && !p.See(lexing.EOF) {
		imp := new(ast.ImportDecl)
		if p.See(Ident) {
			imp.As = p.Shift()
		}
		imp.Path = p.Expect(String)
		imp.Semi = p.ExpectSemi()

		if p.skipErrStmt() {
			continue
		}

		ret.Decls = append(ret.Decls, imp)
	}

	ret.Rparen = p.ExpectOp(")")
	ret.Semi = p.ExpectSemi()

	return ret
}

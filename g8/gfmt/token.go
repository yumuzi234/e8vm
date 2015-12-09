package gfmt

import (
	"fmt"

	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/lex8"
)

func printToken(p *fmt8.Printer, m *matcher, token *lex8.Token) {
	if m != nil {
		m.expect(token)
	}
	fmt.Fprint(p, token.Lit)
}

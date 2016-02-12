package parse

import (
	"e8vm.io/e8vm/lex8"
)

func isLetter(r rune) bool {
	return r == '_' || lex8.IsLetter(r)
}

func lexNumber(x *lex8.Lexer) *lex8.Token {
	// TODO: lex floating point as well

	start := x.Rune()
	if !lex8.IsDigit(start) {
		panic("not starting with a number")
	}

	x.Next()

	r := x.Rune()
	if start == '0' && r == 'x' {
		x.Next()

		for lex8.IsHexDigit(x.Rune()) {
			x.Next()
		}
	} else {
		for lex8.IsDigit(x.Rune()) {
			x.Next()
		}
	}
	return x.MakeToken(Int)
}

func lexIdent(x *lex8.Lexer) *lex8.Token {
	r := x.Rune()
	if !isLetter(r) {
		panic("must start with letter")
	}

	for {
		x.Next()
		r := x.Rune()
		if !isLetter(r) && !lex8.IsDigit(r) {
			break
		}
	}

	ret := x.MakeToken(Ident)
	return ret
}

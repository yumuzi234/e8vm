package parse

import (
	"e8vm.io/e8vm/lexing"
)

func isLetter(r rune) bool {
	return r == '_' || lexing.IsLetter(r)
}

func lexNumber(x *lexing.Lexer) *lexing.Token {
	// TODO: lex floating point as well

	start := x.Rune()
	if !lexing.IsDigit(start) {
		panic("not starting with a number")
	}

	x.Next()

	r := x.Rune()
	if start == '0' && r == 'x' {
		x.Next()

		for lexing.IsHexDigit(x.Rune()) {
			x.Next()
		}
	} else {
		for lexing.IsDigit(x.Rune()) {
			x.Next()
		}
	}
	return x.MakeToken(Int)
}

func lexIdent(x *lexing.Lexer) *lexing.Token {
	r := x.Rune()
	if !isLetter(r) {
		panic("must start with letter")
	}

	for {
		x.Next()
		r := x.Rune()
		if !isLetter(r) && !lexing.IsDigit(r) {
			break
		}
	}

	ret := x.MakeToken(Ident)
	return ret
}

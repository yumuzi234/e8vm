package parse

import (
	"e8vm.io/e8vm/lex8"
)

func isLetter(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	if r >= 'A' && r <= 'Z' {
		return true
	}
	if r == '_' {
		return true
	}
	return false
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isHexDigit(r rune) bool {
	if isDigit(r) {
		return true
	}
	if r >= 'a' && r <= 'f' {
		return true
	}
	if r >= 'A' && r <= 'F' {
		return true
	}
	return false
}

func lexNumber(x *lex8.Lexer) *lex8.Token {
	// TODO: lex floating point as well

	start := x.Rune()
	if !isDigit(start) {
		panic("not starting with a number")
	}

	x.Next()

	r := x.Rune()
	if start == '0' && r == 'x' {
		x.Next()

		for isHexDigit(x.Rune()) {
			x.Next()
		}
	} else {
		for isDigit(x.Rune()) {
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
		if !isLetter(r) && !isDigit(r) {
			break
		}
	}

	ret := x.MakeToken(Ident)
	return ret
}

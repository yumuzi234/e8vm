package parse

import (
	"io"

	"shanhu.io/smlvm/lexing"
)

func lexAsm8(x *lexing.Lexer) *lexing.Token {
	r := x.Rune()
	if x.IsWhite(r) {
		panic("incorrect token start")
	}

	switch r {
	case '\n':
		x.Next()
		return x.MakeToken(Endl)
	case '{':
		x.Next()
		return x.MakeToken(Lbrace)
	case '}':
		x.Next()
		return x.MakeToken(Rbrace)
	case '/':
		x.Next()
		return lexing.LexComment(x)
	case '"':
		return lexing.LexString(x, String, '"')
	}

	if isOperandChar(r) {
		return lexOperand(x)
	}

	x.Errorf("illegal char %q", r)
	x.Next()
	return x.MakeToken(lexing.Illegal)
}

func newLexer(file string, r io.Reader) *lexing.Lexer {
	return lexing.MakeLexer(file, r, lexAsm8)
}

// Tokens parses a file in a token array
func Tokens(f string, r io.Reader) ([]*lexing.Token, []*lexing.Error) {
	x := newLexer(f, r)
	toks := lexing.TokenAll(x)
	if es := x.Errs(); es != nil {
		return nil, es
	}
	return toks, nil
}

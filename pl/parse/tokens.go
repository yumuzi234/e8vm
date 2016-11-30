package parse

import (
	"io"

	"shanhu.io/smlvm/lexing"
)

// Tokens parses a file into a token array
func Tokens(f string, r io.Reader) ([]*lexing.Token, []*lexing.Error) {
	x := makeTokener(f, r, false)
	toks := lexing.TokenAll(x)
	if errs := x.Errs(); errs != nil {
		return nil, errs
	}
	return toks, nil
}

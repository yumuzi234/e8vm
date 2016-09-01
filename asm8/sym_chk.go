package asm8

import (
	"strings"

	"e8vm.io/e8vm/asm8/parse"
	"e8vm.io/e8vm/lexing"
)

// mightBeSymbol just looks at the first rune and see
// if it is *possibly* a symbol
func mightBeSymbol(sym string) bool {
	if sym == "" {
		return false
	}
	r := sym[0]
	if r >= 'a' && r <= 'z' {
		return true
	}
	if r >= 'A' && r <= 'Z' {
		return true
	}
	return false
}

func parseSym(p lexing.Logger, t *lexing.Token) (pack, sym string) {
	if t.Type != parse.Operand {
		panic("symbol not an operand")
	}

	sym = t.Lit
	dot := strings.Index(sym, ".")
	if dot >= 0 {
		pack, sym = sym[:dot], sym[dot+1:]
	}

	if dot >= 0 && !lexing.IsPkgName(pack) {
		p.Errorf(t.Pos, "invalid package name: %q", pack)
	} else if !parse.IsIdent(sym) {
		p.Errorf(t.Pos, "invalid symbol: %q", t.Lit)
	}

	return pack, sym
}

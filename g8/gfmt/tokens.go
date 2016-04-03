package gfmt

import (
	"e8vm.io/e8vm/lex8"
)

type tokens struct {
	toks []*lex8.Token
	cur  int
}

func newTokens(toks []*lex8.Token) *tokens {
	return &tokens{toks: toks}
}

func (t *tokens) get(i int) *lex8.Token {
	if i < 0 {
		return nil
	}
	if i >= len(t.toks) {
		return nil
	}
	return t.toks[i]
}

func (t *tokens) peek() *lex8.Token {
	return t.get(t.cur)
}

func (t *tokens) shift() *lex8.Token {
	ret := t.get(t.cur)
	if ret == nil {
		return nil
	}
	t.cur++
	return ret
}

func (t *tokens) see(tok *lex8.Token) bool {
	return t.peek() == tok
}

func (t *tokens) lineGap() int {
	cur := t.peek()
	last := t.get(t.cur - 1)
	if cur == nil || last == nil {
		return 0
	}

	return cur.Pos.Line - last.Pos.Line
}

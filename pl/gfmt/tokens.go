package gfmt

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/parse"
)

type tokens struct {
	toks []*lexing.Token
	cur  int
}

func newTokens(toks []*lexing.Token) *tokens {
	return &tokens{toks: toks}
}

func (t *tokens) get(i int) *lexing.Token {
	if i < 0 {
		return nil
	}
	if i >= len(t.toks) {
		return nil
	}
	return t.toks[i]
}

func (t *tokens) peek() *lexing.Token {
	return t.get(t.cur)
}

func (t *tokens) shift() *lexing.Token {
	ret := t.get(t.cur)
	if ret == nil {
		return nil
	}
	t.cur++
	return ret
}

func (t *tokens) see(tok *lexing.Token) bool {
	return t.peek() == tok
}

func (t *tokens) lineGap() int {
	cur := t.peek()
	if cur != nil && cur.Type == parse.Semi && cur.Lit == "\n" {
		cur = t.get(t.cur + 1)
	}

	last := t.get(t.cur - 1)
	if cur == nil || last == nil {
		return 0
	}

	return cur.Pos.Line - last.Pos.Line
}

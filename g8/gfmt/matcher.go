package gfmt

import (
	"fmt"

	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/lex8"
)

type matcher struct {
	tokens []*lex8.Token
	offset int
}

func newMatcher(tokens []*lex8.Token) *matcher {
	return &matcher{tokens, 0}
}

func (m *matcher) expect(token *lex8.Token) {
	t := m.next()
	if t == nil {
		panic(fmt.Errorf("unexpected token: %v", token))
	}
	if t != token {
		panic(fmt.Errorf("unmatched token: got %v, expected %v", t, token))
	}
}

func (m *matcher) next() *lex8.Token {
	for m.offset < len(m.tokens) {
		token := m.tokens[m.offset]
		m.offset++
		if token.Type == parse.Semi {
			continue // ignore semi-s
		}
		if token.Type == lex8.Comment {
			// TODO(kcnm): Emit comments.
			continue
		}
		return token
	}
	return nil
}

func (m *matcher) finish() {
	token := m.next()
	if token.Type != lex8.EOF {
		panic(fmt.Errorf("unmatched token: got %v, expected EOF", token))
	}
	if m.offset < len(m.tokens) {
		panic(fmt.Errorf("unfinished tokens: %v", m.tokens[m.offset:]))
	}
}

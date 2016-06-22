package lex8

import (
	"strings"
	"testing"
)

func TestWordLexer(t *testing.T) {
	x := NewWordLexer("t1.txt", strings.NewReader("hello, world!"))
	var toks []*Token

	for {
		t := x.Token()
		toks = append(toks, t)
		if t.Type == EOF {
			break
		}
	}

	if len(toks) != 5 {
		t.Errorf("want 5 tokens, got %d", len(toks))
		return
	}

	for i, s := range []string{"hello", ",", "world", "!", ""} {
		lit := toks[i].Lit
		if s != lit {
			t.Errorf("token %d want %q, got %q", i, s, lit)
		}
	}

	for i, want := range []int{Word, Punc, Word, Punc, EOF} {
		typ := toks[i].Type
		if want != typ {
			t.Errorf("token %d want type %d, got %d", i, want, typ)
		}
	}

	x = NewWordLexer("t2.txt", strings.NewReader("123	#a $%\n\rd^"))
	toks = nil
	for {
		t := x.Token()
		toks = append(toks, t)
		if t.Type == EOF {
			break
		}
	}
	if len(toks) != 8 {
		t.Errorf("want 7 tokens, got %d", len(toks))
		return
	}
	for i, s := range []string{"123", "#", "a", "$", "%", "d", "^", ""} {
		lit := toks[i].Lit
		if s != lit {
			t.Errorf("token %d want %q, got %q", i, s, lit)
		}
	}

	for i, want := range []int{Word, Punc, Word, Punc, Punc, Word, Punc, EOF} {
		typ := toks[i].Type
		if want != typ {
			t.Errorf("token %d want type %d, got %d", i, want, typ)
		}
	}
}

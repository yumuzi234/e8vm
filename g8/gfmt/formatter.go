package gfmt

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/lex8"
)

type formatter struct {
	*fmt8.Printer
	tokens []*lex8.Token
	offset int
}

func newFormatter(out io.Writer, tokens []*lex8.Token) *formatter {
	p := fmt8.NewPrinter(out)
	return &formatter{p, tokens, 0}
}

func (f *formatter) printStr(s string) {
	fmt.Fprint(f.Printer, s)
}

func (f *formatter) printSpace() {
	f.printStr(" ")
}

func (f *formatter) printEndl() {
	fmt.Fprintln(f.Printer)
}

func (f *formatter) printToken(token *lex8.Token) {
	if f.tokens != nil {
		f.expect(token)
	}
	f.printStr(token.Lit)
}

func (f *formatter) expect(token *lex8.Token) {
	t := f.next()
	if t == nil {
		panic(fmt.Errorf("unexpected token: %v", token))
	}
	if t != token {
		panic(fmt.Errorf("unmatched token: got %v, expected %v", t, token))
	}
}

func (f *formatter) next() *lex8.Token {
	for f.offset < len(f.tokens) {
		token := f.tokens[f.offset]
		f.offset++
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

func (f *formatter) finish() {
	token := f.next()
	if token.Type != lex8.EOF {
		panic(fmt.Errorf("unmatched token: got %v, expected EOF", token))
	}
	if f.offset < len(f.tokens) {
		panic(fmt.Errorf("unfinished tokens: %v", f.tokens[f.offset:]))
	}
}

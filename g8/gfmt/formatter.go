package gfmt

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/lex8"
)

type formatter struct {
	*fmt8.Printer
	toks *tokens

	exprFunc func(f *formatter, expr ast.Expr)
}

func (f *formatter) printExpr(expr ast.Expr) {
	f.exprFunc(f, expr)
}

func (f *formatter) printExprs(args ...interface{}) {
	for _, arg := range args {
		f.printExpr(arg)
	}
}

func newFormatter(out io.Writer, toks []*lex8.Token) *formatter {
	p := fmt8.NewPrinter(out)
	return &formatter{
		Printer: p,
		toks:    newTokens(toks),
	}
}

func (f *formatter) printStr(s string) { fmt.Fprint(f.Printer, s) }
func (f *formatter) printSpace()       { f.printStr(" ") }
func (f *formatter) printEndl()        { fmt.Fprintln(f.Printer) }

func (f *formatter) peek() *lex8.Token {
	for {
		cur := f.toks.peek()
		if cur == nil {
			return nil
		}
		if cur.Type == parse.Semi {
			f.toks.shift()
			continue
		}
		return cur
	}
}

func (f *formatter) cue() *lex8.Token {
	for {
		cur := f.peek()
		if cur == nil {
			panic(fmt.Errorf("unmatched token: %v", cur))
		}

		if cur.Type == lex8.Comment {
			f.printStr(formatComment(cur.Lit))
			f.toks.shift()
			f.printEndlPlus(true, false)
			continue
		}

		return cur
	}
}

func (f *formatter) cueTo(token *lex8.Token) {
	cur := f.cue()
	if cur != token {
		panic(fmt.Errorf("unmatched token %v, got %v", token, cur))
	}
}

func (f *formatter) expect(token *lex8.Token) {
	f.cueTo(token)
	f.toks.shift()
}

// printEndlPlus prints one endline plus some number of empty lines if
// possible. This number is usually 0 or 1 depending on the line diffs between
// last token and next, but can be overriden by minEmptyLines.
func (f *formatter) printEndlPlus(plus, paraGap bool) {
	f.printEndl()
	if !plus {
		return
	}

	if paraGap {
		f.printEndl()
	}

	if f.toks.lineGap() >= 2 {
		f.printEndl()
	}
}

func (f *formatter) printToken(t *lex8.Token) {
	f.expect(t)
	f.printStr(t.Lit)
	f.printSameLineComments(t.Pos.Line)
}

func (f *formatter) printSameLineComments(line int) {
	for {
		tok := f.peek()
		if tok == nil {
			break
		}

		if !(tok.Type == lex8.Comment && tok.Pos.Line == line) {
			return
		}

		f.printSpace()
		f.printStr(formatComment(tok.Lit))
		f.toks.shift()
	}
}

func (f *formatter) finish() {
	tok := f.cue()
	if tok.Type != lex8.EOF {
		panic(fmt.Errorf("unmatched token: got %v, expected EOF", tok))
	}
	f.toks.shift()

	if f.toks.peek() != nil {
		panic(fmt.Errorf("unfinished tokens: %v", tok))
	}
}

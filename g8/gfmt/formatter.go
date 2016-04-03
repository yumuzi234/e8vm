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

func (f *formatter) printToken(token *lex8.Token) {
	f.cueToken(token)
	t := f.shift()
	if t != token {
		panic("bug")
	}

	f.printStr(t.Lit)
	f.printSameLineComments(t.Pos.Line)
}

func (f *formatter) cueToken(token *lex8.Token) {
	for {
		t := f.peek()
		if t == nil {
			panic(fmt.Errorf("unmatched token: %v", token))
		}

		if t.Type == lex8.Comment {
			f.printStr(formatComment(t.Lit))
			f.toks.shift()
			f.printEndlPlus(true, false)
			continue
		}

		if t != token {
			panic(fmt.Errorf(
				"unmatched token: got %v, expected %v", t, token,
			))
		}

		return
	}
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

func (f *formatter) expect(token *lex8.Token) {
	t := f.shift()
	if t == nil {
		panic(fmt.Errorf("unmatched token: %v", token))
	}
	if t != token {
		panic(fmt.Errorf("unmatched token: got %v, expected %v", t, token))
	}
}

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

func (f *formatter) shift() *lex8.Token {
	for {
		ret := f.toks.shift()
		if ret == nil {
			return nil
		}
		if ret.Type == parse.Semi {
			continue
		}

		if ret.Type == lex8.Comment {
			f.printStr(formatComment(ret.Lit))
			f.printEndlPlus(true, false)
			continue
		}

		return ret
	}
}

func (f *formatter) finish() {
	tok := f.shift()
	if tok.Type != lex8.EOF {
		panic(fmt.Errorf("unmatched token: got %v, expected EOF", tok))
	}

	tok = f.shift()
	if tok != nil {
		panic("unfinished tokens")
	}
}

package gfmt

import (
	"fmt"
	"io"

	"shanhu.io/smlvm/fmtutil"
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/parse"
)

type formatter struct {
	*fmtutil.Printer
	toks *tokens
	err  *lexing.Error
	last string // last string printed

	exprFunc func(f *formatter, expr interface{})
}

func (f *formatter) errs() []*lexing.Error {
	if f.err != nil {
		return []*lexing.Error{f.err}
	}
	return nil
}

func (f *formatter) errorf(pos *lexing.Pos, s string, args ...interface{}) {
	if f.err != nil {
		return
	}
	f.err = &lexing.Error{
		Pos: pos,
		Err: fmt.Errorf(s, args...),
	}
}

func (f *formatter) printExpr(expr interface{}) {
	f.exprFunc(f, expr)
}

func (f *formatter) printExprs(args ...interface{}) {
	for _, arg := range args {
		f.printExpr(arg)
	}
}

func newFormatter(out io.Writer, toks []*lexing.Token) *formatter {
	p := fmtutil.NewPrinter(out)
	return &formatter{
		Printer: p,
		toks:    newTokens(toks),
	}
}

func (f *formatter) printStr(s string) {
	if f.last == " " && s == " " {
		return // consume consecutive spaces
	}
	fmt.Fprint(f.Printer, s)
	f.last = s
}

func (f *formatter) printSpace() { f.printStr(" ") }

func (f *formatter) printEndl() {
	f.last = "\n"
	fmt.Fprintln(f.Printer)
}

func (f *formatter) peek() *lexing.Token {
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

func (f *formatter) cue() *lexing.Token {
	for {
		cur := f.peek()
		if cur == nil {
			return nil
		}

		if cur.Type == lexing.Comment {
			f.printStr(formatComment(cur.Lit))
			f.toks.shift()
			f.printGap()
			continue
		}

		return cur
	}
}

func (f *formatter) cueTo(token *lexing.Token) {
	cur := f.cue()
	if cur != token {
		f.errorf(token.Pos, "unmatched token %v, got %v", token, cur)
	}
}

func (f *formatter) expect(token *lexing.Token) {
	f.cueTo(token)
	f.toks.shift()
}

// printGap will print 2 endl when there is 2 or more endl originallly
// otherwise print 1 endl
func (f *formatter) printGap() {
	f.printEndl()
	if f.toks.lineGap() >= 2 {
		f.printEndl()
	}
}

func (f *formatter) printToken(t *lexing.Token) {
	f.expect(t)
	f.printStr(t.Lit)
	f.printSameLineComments(t.Pos.Line)
}

func (f *formatter) omitToken(t *lexing.Token) {
	f.expect(t)
	f.printSameLineComments(t.Pos.Line)
}

func (f *formatter) printSameLineComments(line int) {
	for {
		tok := f.peek()
		if tok == nil {
			break
		}

		if !(tok.Type == lexing.Comment && tok.Pos.Line == line) {
			return
		}

		f.printStr(" ")
		f.printStr(formatComment(tok.Lit))
		f.toks.shift()

		if f.toks.lineGap() == 0 {
			f.printStr(" ")
		}
	}
}

func (f *formatter) finish() {
	tok := f.cue()
	if tok.Type != lexing.EOF {
		f.errorf(tok.Pos, "unmatched token: got %v, expected EOF", tok)
		return
	}
	f.toks.shift()

	if f.toks.peek() != nil {
		f.errorf(tok.Pos, "unfinished tokens: %v", tok)
	}
}

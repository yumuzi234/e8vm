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
	tokens    []*lex8.Token
	lastToken *lex8.Token
	offset    int

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

func newFormatter(out io.Writer, tokens []*lex8.Token) *formatter {
	p := fmt8.NewPrinter(out)
	return &formatter{
		Printer: p,
		tokens:  tokens,
	}
}

func (f *formatter) printStr(s string) { fmt.Fprint(f.Printer, s) }
func (f *formatter) printSpace()       { f.printStr(" ") }
func (f *formatter) printEndl()        { fmt.Fprintln(f.Printer) }

// printEndlPlus prints one endline plus some number of empty lines if
// possible. This number is usually 0 or 1 depending on the line diffs between
// last token and next, but can be overriden by minEmptyLines.
func (f *formatter) printEndlPlus(plus bool, minEmptyLines int) {
	f.printEndl()
	if !plus {
		return
	}
	lines := 0
	if f.peek() != nil && f.lastToken != nil {
		lines = f.peek().Pos.Line - f.lastToken.Pos.Line - 1
	}
	if lines > 1 {
		lines = 1
	}
	if lines < minEmptyLines {
		lines = minEmptyLines
	}
	for i := 0; i < lines; i++ {
		f.printEndl()
	}
}

func (f *formatter) printToken(token *lex8.Token) {
	f.cueToken(token)
	f.printStr(token.Lit)
	f.printSameLineComments(token.Pos.Line)
}

func (f *formatter) cueToken(token *lex8.Token) {
	if f.tokens != nil {
		f.expect(token)
	}
}

func (f *formatter) printSameLineComments(line int) {
	for {
		tok := f.peek()
		if tok == nil {
			break
		}

		if tok.Type != lex8.Comment || tok.Pos.Line != line {
			return
		}

		f.printSpace()
		lit := tok.Lit
		if tok.Type == lex8.Comment {
			lit = formatComment(lit)
		}
		f.printStr(lit)
		f.offset++
	}
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

func (f *formatter) peek() *lex8.Token {
	for f.offset < len(f.tokens) {
		if token := f.tokens[f.offset]; token.Type != parse.Semi {
			return token
		}
		f.offset++ // ignore semi's
	}
	return nil
}

func (f *formatter) next() *lex8.Token {
	for f.offset < len(f.tokens) {
		f.lastToken = f.tokens[f.offset]
		f.offset++
		if f.lastToken.Type == parse.Semi {
			continue // ignore semi's
		}
		if f.lastToken.Type == lex8.Comment {
			lit := formatComment(f.lastToken.Lit)
			f.printStr(lit)
			f.printEndlPlus(true, 0)
			continue
		}
		return f.lastToken
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

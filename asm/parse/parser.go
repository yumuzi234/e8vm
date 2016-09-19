package parse

import (
	"io"

	"shanhu.io/smlvm/lexing"
)

// Parser parses a file input stream into top-level syntax blocks.
type parser struct {
	x lexing.Tokener
	*lexing.Parser
}

func newParser(f string, r io.Reader) (*parser, *lexing.Recorder) {
	ret := new(parser)

	var x lexing.Tokener = newLexer(f, r)
	x = newSemiInserter(x)
	rec := lexing.NewRecorder(x)
	ret.x = lexing.NewCommentRemover(rec)
	ret.Parser = lexing.NewParser(ret.x, Types)
	return ret, rec
}

func (p *parser) SeeKeyword(kw string) bool {
	return p.SeeLit(Keyword, kw)
}

func (p *parser) ExpectKeyword(kw string) *lexing.Token {
	return p.ExpectLit(Keyword, kw)
}

func (p *parser) skipErrStmt() bool {
	return p.SkipErrStmt(Semi)
}

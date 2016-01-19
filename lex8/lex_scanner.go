package lex8

import (
	"bytes"
	"io"
)

// lexScanner parses a file input stream into tokens.
type lexScanner struct {
	s     *runeScanner
	errs  *ErrorList
	valid bool

	pos *Pos
	buf *bytes.Buffer
}

// newLexScanner creates a new lexer.
func newLexScanner(file string, r io.Reader) *lexScanner {
	ret := &lexScanner{
		s:    newRuneScanner(file, r),
		errs: NewErrorList(),
		buf:  new(bytes.Buffer),
	}
	ret.pos = ret.s.pos()
	return ret
}

// next pushes the current rune (if valid) into the buffer,
// and returns the next rune or error from scanning the input
// stream.
func (s *lexScanner) next() (rune, error) {
	if s.valid {
		s.buf.WriteRune(s.s.Rune) // push into the buffer
		s.valid = false
	}

	if !s.s.scan() {
		if s.s.Err != nil {
			return 0, s.s.Err
		}

		return 0, io.EOF // signal end of file
	}

	s.valid = true
	return s.s.Rune, nil
}

// accept returns the string buffered, and the starting position
// of the string.
func (s *lexScanner) accept() (string, *Pos) {
	ret := s.buf.String()
	s.buf.Reset()
	pos := s.pos

	s.pos = s.s.pos()
	return ret, pos
}

// buffered returns the current buffered string in the
// scanner
func (s *lexScanner) buffered() string {
	return s.buf.String()
}

// startPos returns the position of the buffer start.
func (s *lexScanner) startPos() *Pos { return s.pos }

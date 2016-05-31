package lex8

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type errorScanner struct {
	scanner *runeScanner
	count   int
}

func newErrorScanner(filename string, r io.Reader) *errorScanner {
	return &errorScanner{
		scanner: newRuneScanner(filename, r),
		count:   1,
	}
}

var testErr = errors.New("timeout")

func (s *errorScanner) err() error {
	if s.count == 1 {
		for s.scanner.scan() {
		}
		s.count++
		return s.scanner.Err
	}
	return testErr
}

func TestRuneScanner(t *testing.T) {

	testCase := []struct {
		r    rune
		line int
		col  int
	}{
		{'a', 1, 1},
		{'~', 1, 2},
		{' ', 1, 3},
		{'\n', 1, 4},
		{'\n', 2, 1},
		{'1', 3, 1},
		{'A', 3, 2},
	}
	reader := strings.NewReader("a~ \n\n1A")
	fName := "testFileName"
	scanner := newRuneScanner(fName, reader)
	for _, tc := range testCase {
		if !scanner.scan() {
			t.Errorf("scan() FALSE, wants TRUE")
		}
		if scanner.pos().Col != tc.col {
			t.Errorf("col=%d, wants %d", scanner.pos().Col, tc.col)
		}
		if scanner.pos().Line != tc.line {
			t.Errorf("line=%d, wants %d", scanner.pos().Line, tc.line)
		}
		if scanner.pos().File != fName {
			t.Errorf("filename=%s, wants %s", scanner.pos().File, fName)
		}
		if scanner.Rune != tc.r {
			t.Errorf("rune=%c, wants %c", scanner.Rune, tc.r)
		}
	}
	if scanner.scan() {
		t.Errorf("scanner.scan()=%d, wants False", scanner.closed)
	}
	if !scanner.closed {
		t.Errorf("scanner close=%d, wants TRUE", scanner.closed)
	}
	if scanner.Err != nil {
		t.Errorf("Err=%v, wants nil", scanner.Err)
	}

	//test for errors
	reader = strings.NewReader("")
	ts := newErrorScanner(fName, reader)
	if ts.err() != nil {
		t.Errorf("Err=%v, wants nil", scanner.Err)
	}
	if ts.err() != testErr {
		t.Errorf("Err=%v, wants timeout", scanner.Err)
	}
}

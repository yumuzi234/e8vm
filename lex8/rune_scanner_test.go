package lex8

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

type errorScanner struct {
	s *runeScanner
	n int
}

func newErrorScanner(f string, r io.Reader) *errorScanner {
	return &errorScanner{
		s: newRuneScanner(f, r),
		n: 1,
	}
}

var errTest = errors.New("timeout")

func (s *errorScanner) err() error {
	if s.n == 1 {
		for s.s.scan() {
		}
		s.n++
		return s.s.Err
	}
	return errTest
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

	r := strings.NewReader("a~ \n\n1A")
	file := "a.txt"
	s := newRuneScanner(file, r)
	for _, tc := range testCase {
		if !s.scan() {
			t.Fatal("scan failed")
		}
		p := s.pos()
		want := &Pos{
			Col:  tc.col,
			Line: tc.line,
			File: file,
		}
		if !reflect.DeepEqual(p, want) {
			t.Errorf("pos got %v, want %v", p, want)
		}
		if s.Rune != tc.r {
			t.Errorf("rune got %c, want %c", s.Rune, tc.r)
		}
	}
	if s.scan() {
		t.Error("s.scan() got false, want true")
	}
	if !s.closed {
		t.Error("s close got false, want true")
	}
	if s.Err != nil {
		t.Errorf("expected error %v", s.Err)
	}

	// test for errors
	r = strings.NewReader("")
	ts := newErrorScanner(file, r)
	if ts.err() != nil {
		t.Errorf("Err=%v, wants nil", s.Err)
	}
	if ts.err() != errTest {
		t.Errorf("Err=%v, wants timeout", s.Err)
	}
}

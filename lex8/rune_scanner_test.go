package lex8

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

type errorReader struct {
	err error
	n   int
}

func newErrorReader(e error) *errorReader {
	return &errorReader{
		err: e,
		n:   1,
	}
}

var errTest = errors.New("timeout")

func (eR *errorReader) writeError(rS *runeScanner) {
	if eR.n == 1 {
		for rS.scan() {
		}
		eR.n++
	} else {
		rS.Err = eR.err
	}
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
	s = newRuneScanner(file, r)
	eR := newErrorReader(errTest)
	eR.writeError(s)
	if s.Err != nil {
		t.Errorf("Err=%v, wants nil", s.Err)
	}
	eR.writeError(s)
	if s.Err != errTest {
		t.Errorf("Err=%v, wants timeout", s.Err)
	}
}

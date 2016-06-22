package lex8

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

// LexComment lexes a c style comment.
func TestCommentLexer(t *testing.T) {
	e1 := fmt.Errorf("unexpected eof in block comment")
	e2 := fmt.Errorf("illegal char 'a'")
	var a rune
	a = 0
	e3 := fmt.Errorf("illegal char %q", a)
	testCase := []struct {
		Type int
		err  error
		e    error
		r    rune
	}{
		{Comment, e1, io.EOF, 0},
		{Comment, e1, io.EOF, 0},
		{Comment, nil, io.EOF, 0},
		{Comment, nil, nil, '\n'},
		{Comment, nil, nil, 'a'},
		{Comment, nil, io.EOF, 0},
		{Illegal, e2, nil, 'a'},
		{Illegal, e3, io.EOF, 0},
	}
	testString := []string{"/*//*a/", "/*\n\n", "//abc", "//abc\nabc", "/*abc*/abc", "/*abc*/", "/abc", "/"}
	i := 0
	for _, s := range testString {
		r := strings.NewReader(s)
		x := NewLexer("t1.txt", r)
		x.LexFunc = LexComment
		x.Next()
		tok := LexComment(x)
		tc := testCase[i]
		fmt.Printf("i=%d, r=%c\n", i, x.Rune())
		if tok.Type != tc.Type {
			t.Errorf("token want %q, got %q", tc.Type, tok.Type)
		}
		if tc.err != nil {
			if !reflect.DeepEqual(x.Errs()[0].Err, tc.err) {
				t.Errorf("want %q, got %q", tc.err, x.Errs()[0].Err)
			}
		} else if x.Errs() != nil {
			t.Errorf("unexpected error %q", x.Errs()[0].Err)
		}
		if tc.e != nil {
			if !reflect.DeepEqual(x.e, tc.e) {
				t.Errorf("want %q, got %q", tc.e, x.e)
			}
		} else if x.e != nil {
			t.Errorf("unexpected error %q", x.e)
		}

		if x.Rune() != tc.r {
			t.Errorf("want %c, got %c", tc.r, x.Rune())
		}
		i++
	}
}

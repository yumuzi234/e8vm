package lex8

import (
	"testing"
	"strings"
	"fmt"
)

// LexComment lexes a c style comment.
func TestCommitLexer(t *testing.T) {
	r:=strings.NewReader("/*//*abc/")
	x:=NewLexer("t1.txt",r)
	x.LexFunc=LexComment
	x.Next();
	fmt.Printf("x.Buffered()= %s \n", x.Buffered())
	fmt.Printf("x.Buffered()= %s \n", x.Buffered())
	tok:=x.Token()
	if (tok.Type!=Comment) {
		t.Errorf("token want Comment, got %q",tok.Type)
	} 
	if (x.errs.errs[0].Err!=fmt.Errorf("unexpected eof in block comment")) {
		t.Errorf("want unexpected eof in block comment, got %q",x.errs.errs[0].Err)
	}
}

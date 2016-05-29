package lex8

import (
	"strings"
	"testing"
)


func TestRuneScanner(t *testing.T) {
    
    testCase :=[] struct {
        r rune
        line int
        col int
    }{
        {'a', 1, 1},
        {'~', 1, 2},
        {' ', 1, 3},
        {'1', 3, 1},
        {'A', 3, 2},
    }
	reader := strings.NewReader("a~ /n/n1A")
    fName:="testFileName"
	scanner := newRuneScanner(fName, reader)
    for _, tc:=range testCase {
        if !scanner.scan() {
            t.Errorf("scan() FALSE, wants TRUE")
        }
        if scanner.pos().Col!=tc.col {
             t.Errorf("col=%d, wants %d", scanner.pos().Col, tc.col) 
        }
        if scanner.pos().Line!=tc.line {
             t.Errorf("line=%d, wants %d", scanner.pos().Line, tc.line) 
        }
        if scanner.pos().File!=fName {
             t.Errorf("filename=%s, wants %s", scanner.pos().File, fName) 
        }
        if scanner.Rune!=tc.r {
             t.Errorf("rune=%v, wants %v", scanner.Rune, tc.r) 
        }
    }
	
}
package lex8

import (
	"testing"
)

func TestIsLetter(t *testing.T) {
	for _, r := range "abzdATZ" {
		if !IsLetter(r) {
			t.Errorf("%v should be a letter", r)
		}
	}

	for _, r := range "013_%~-" {
		if IsLetter(r) {
			t.Errorf("%v should not be a letter", r)
		}
	}
}

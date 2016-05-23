
package lex8

import (
	"testing"
)

func TestIsLetter(t *testing.T) {
	t.Log ("Testing IsLetter")
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

func TestIsDigit(t *testing.T) {
	t.Log ("Testing IsDigit")
        for _, r := range "0123456789" {
                if !IsDigit(r) {
                        t.Errorf("%v should be a digit", r)
                }
        }

	for _, r := range "abzATZ#%~" {
		if IsDigit(r) {
                        t.Errorf("%v should not be a digit", r)
                }
	}
}


func TestIsHexDigit(t *testing.T) {
        t.Log ("Testing IsHexDigit")
        for _, r := range "0123456789" {
                if !IsHexDigit(r) {
                        t.Errorf("%v should be a hexdigit", r)
                }
        }

        for _, r := range "abcdefABCDEF"{
                if !IsHexDigit(r) {
                        t.Errorf("%v should be a hexdigit", r)
                }
        }

        for _, r := range "gJmXY!@*"{
                if IsHexDigit(r) {
                        t.Errorf("%v should not be a hexdigit", r)
                }
        }
}


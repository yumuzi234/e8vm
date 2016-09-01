package asm8

import (
	"e8vm.io/e8vm/asm8/parse"
	"e8vm.io/e8vm/lexing"
)

func isLabelStart(s string) bool {
	return len(s) > 0 && s[0] == '.'
}

func isLabel(s string) bool {
	if len(s) <= 1 || s[0] != '.' {
		return false
	}

	for i, r := range s[1:] {
		if r >= '0' && r <= '9' && i > 0 {
			continue
		}
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r == '_' {
			continue
		}
		return false
	}
	return true
}

func checkLabel(log lexing.Logger, t *lexing.Token) bool {
	if t.Type != parse.Operand {
		panic("not an operand")
	}

	lab := t.Lit
	if !isLabelStart(lab) {
		return false
	}

	if !isLabel(lab) {
		log.Errorf(t.Pos, "invalid label: %q", lab)
	}

	return true
}

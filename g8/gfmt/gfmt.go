package gfmt

import (
	"bytes"
	"strings"

	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/lex8"
)

// Format formats a file of a particular file name and content.
func Format(fname, content string) (string, []*lex8.Error) {
	r := strings.NewReader(content)
	ast, rec, errs := parse.File(fname, r, false)
	if errs != nil {
		return "", errs
	}

	out := new(bytes.Buffer)
	FprintFile(out, ast, rec)
	return out.String(), nil
}

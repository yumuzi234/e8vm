package gfmt

import (
	"bytes"
	"io"

	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/lexing"
)

func makeFormatter(out io.Writer, tokens []*lexing.Token) *formatter {
	ret := newFormatter(out, tokens)
	ret.exprFunc = printExpr
	return ret
}

// FileTo formats a g language file and output the formatted
// program via a writer.
func FileTo(fname string, file []byte, out io.Writer) []*lexing.Error {
	fast, rec, errs := parse.File(fname, bytes.NewBuffer(file), false)
	if errs != nil {
		return errs
	}

	f := makeFormatter(out, rec.Tokens())
	printFile(f, fast)
	return f.errs()
}

// File formats a g language file in bytes.
func File(fname string, file []byte) ([]byte, []*lexing.Error) {
	out := new(bytes.Buffer)
	errs := FileTo(fname, file, out)
	if errs != nil {
		return nil, errs
	}
	return out.Bytes(), nil
}

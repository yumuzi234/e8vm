package gfmt

import (
	"bytes"
	"io"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/lex8"
)

func makeFormatter(out io.Writer, tokens []*lex8.Token) *formatter {
	ret := newFormatter(out, tokens)
	ret.exprFunc = printExpr
	return ret
}

// FprintStmts prints the statements out to a writer
func FprintStmts(out io.Writer, stmts []ast.Stmt) {
	f := makeFormatter(out, nil) // TODO(h8liu): nil tokens, this will break
	printStmt(f, stmts)
}

// FprintFile prints a file
func FprintFile(out io.Writer, file *ast.File, rec *lex8.Recorder) {
	f := makeFormatter(out, rec.Tokens())
	printFile(f, file)
}

// File formats a g language file.
func File(fname string, file []byte) ([]byte, []*lex8.Error) {
	f, rec, errs := parse.File(fname, bytes.NewBuffer(file), false)
	if errs != nil {
		return nil, errs
	}

	out := new(bytes.Buffer)
	FprintFile(out, f, rec)
	return out.Bytes(), nil
}

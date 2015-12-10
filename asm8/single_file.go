package asm8

import (
	"bytes"
	"io"

	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/link8"
)

// BuildSingleFile builds a package named "main" from a single file.
func BuildSingleFile(f string, rc io.ReadCloser) ([]byte, []*lex8.Error) {
	path := "_"

	pinfo := build8.SimplePkg(path, f, rc)

	pkg, errs := Lang().Compile(pinfo)
	if errs != nil {
		return nil, errs
	}

	buf := new(bytes.Buffer)
	err := link8.LinkSingle(buf, pkg.Lib, pkg.Main)
	if err != nil {
		return nil, lex8.SingleErr(err)
	}

	return buf.Bytes(), nil
}

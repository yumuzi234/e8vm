package asm8

import (
	"bytes"
	"io"

	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/e8"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/link8"
)

// BuildSingleFile builds a package named "main" from a single file.
func BuildSingleFile(f string, rc io.ReadCloser) ([]byte, []*lex8.Error) {
	path := "_"

	pinfo := build8.SimplePkg(path, f, rc)

	pkg, errs := Lang().Compile(pinfo, new(build8.Options))
	if errs != nil {
		return nil, errs
	}

	secs, err := link8.LinkSinglePkg(pkg.Lib, pkg.Main)
	if err != nil {
		return nil, lex8.SingleErr(err)
	}

	buf := new(bytes.Buffer)
	if err := e8.Write(buf, secs); err != nil {
		return nil, lex8.SingleErr(err)
	}
	return buf.Bytes(), nil
}

package asm

import (
	"bytes"
	"io"

	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/image"
	"e8vm.io/e8vm/lexing"
	"e8vm.io/e8vm/link"
)

// BuildSingleFile builds a package named "main" from a single file.
func BuildSingleFile(f string, rc io.ReadCloser) ([]byte, []*lexing.Error) {
	path := "_"

	pinfo := build8.SimplePkg(path, f, rc)

	pkg, errs := Lang().Compile(pinfo)
	if errs != nil {
		return nil, errs
	}

	secs, err := link.SinglePkg(pkg.Lib, pkg.Main)
	if err != nil {
		return nil, lexing.SingleErr(err)
	}

	buf := new(bytes.Buffer)
	if err := image.Write(buf, secs); err != nil {
		return nil, lexing.SingleErr(err)
	}
	return buf.Bytes(), nil
}

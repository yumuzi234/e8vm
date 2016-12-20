package asm

import (
	"bytes"

	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/image"
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/link"
)

// BuildSingleFile builds a package named "main" from a single file.
func BuildSingleFile(f string, o builds.FileOpener) ([]byte, []*lexing.Error) {
	path := "_"

	pinfo := builds.SimplePkg(path, f, o)

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

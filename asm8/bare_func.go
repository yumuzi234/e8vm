package asm8

import (
	"io"

	"e8vm.io/e8vm/asm8/parse"
	"e8vm.io/e8vm/lexing"
	"e8vm.io/e8vm/link"
)

// BuildBareFunc builds a function body into an image.
func BuildBareFunc(f string, rc io.ReadCloser) ([]byte, []*lexing.Error) {
	fn, es := parse.BareFunc(f, rc)
	if es != nil {
		return nil, es
	}

	// resolving pass
	log := lexing.NewErrorList()
	rfunc := resolveFunc(log, fn)
	if es := log.Errs(); es != nil {
		return nil, es
	}

	// building pass
	b := newBuilder("main")
	fobj := buildFunc(b, rfunc)
	if es := b.Errs(); es != nil {
		return nil, es
	}

	ret, e := link.BareFunc(fobj)
	if e != nil {
		return nil, lexing.SingleErr(e)
	}

	return ret, nil
}

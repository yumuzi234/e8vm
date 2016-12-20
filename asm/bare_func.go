package asm

import (
	"shanhu.io/smlvm/asm/parse"
	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/link"
)

// BuildBareFunc builds a function body into an image.
func BuildBareFunc(f string, o builds.FileOpener) ([]byte, []*lexing.Error) {
	rc, err := o.Open()
	if err != nil {
		return nil, lexing.SingleErr(err)
	}
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

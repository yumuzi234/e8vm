package link8

import (
	"bytes"
	"fmt"
	"io"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/e8"
)

// Job is a linking job.
type Job struct {
	Pkgs   map[string]*Pkg
	Funcs  []*PkgSym
	InitPC uint32
}

// NewJob creates a new linking job which init pc is the default one.
func NewJob(pkgs map[string]*Pkg, funcs []*PkgSym) *Job {
	return &Job{
		Pkgs:   pkgs,
		Funcs:  funcs,
		InitPC: arch8.InitPC,
	}
}

// LinkMain is a short hand for NewJob(pkgs, path, start).Link(out)
func LinkMain(out io.Writer, pkgs map[string]*Pkg, funcs []*PkgSym) error {
	return NewJob(pkgs, funcs).Link(out)
}

// LinkSingle call LinkMain with only one single package.
func LinkSingle(out io.Writer, pkg *Pkg, start string) error {
	path := pkg.Path()
	pkgs := map[string]*Pkg{path: pkg}
	return LinkMain(out, pkgs, []*PkgSym{{path, start}})
}

// Link performs the linking job and writes the output to out.
func (j *Job) Link(out io.Writer) error {
	pkgs := j.Pkgs
	used := traceUsed(pkgs, j.Funcs)
	main := wrapMain(j.Funcs)
	funcs, vars, zeros, err := layout(pkgs, used, main, j.InitPC)
	if err != nil {
		return err
	}

	var secs []*e8.Section
	buf := new(bytes.Buffer)
	w := newWriter(pkgs, buf)
	w.writeFunc(main)
	for _, f := range funcs {
		w.writeFunc(pkgFunc(pkgs, f))
	}
	if err := w.Err(); err != nil {
		return err
	}

	if buf.Len() > 0 {
		secs = append(secs, &e8.Section{
			Header: &e8.Header{
				Type: e8.Code,
				Addr: j.InitPC,
			},
			Bytes: buf.Bytes(),
		})
	}

	if len(vars) > 0 {
		buf := new(bytes.Buffer)
		w := newWriter(pkgs, buf)
		for _, v := range vars {
			w.writeVar(pkgVar(pkgs, v))
		}
		if err := w.Err(); err != nil {
			return err
		}

		if buf.Len() > 0 {
			secs = append(secs, &e8.Section{
				Header: &e8.Header{
					Type: e8.Data,
					Addr: pkgVar(pkgs, vars[0]).addr,
				},
				Bytes: buf.Bytes(),
			})
		}
	}

	if len(zeros) > 0 {
		start := pkgVar(pkgs, zeros[0]).addr
		lastVar := pkgVar(pkgs, zeros[len(zeros)-1])
		end := lastVar.addr + lastVar.Size()
		secs = append(secs, &e8.Section{
			Header: &e8.Header{
				Type: e8.Zeros,
				Addr: start,
				Size: end - start,
			},
		})
	}

	return e8.Write(out, secs)
}

// LinkBareFunc produces a image of a single function that has no links.
func LinkBareFunc(f *Func) ([]byte, error) {
	if f.TooLarge() {
		return nil, fmt.Errorf("code section too large")
	}

	buf := new(bytes.Buffer)
	w := newWriter(make(map[string]*Pkg), buf)
	w.writeBareFunc(f)
	if err := w.Err(); err != nil {
		return nil, err
	}

	image := new(bytes.Buffer)
	sec := &e8.Section{
		Header: &e8.Header{
			Type: e8.Code,
			Addr: arch8.InitPC,
		},
		Bytes: buf.Bytes(),
	}
	if err := e8.Write(image, []*e8.Section{sec}); err != nil {
		return nil, err
	}

	return image.Bytes(), nil
}

package link

import (
	"bytes"
	"fmt"

	"e8vm.io/e8vm/arch"
	"e8vm.io/e8vm/image"
)

// Job is a linking job.
type Job struct {
	Pkgs   map[string]*Pkg
	Funcs  []*PkgSym
	InitPC uint32

	FuncDebug func(pkg, name string, addr, size uint32)
}

// NewJob creates a new linking job which init pc is the default one.
func NewJob(pkgs map[string]*Pkg, funcs []*PkgSym) *Job {
	return &Job{
		Pkgs:   pkgs,
		Funcs:  funcs,
		InitPC: arch.InitPC,
	}
}

// Main is a short hand for NewJob(pkgs, path, start).Link(out)
func Main(pkgs map[string]*Pkg, funcs []*PkgSym) ([]*image.Section, error) {
	return NewJob(pkgs, funcs).Link()
}

// SinglePkg call LinkMain with only one single package.
func SinglePkg(pkg *Pkg, start string) ([]*image.Section, error) {
	path := pkg.Path()
	pkgs := map[string]*Pkg{path: pkg}
	return Main(pkgs, []*PkgSym{{path, start}})
}

// Link performs the linking job and writes the output to out.
func (j *Job) Link() ([]*image.Section, error) {
	pkgs := j.Pkgs
	used := traceUsed(pkgs, j.Funcs)
	main := wrapMain(j.Funcs)
	funcs, vars, zeros, err := layout(pkgs, used, main, j.InitPC)
	if err != nil {
		return nil, err
	}

	var secs []*image.Section
	buf := new(bytes.Buffer)
	w := newWriter(pkgs, buf)
	w.writeFunc(main)
	for _, ps := range funcs {
		f := pkgFunc(pkgs, ps)
		if j.FuncDebug != nil {
			j.FuncDebug(ps.Pkg, ps.Sym, f.addr, f.Size())
		}
		w.writeFunc(f)
	}
	if err := w.Err(); err != nil {
		return nil, err
	}

	if buf.Len() > 0 {
		secs = append(secs, &image.Section{
			Header: &image.Header{
				Type: image.Code,
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
			return nil, err
		}

		if buf.Len() > 0 {
			secs = append(secs, &image.Section{
				Header: &image.Header{
					Type: image.Data,
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
		secs = append(secs, &image.Section{
			Header: &image.Header{
				Type: image.Zeros,
				Addr: start,
				Size: end - start,
			},
		})
	}

	return secs, nil
}

// BareFunc produces a image of a single function that has no links.
func BareFunc(f *Func) ([]byte, error) {
	if f.TooLarge() {
		return nil, fmt.Errorf("code section too large")
	}

	buf := new(bytes.Buffer)
	w := newWriter(make(map[string]*Pkg), buf)
	w.writeBareFunc(f)
	if err := w.Err(); err != nil {
		return nil, err
	}

	im := new(bytes.Buffer)
	sec := &image.Section{
		Header: &image.Header{
			Type: image.Code,
			Addr: arch.InitPC,
		},
		Bytes: buf.Bytes(),
	}
	if err := image.Write(im, []*image.Section{sec}); err != nil {
		return nil, err
	}

	return im.Bytes(), nil
}

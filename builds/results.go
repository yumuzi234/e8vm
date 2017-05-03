package builds

import (
	"io"
	"path"
)

type results struct {
	out Output
}

func newResults(out Output) *results {
	return &results{
		out: out,
	}
}

func (r *results) bin(p string) (io.WriteCloser, error) {
	rel := relPath(p)
	return r.out.Create(path.Join("bin", rel+".bin"))
}

func (r *results) testBin(p string) (io.WriteCloser, error) {
	rel := relPath(p)
	return r.out.Create(path.Join("test", rel+".bin"))
}

func (r *results) output(p, name string) (io.WriteCloser, error) {
	rel := relPath(p)
	return r.out.Create(path.Join("out", rel, name))
}

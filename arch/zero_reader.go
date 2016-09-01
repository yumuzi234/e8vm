package arch

import (
	"io"
)

type zeroReader struct {
	n uint32
}

func (r *zeroReader) Read(buf []byte) (int, error) {
	n := len(buf)
	if n == 0 {
		return 0, nil
	}
	if r.n == 0 {
		// nothing left
		return 0, io.EOF
	}

	if n > PageSize {
		// we read at most a page at a time
		n = PageSize
	}

	if uint32(n) < r.n {
		// have something left
		r.n -= uint32(n)
		copy(buf, make([]byte, n))
		return n, nil
	}

	copy(buf, make([]byte, r.n))
	r.n = 0
	return int(r.n), nil
}

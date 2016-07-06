package build8

import (
	"testing"
)

func TestIsPkgPath(t *testing.T) {
	o := func(p string) {
		if !isPkgPath(p) {
			t.Errorf("%q should be a valid package path, but not", p)
		}
	}

	o("asm/builtin")
	o("hello")
	o("/something")
	o("/h8liu/std")

	e := func(p string) {
		if isPkgPath(p) {
			t.Errorf("%q should be an invalid package path, but not", p)
		}
	}

	e("")
	e("/")
	e("/h8liu/")
	e("//")
	e("h8liu/")
	e("   ")
	e("3435")
	e("x/3435")
	e("3435/x")
}

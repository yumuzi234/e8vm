package builds

import (
	"testing"
)

func TestIsPkgPath(t *testing.T) {
	o := func(p string) {
		if !IsPkgPath(p) {
			t.Errorf("%q should be a valid package path, but not", p)
		}
	}

	o("/std/asm/builtin")
	o("something/nothing")
	o("hello")
	o("/something")
	o("/h8liu/std")

	e := func(p string) {
		if IsPkgPath(p) {
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

func TestIsParentPkg(t *testing.T) {
	o := func(p, sub string) {
		if !IsParentPkg(p, sub) {
			t.Errorf("%q should be a parent package of %q, but not", p, sub)
		}
	}

	o("", "something")
	o("", "something/something")
	o("", "/x")
	o("x", "x")
	o("x", "x/y")
	o("/", "/x")
	o("/", "/")
	o("/xxx", "/xxx/yyy")

	e := func(p, sub string) {
		if IsParentPkg(p, sub) {
			t.Errorf("%q should not be a sub package of %q, but is", p, sub)
		}
	}

	e("/", "")
	e("/", "something")
	e("/x", "x")
	e("x", "xxx/yyy")
	e("/x", "/xxx/yyy")
}

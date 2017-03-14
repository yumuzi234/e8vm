package builds

import (
	"reflect"
	"testing"
)

var _ Input2 = new(MemFS)
var _ Output2 = new(MemFS)

func TestMemFS(t *testing.T) {
	ne := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	as := func(cond bool) {
		if !cond {
			t.Fatal("assertion failed")
		}
	}

	fs := NewMemFS()
	ne(fs.AddTextFile("std/asm/builtin/builtin.s", "something"))
	ok, err := fs.HasDir("std/asm/builtin")
	ne(err)
	as(ok)

	dirs, err := fs.ListDirs("")
	ne(err)
	as(reflect.DeepEqual(dirs, []string{"std"}))
}

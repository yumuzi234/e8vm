package ir

type testList struct {
	pkg, sym string

	name  string
	funcs []*Func
}

func newTestList(name string, funcs []*Func) *testList {
	return &testList{name: name, funcs: funcs}
}

func (lst *testList) String() string     { return lst.name }
func (lst *testList) RegSizeAlign() bool { return true }
func (lst *testList) Size() int32 {
	return regSize * int32(len(lst.funcs))
}

package ir

type testList struct {
	pkg, name string
	funcs     []*Func
}

func newTestList(pkg, name string, funcs []*Func) *testList {
	return &testList{pkg: pkg, name: name, funcs: funcs}
}

func (lst *testList) String() string     { return lst.name }
func (lst *testList) RegSizeAlign() bool { return true }
func (lst *testList) Size() int32 {
	return regSize * int32(len(lst.funcs))
}

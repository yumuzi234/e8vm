package sempass

import (
	"fmt"

	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/types"
	"shanhu.io/smlvm/syms"
)

type canAssignResult struct {
	err      bool
	needCast bool
	castMask []bool
}

func (r *canAssignResult) add(ok, needCast bool) {
	r.err = r.err || !ok
	r.needCast = r.needCast || needCast
	r.castMask = append(r.castMask, needCast)
}

func canAssignType(
	b *builder, pos *lexing.Pos, t types.T, right []types.T,
	in string,
) *canAssignResult {
	var ts []types.T
	for range right {
		ts = append(ts, t)
	}
	return canAssigns(b, pos, ts, right, in)
}

func canAssigns(
	b *builder, pos *lexing.Pos, left, right []types.T, in string,
) *canAssignResult {
	if len(left) != len(right) {
		panic("length mismatch")
	}

	res := new(canAssignResult)
	for i, t := range right {
		res.add(canAssign(b, pos, left[i], t, in))
	}
	return res
}

func canAssign(
	b *builder, p *lexing.Pos, left, right types.T, in string,
) (ok bool, needCast bool) {
	if i, ok := left.(*types.Interface); ok {
		if !assignInterface(b, p, i, right, in) {
			return false, false
		}
		return true, true
	}
	ok, needCast = types.CanAssign(left, right)
	if !ok {
		b.CodeErrorf(p, "pl.cannotAssign.typeMismatch",
			"cannot use %s as %s in %s", left, right, in)
		return false, false
	}
	return ok, needCast
}

func assignInterface(
	b *builder, p *lexing.Pos, i *types.Interface, right types.T,
	in string,
) bool {
	flag := true
	var syms *syms.Table
	if t, ok := types.PointerOf(right).(*types.Struct); ok {
		syms = t.Syms
	} else if t, ok := right.(*types.Interface); ok {
		syms = t.Syms
	} else {
		b.CodeErrorf(p, "pl.cannotAssign.interface",
			"cannot use %s as interface %s in %s, "+
				"not a struct pointer or interface",
			right, i, in)
		return false
	}

	errorf := func(f string, a ...interface{}) {
		m := fmt.Sprintf(f, a...)
		b.CodeErrorf(p, "pl.cannotAssign.interface",
			"cannot use %s as interface %s in %s, %s", right, i, in, m)
		flag = false
	}

	funcs := i.Syms.List()
	for _, f := range funcs {
		sym := syms.Query(f.Name())
		if sym == nil {
			errorf("function %s not implemented in %s", f.Name(), right)
			continue
		}
		t2, ok := sym.ObjType.(*types.Func)
		if !ok {
			errorf("%s is not a function in %s", sym.Name(), right)
			continue
		}
		if !t2.IsBond {
			t2 = t2.MethodFunc
		}
		t1 := f.ObjType.(*types.Func)
		if !types.SameType(t1, t2) {
			errorf("function mismatch, want %q, have %q", t1, t2)
		}
	}
	return flag
}

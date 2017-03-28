package types

import "fmt"

// CanAssign checks if right can be assigned to left
func CanAssign(left, right T) bool {
	if c, ok := right.(*Const); ok {
		if _, ok := c.Type.(Number); ok {
			return InRange(c.Value.(int64), left)
		}
		right = c.Type
	}

	if _, ok := left.(*Type); ok {
		return false
	}
	if _, ok := left.(*Pkg); ok {
		return false
	}

	if IsNil(right) {
		switch left := left.(type) {
		case *Pointer:
			return true
		case *Slice:
			return true
		case *Func:
			if left.IsBond {
				return false
			}
			return true
		case *Interface:
			return true
		}
		return false
	}

	return SameType(left, right)
}

// SameType checks if two types are of the same type
func SameType(t1, t2 T) bool {
	if t1 == t2 {
		return true
	}
	switch t1 := t1.(type) {
	case null:
		return false
	case *Const:
		return false
	case Basic:
		if t2, ok := t2.(Basic); ok {
			return t2 == t1
		}
		return false
	case *Pointer:
		if t2, ok := t2.(*Pointer); ok {
			return SameType(t1.T, t2.T)
		}
		return false
	case *Slice:
		if t2, ok := t2.(*Slice); ok {
			return SameType(t1.T, t2.T)
		}
		return false
	case *Array:
		if t2, ok := t2.(*Array); ok {
			return t1.N == t2.N && SameType(t1.T, t2.T)
		}
		return false
	case *Func:
		t2, ok := t2.(*Func)
		if !ok {
			return false
		}
		if t2.IsBond {
			return false
		}
		if len(t1.Args) != len(t2.Args) {
			return false
		}
		if len(t1.Rets) != len(t2.Rets) {
			return false
		}

		for i, t := range t1.Args {
			if !SameType(t.T, t2.Args[i].T) {
				return false
			}
		}

		for i, t := range t1.Rets {
			if !SameType(t.T, t2.Rets[i].T) {
				return false
			}
		}

		return true
	case *Struct:
		if t2, ok := t2.(*Struct); ok {
			return t1 == t2
		}
		return false
	case *Interface:
		if t2, ok := t2.(*Interface); ok {
			return t1 == t2
		}
		return false
	default:
		panic(fmt.Errorf("invalid type: %T", t1))
	}
}

// BothPointer checks if the internal type are the same pointer types.
// If both are of the same pointer type, it returns true.
// If one is nil, but the other one is a pointer, it returns true.
// Otherwise it returns false.
func BothPointer(t1, t2 T) bool {
	p1 := PointerOf(t1)
	p2 := PointerOf(t2)
	if IsNil(t1) && p2 != nil {
		return true
	} else if IsNil(t2) && p1 != nil {
		return true
	} else if p1 == nil || p2 == nil {
		return false
	}
	return SameType(p1, p2)
}

// BothFuncPointer checks if the two types are comparable func pointers.
// If they are the same type, it returns true.
// If one is nil, but the other one is a func pointer, it returns true.
// Otherwise it returns false.
func BothFuncPointer(t1, t2 T) bool {
	b1 := IsFuncPointer(t1)
	b2 := IsFuncPointer(t2)
	if IsNil(t1) && b2 {
		return true
	} else if IsNil(t2) && b1 {
		return true
	} else if !b1 || !b2 {
		return false
	}

	return SameType(t1, t2)
}

// BothSlice checks if the internal are of the same slice types
// If one of them is nil, but the other one is not, it returns the element
// type of the slice type.
// If one of t1 and t2 is not a slice and it is not a nil, it returns nil.
// If t1 and t2 are of different slice types, it returns nil.
func BothSlice(t1, t2 T) bool {
	p1 := SliceOf(t1)
	p2 := SliceOf(t2)
	if IsNil(t1) && p2 != nil {
		return true
	} else if IsNil(t2) && p1 != nil {
		return true
	} else if p1 == nil || p2 == nil {
		return false
	}

	return SameType(p1, p2)
}

// CastConst checks if a const can be used to define a const with type T.
// and return the Const type
func CastConst(ct *Const, t T) *Const {
	if SameType(t, ct.Type) {
		return &Const{Value: ct.Value, Type: t}
	}
	_, ok := ct.Type.(Number)
	if IsInteger(t) && ok {
		ret, e := NewConstInt(ct.Value.(int64), t)
		if e == nil {
			return ret
		}
	}
	return nil
}

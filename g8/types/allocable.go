package types

// IsAllocable tells if a type can be allocated in memory.
func IsAllocable(t T) bool {
	switch t := t.(type) {
	case *Pkg:
		return false
	case *Type:
		return false
	case *BuiltInFunc:
		return false
	case *Func:
		return !t.IsBond
	}

	return true
}

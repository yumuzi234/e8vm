package g8

import (
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/link8"
	"e8vm.io/e8vm/sym8"
)

var (
	refTrue  = newRef(types.Bool, ir.Byt(1))
	refFalse = newRef(types.Bool, ir.Byt(0))
	refNil   = newRef(types.Nil(), nil)
)

func declareBuiltin(b *builder, builtin *link8.Pkg) {
	path := builtin.Path()
	e := b.p.HookBuiltin(builtin)
	if e != nil {
		b.Errorf(nil, e.Error())
		return
	}

	o := func(name, as string, t *types.Func) ir.Ref {
		sym := builtin.SymbolByName(name)
		if sym == nil {
			b.Errorf(nil, "builtin symbol %s missing", name)
			return nil
		} else if sym.Type != link8.SymFunc {
			b.Errorf(nil, "builtin symbol %s is not a function", name)
			return nil
		}

		ref := ir.NewFuncSym(path, name, nil)
		obj := &objFunc{as, newRef(t, ref), nil, false}
		pre := b.scope.Declare(sym8.Make(b.symPkg, as, symFunc, obj, nil))
		if pre != nil {
			b.Errorf(nil, "builtin symbol %s declare failed", name)
			return nil
		}
		return ref
	}

	// TODO: should export these as generic function pointer symbols, and
	// convert them into G functions in g files, rahter than inside here in
	// the compiler.
	o("PrintInt32", "printInt", types.NewVoidFunc(types.Int))
	o("PrintUint32", "printUint", types.NewVoidFunc(types.Uint))
	o("PrintChar", "printChar", types.NewVoidFunc(types.Int8))
	o("Vtable", "vtable", types.NewVoidFunc(types.Uint))
	b.panicFunc = o("Panic", "panic", types.VoidFunc)
	o("Halt", "halt", types.VoidFunc)

	o("Syscall", "syscall", types.NewFuncUnamed(
		[]types.T{types.Uint, types.Uint, types.Uint},
		[]types.T{types.Uint, types.Uint, types.Uint},
	))

	// TODO: these are just hacks for context switch
	oe := func(name, as string, t *types.Func) {
		sym := builtin.SymbolByName(name)
		if sym == nil {
			return
		}

		ref := ir.NewFuncSym(path, name, nil)
		obj := &objFunc{as, newRef(t, ref), nil, false}
		pre := b.scope.Declare(sym8.Make(b.symPkg, as, symFunc, obj, nil))
		if pre != nil {
			b.Errorf(nil, "builtin symbol %s declare failed", name)
			return
		}
	}
	oe("Ienter", "_ienter", types.NewVoidFunc())
	oe("SysEnter", "_sysenter", types.NewVoidFunc())
	oe("Ustart", "ustart", types.NewVoidFunc(types.Uint, types.Uint))

	ov := func(name, as string, t types.T) {
		sym := builtin.SymbolByName(name)
		if sym == nil {
			return
		}
		ref := ir.NewHeapSym(path, name, t.Size(), false, true)
		obj := &objVar{name: name, ref: newAddressableRef(t, ref)}
		pre := b.scope.Declare(sym8.Make(b.symPkg, as, symVar, obj, nil))
		if pre != nil {
			b.Errorf(nil, "builtin symbol %s declare failed", name)
			return
		}
	}
	ov("Ientry", "_ientry", types.Uint)
	ov("SysEntry", "_sysentry", types.Uint)

	bi := func(name string) {
		obj := &objFunc{name, newRef(types.NewBuiltInFunc(name), nil),
			nil, false,
		}
		s := sym8.Make(b.symPkg, name, symFunc, obj, nil)
		pre := b.scope.Declare(s)
		if pre != nil {
			b.Errorf(nil, "builtin symbol %s declare failed", name)
		}
	}

	bi("len")
	bi("make")

	c := func(name string, r *ref) {
		// TODO: declare these as typed consts
		obj := &objConst{name, r}
		s := sym8.Make(b.symPkg, name, symConst, obj, nil)
		pre := b.scope.Declare(s)
		if pre != nil {
			b.Errorf(nil, "builtin symbol %s declare failed", name)
		}
	}

	c("true", refTrue)
	c("false", refFalse)
	c("nil", refNil)

	t := func(name string, t types.T) {
		obj := &objType{name, newTypeRef(t)}
		s := sym8.Make(b.symPkg, name, symType, obj, nil)
		pre := b.scope.Declare(s)
		if pre != nil {
			b.Errorf(nil, "builtin symbol %s declare failed", name)
		}
	}

	t("int", types.Int)
	t("uint", types.Uint)
	t("int32", types.Int)
	t("uint32", types.Uint)
	t("int8", types.Int8)
	t("uint8", types.Uint8)
	t("char", types.Int8)
	t("byte", types.Uint8)
	t("bool", types.Bool)
	t("float", types.Float32)
	t("string", types.String)
	// t("ptr", &types.Pointer{types.Uint8})
	// t("float32", types.Float32)
	// t("string", types.String)
}

func isBasicType(t string) bool {
	switch t {
	case "int", "uint", "int32", "uint32",
		"int8", "uint8", "char", "byte",
		"bool", "float", "string":
		return true
	}
	return false
}

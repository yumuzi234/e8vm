package glang

import (
	"e8vm.io/e8vm/glang/codegen"
	"e8vm.io/e8vm/glang/tast"
	"e8vm.io/e8vm/glang/types"
	"e8vm.io/e8vm/link8"
	"e8vm.io/e8vm/sym8"
)

var (
	refTrue  = newRef(types.Bool, codegen.Byt(1, true))
	refFalse = newRef(types.Bool, codegen.Byt(0, true))
	refNil   = newRef(types.Nil(), nil)
)

func declareBuiltin(b *builder, builtin *link8.Pkg) {
	path := builtin.Path()
	e := b.p.HookBuiltin(builtin)
	if e != nil {
		b.Errorf(nil, e.Error())
		return
	}

	o := func(name, as string, t *types.Func) codegen.Ref {
		sym := builtin.SymbolByName(name)
		if sym == nil {
			b.Errorf(nil, "builtin symbol %s missing", name)
			return nil
		} else if sym.Type != link8.SymFunc {
			b.Errorf(nil, "builtin symbol %s is not a function", name)
			return nil
		}

		ref := codegen.NewFuncSym(path, name, makeFuncSig(t))
		obj := &objFunc{name: as, ref: newRef(t, ref)}
		s := sym8.Make(b.path, as, tast.SymFunc, obj, t, nil)
		pre := b.scope.Declare(s)
		if pre != nil {
			b.Errorf(nil, "builtin symbol %s declare failed", name)
			return nil
		}
		return ref
	}

	// TODO: export these as generic function pointer symbols, and convert
	// them into G functions in g files, rather than inside here in the
	// compiler.
	o("PrintInt32", "printInt", types.NewVoidFunc(types.Int))
	o("PrintUint32", "printUint", types.NewVoidFunc(types.Uint))
	o("PrintChar", "printChar", types.NewVoidFunc(types.Int8))
	b.panicFunc = o("Panic", "panic", types.VoidFunc)
	o("Sleep", "sleep", types.VoidFunc)
	o("Assert", "assert", types.NewVoidFunc(types.Bool))

	bi := func(name string) {
		t := types.NewBuiltInFunc(name)
		obj := &objFunc{name: name, ref: newRef(t, nil)}
		s := sym8.Make(b.path, name, tast.SymFunc, obj, t, nil)
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
		s := sym8.Make(b.path, name, tast.SymConst, obj, r.Type(), nil)
		pre := b.scope.Declare(s)
		if pre != nil {
			b.Errorf(nil, "builtin symbol %s declare failed", name)
		}
	}

	c("true", refTrue)
	c("false", refFalse)
	c("nil", refNil)

	t := func(name string, t types.T) {
		s := sym8.Make(b.path, name, tast.SymType, nil, &types.Type{t}, nil)
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
	t("uintptr", types.Uint)
	// t("float32", types.Float32)
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

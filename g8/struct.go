package g8

import (
	"fmt"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/sym8"
)

func declareStruct(b *builder, d *ast.Struct) *structInfo {
	info := newStructInfo(d)
	name := info.Name()
	pos := info.name.Pos
	obj := &objType{name, newTypeRef(info.t)}
	pre := b.scope.Declare(sym8.Make(b.symPkg, name, symStruct, obj, pos))
	if pre != nil {
		b.Errorf(pos, "%s already declared", name)
		b.Errorf(pre.Pos, "previously declared here as a %s",
			symStr(pre.Type),
		)
		return nil
	}

	return info
}

func defineStructFields(b *builder, info *structInfo) {
	s := info.ast
	t := info.t

	// build fields
	for _, f := range s.Fields {
		ft := b.buildType(f.Type)
		if ft == nil {
			return
		}

		for _, id := range f.Idents.Idents {
			name := id.Lit
			field := new(types.Field)
			field.Name = name
			field.T = ft

			obj := &objField{name, field}
			sym := sym8.Make(b.symPkg, name, symField, obj, id.Pos)
			pre := t.Syms.Declare(sym)
			if pre != nil {
				b.Errorf(id.Pos, "field %s already defined", id.Lit)
				b.Errorf(pre.Pos, "previously defined here")
				continue
			}

			t.AddField(field)
		}
	}
}

func declareStructMethods(b *builder, info *structInfo) {
	for _, m := range info.ast.Methods {
		obj := declareMethod(b, info, m)
		if obj != nil {
			info.methodObjs = append(info.methodObjs, obj)
		}
	}
}

func declareMethod(b *builder, info *structInfo, f *ast.Func) *objFunc {
	t := buildFuncType(b, info, f.FuncSig)
	if t == nil {
		return nil
	}

	name := f.Name.Lit
	ret := &objFunc{name, nil, f, true}

	sym := sym8.Make(b.symPkg, name, symFunc, ret, f.Name.Pos)
	pre := info.t.Syms.Declare(sym)
	if pre != nil {
		b.Errorf(f.Name.Pos, "member %s already defined", name)
		b.Errorf(pre.Pos, "previously defined here")
		return nil
	}

	fullName := fmt.Sprintf("%s:%s", info.name.Lit, name)
	irFunc := b.p.NewFunc(fullName, f.Name.Pos, makeFuncSig(t))
	ret.ref = newRef(t, irFunc)

	return ret
}

func buildMethods(b *builder, info *structInfo) {
	if !b.golike {
		b.scope.PushTable(info.t.Syms)
		defer b.scope.Pop()
	}

	for _, m := range info.methodObjs {
		buildMethodFunc(b, info, m)
	}
}

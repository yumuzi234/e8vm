package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

func declareConst(b *builder, tok *lex8.Token, t types.T) *objConst {
	name := tok.Lit
	v := &objConst{name: name}
	s := sym8.Make(b.path, name, tast.SymConst, v, t, tok.Pos)
	conflict := b.scope.Declare(s)
	if conflict != nil {
		b.Errorf(tok.Pos, "%q already declared as a %s",
			name, tast.SymStr(conflict.Type),
		)
		return nil
	}
	return v
}

func buildGlobalConstDecl(b *builder, info *constInfo) {
	if info.typ != nil {
		b.Errorf(ast.ExprPos(info.typ), "typed const not implemented yet")
		return
	}

	right := buildConstExpr(b, info.expr)
	if right == nil {
		return
	}

	if !right.IsSingle() {
		b.Errorf(ast.ExprPos(info.expr), "must be single expression")
		return
	}

	t := right.Type()
	if !types.IsConst(t) {
		b.Errorf(ast.ExprPos(info.expr), "not a const")
		return
	}

	obj := declareConst(b, info.name, t)
	if obj == nil {
		return
	}

	obj.ref = right
}

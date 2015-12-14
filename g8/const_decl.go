package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

func buildConstExprList(b *builder, list *ast.ExprList) *ref {
	n := list.Len()

	ret := new(ref)
	if n == 0 {
		return ret
	}
	if n == 1 {
		return b.buildConstExpr(list.Exprs[0])
	}

	for _, expr := range list.Exprs {
		ref := b.buildConstExpr(expr)
		if ref == nil {
			return nil
		}
		if !ref.IsSingle() {
			b.Errorf(ast.ExprPos(expr), "cannot composite list in a list")
			return nil
		}
		ref.addressable = false
		ret = appendRef(ret, ref)
	}

	return ret
}

func declareConst(b *builder, tok *lex8.Token) *objConst {
	name := tok.Lit
	v := &objConst{name: name}
	s := sym8.Make(b.path, name, symConst, v, tok.Pos)
	conflict := b.scope.Declare(s)
	if conflict != nil {
		b.Errorf(tok.Pos, "%q already declared as a %s",
			name, symStr(conflict.Type),
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

	obj := declareConst(b, info.name)
	if obj == nil {
		return
	}

	obj.ref = right
}

func buildConstDecl(b *builder, d *ast.ConstDecl) {
	if d.Type != nil {
		b.Errorf(ast.ExprPos(d.Type), "typed const not implemented yet")
		return
	}

	right := buildConstExprList(b, d.Exprs)
	if right == nil {
		return
	}

	nright := right.Len()
	idents := d.Idents.Idents
	nleft := len(idents)
	if nleft != nright {
		b.Errorf(d.Eq.Pos, "%d values for %d identifiers",
			nright, nleft,
		)
		return
	}

	for i, ident := range idents {
		rightRef := right.At(i)
		t := rightRef.Type()
		if !types.IsConst(t) {
			b.Errorf(ast.ExprPos(d.Exprs.Exprs[i]), "not a const")
			return
		}

		obj := declareConst(b, ident)
		if obj == nil {
			return
		}
		obj.ref = newRef(t, rightRef.IR())
	}
}

func buildConstDecls(b *builder, decls *ast.ConstDecls) {
	for _, d := range decls.Decls {
		buildConstDecl(b, d)
	}
}

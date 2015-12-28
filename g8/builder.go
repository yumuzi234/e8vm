package g8

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/sempass"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

// builder builds a package
type builder struct {
	*lex8.ErrorList
	path string

	p         *ir.Pkg
	f         *ir.Func
	fretNamed bool
	fretRef   *ref

	golike bool

	b     *ir.Block
	scope *sym8.Scope

	continues *blockStack
	breaks    *blockStack

	exprFunc  func(b *builder, expr tast.Expr) *ref
	stmtFunc  func(b *builder, stmt ast.Stmt)
	stmtFunc2 func(b *builder, stmt tast.Stmt)
	irLog     io.WriteCloser

	panicFunc ir.Ref

	// this pointer, only valid when building a method.
	this *ref

	anonyCount int

	rand *rand.Rand

	spass *sempass.Builder
}

func newRand() *rand.Rand {
	var buf [8]byte
	_, err := crand.Read(buf[:])
	if err != nil {
		panic(err)
	}
	seed := int64(binary.LittleEndian.Uint64(buf[:]))
	return rand.New(rand.NewSource(seed))
}

func newBuilder(path string, golike bool) *builder {
	s := sym8.NewScope()
	return &builder{
		ErrorList: lex8.NewErrorList(),
		path:      path,
		p:         ir.NewPkg(path),
		scope:     s, // package scope
		golike:    golike,

		continues: newBlockStack(),
		breaks:    newBlockStack(),

		rand: newRand(),

		spass: sempass.NewBuilder(path, s),
	}
}

func (b *builder) anonyName(name string) string {
	if name == "_" {
		name = fmt.Sprintf("_:%d", b.anonyCount)
		b.anonyCount++
	}
	return name
}

func (b *builder) newTempIR(t types.T) ir.Ref {
	return b.f.NewTemp(t.Size(), types.IsByte(t), t.RegSizeAlign())
}

func (b *builder) newTemp(t types.T) *ref { return newRef(t, b.newTempIR(t)) }

func (b *builder) newCond() ir.Ref { return b.f.NewTemp(1, true, false) }
func (b *builder) newPtr() ir.Ref  { return b.f.NewTemp(4, true, true) }

func (b *builder) newAddressableTemp(t types.T) *ref {
	return newAddressableRef(t, b.newTempIR(t))
}

func (b *builder) newLocal(t types.T, name string) ir.Ref {
	return b.f.NewLocal(t.Size(), name,
		types.IsByte(t), t.RegSizeAlign(),
	)
}

func (b *builder) newGlobalVar(t types.T, name string) ir.Ref {
	name = b.anonyName(name)
	return b.p.NewGlobalVar(t.Size(), name, types.IsByte(t), t.RegSizeAlign())
}

func (b *builder) buildExpr(expr tast.Expr) *ref {
	if b.exprFunc == nil {
		panic("exprFunc2 function not assigned")
	}
	return b.exprFunc(b, expr)
}

func (b *builder) buildStmts(stmts []ast.Stmt) {
	if b.stmtFunc == nil {
		return
	}

	for _, stmt := range stmts {
		b.stmtFunc(b, stmt)
	}
}

func (b *builder) buildStmt(stmt ast.Stmt) { b.stmtFunc(b, stmt) }

func (b *builder) buildStmt2(stmt tast.Stmt) { b.stmtFunc2(b, stmt) }

func (b *builder) Errs() []*lex8.Error {
	if errs := b.spass.Errs(); errs != nil {
		return errs
	}
	return b.ErrorList.Errs()
}

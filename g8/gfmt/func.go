package gfmt

import (
	"e8vm.io/e8vm/g8/ast"
)

func printFunc(f *formatter, fn *ast.Func) {
	f.printExprs(fn.Kw, " ")
	if r := fn.Recv; r != nil {
		f.printExprs(
			r.Lparen, r.Recv, " ", r.Star, r.StructName, r.Rparen, " ",
		)
	}
	f.printExprs(fn.Name)
	printFuncSig(f, fn.FuncSig)
	f.printSpace()
	printStmt(f, fn.Body)
}

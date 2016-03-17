package gfmt

import (
	"e8vm.io/e8vm/g8/ast"
)

func printFunc(f *formatter, fn *ast.Func) {
	f.printExprs(fn.Kw, " ", fn.Name)
	printFuncSig(f, fn.FuncSig)
	f.printSpace()
	printStmt(f, fn.Body)
}

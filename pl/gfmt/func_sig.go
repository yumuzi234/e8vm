package gfmt

import (
	"shanhu.io/smlvm/pl/ast"
)

func printParaList(f *formatter, lst *ast.ParaList) {
	f.printToken(lst.Lparen)
	for i, para := range lst.Paras {
		if para.Ident != nil {
			f.printToken(para.Ident)
			if para.Type != nil {
				f.printStr(" ")
			}
		}

		if para.Type != nil {
			f.printExpr(para.Type)
		}
		if i < len(lst.Commas) {
			f.printExprs(lst.Commas[i], " ")
		}
	}
	f.printToken(lst.Rparen)
}

func printFuncSig(f *formatter, fsig *ast.FuncSig) {
	printParaList(f, fsig.Args)
	if fsig.RetType != nil {
		f.printExprs(" ", fsig.RetType)
	} else if fsig.Rets != nil {
		f.printStr(" ")
		printParaList(f, fsig.Rets)
	}
}

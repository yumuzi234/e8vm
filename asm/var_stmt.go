package asm

import (
	"e8vm.io/e8vm/asm/ast"
	"e8vm.io/e8vm/lexing"
)

type varStmt struct {
	*ast.VarStmt

	align uint32
	data  []byte
}

func resolveVarStmt(log lexing.Logger, v *ast.VarStmt) *varStmt {
	ret := new(varStmt)
	ret.VarStmt = v
	ret.data, ret.align = resolveData(log, v.Type, v.Args)
	return ret
}

func resolveData(p lexing.Logger, t *lexing.Token, args []*lexing.Token) (
	[]byte, uint32,
) {
	switch t.Lit {
	case "str":
		return parseDataStr(p, args)
	case "x":
		return parseDataHex(p, args)
	case "u32":
		return parseDataNums(p, args, modeWord)
	case "i32":
		return parseDataNums(p, args, modeWord|modeSigned)
	case "u8", "byte":
		return parseDataNums(p, args, 0)
	case "i8":
		return parseDataNums(p, args, modeSigned)
	case "f32":
		return parseDataNums(p, args, modeWord|modeFloat)
	default:
		p.Errorf(t.Pos, "unknown data type %q", t.Lit)
		return nil, 0
	}
}

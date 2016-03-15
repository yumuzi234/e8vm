package ir

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/fmt8"
)

func printOp(p io.Writer, op Op) {
	switch op := op.(type) {
	case *Comment:
		fmt.Fprintf(p, "// %s\n", op.Str)
	case *ArithOp:
		if op.A == nil {
			if op.Op == "" {
				fmt.Fprintf(p, "%s = %s\n", op.Dest, op.B)
			} else if op.Op == "0" {
				fmt.Fprintf(p, "%s = 0\n", op.Dest)
			} else {
				fmt.Fprintf(p, "%s = %s %s\n", op.Dest, op.Op, op.B)
			}
		} else {
			fmt.Fprintf(p, "%s = %s %s %s\n",
				op.Dest, op.A, op.Op, op.B,
			)
		}
	case *CallOp:
		var args string
		if op.Args != nil {
			args = fmt8.Join(op.Args, ",")
		}
		fmt.Fprintf(p, "%s = %s(%s)\n", op.Dest, op.F, args)
	default:
		panic(fmt.Errorf("invalid or unknown IR op: %T", op))
	}
}

func printJump(p io.Writer, j *blockJump) {
	if j == nil {
		return
	}

	switch j.typ {
	case jmpAlways:
		fmt.Fprintf(p, "goto %s\n", j.to)
	case jmpIf:
		fmt.Fprintf(p, "if %s goto %s\n", j.cond, j.to)
	case jmpIfNot:
		fmt.Fprintf(p, "if !%s goto %s\n", j.cond, j.to)
	default:
		panic("invalid jump type")
	}
}

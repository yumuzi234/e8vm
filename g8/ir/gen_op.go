package ir

import (
	"math"
)

var basicOpFuncs = map[string]func(
	dest, r1, r2 uint32,
) uint32{
	"+":   asm.add,
	"-":   asm.sub,
	"*":   asm.mul,
	"/":   asm.div,
	"%":   asm.mod,
	"&":   asm.and,
	"|":   asm.or,
	"^":   asm.xor,
	"nor": asm.nor,
}

func genArithOp(g *gener, b *Block, op *arithOp) {
	if op.dest == nil {
		panic("arith with no destination")
	}

	if op.op == "0" {
		zeroRef(g, b, op.dest)
		return
	}

	if op.a != nil {
		// binary arith op
		loadRef(b, _4, op.a)
		loadRef(b, _1, op.b)

		fn := basicOpFuncs[op.op]
		if fn != nil {
			b.inst(fn(_4, _4, _1))
		} else {
			switch op.op {
			case "==":
				b.inst(asm.xor(_4, _4, _1))  // the diff
				b.inst(asm.sltu(_4, _0, _4)) // if _4 is 0, _4 <= 0
				b.inst(asm.xori(_4, _4, 1))  // flip
			case "!=":
				b.inst(asm.xor(_4, _4, _1))  // the diff
				b.inst(asm.sltu(_4, _0, _4)) // if _4 is 0, _4 <= 0
			case ">":
				b.inst(asm.slt(_4, _1, _4))
			case "<":
				b.inst(asm.slt(_4, _4, _1)) // delta = b - a
			case ">=":
				b.inst(asm.slt(_4, _4, _1))
				b.inst(asm.xori(_4, _4, 1)) // flip
			case "<=":
				b.inst(asm.slt(_4, _1, _4))
				b.inst(asm.xori(_4, _4, 1)) // flip
			case "u>":
				b.inst(asm.sltu(_4, _1, _4))
			case "u<":
				b.inst(asm.sltu(_4, _4, _1))
			case "u>=":
				b.inst(asm.sltu(_4, _4, _1))
				b.inst(asm.xori(_4, _4, 1))
			case "u<=":
				b.inst(asm.sltu(_4, _1, _4))
				b.inst(asm.xori(_4, _4, 1))
			case "<<":
				b.inst(asm.sllv(_4, _4, _1))
			case ">>":
				b.inst(asm.srla(_4, _4, _1))
			case "u>>":
				b.inst(asm.srlv(_4, _4, _1))
			default:
				panic("unknown arith op: " + op.op)
			}
		}

		saveRef(b, _4, op.dest, _1)
	} else if op.op == "" {
		copyRef(g, b, op.dest, op.b, false)
	} else if op.op == "cast" {
		loadRef(b, _4, op.b)
		saveRef(b, _4, op.dest, _1)
	} else if op.op == "makeStr" {
		s := op.b.(*strConst)
		n := len(s.str)
		if n > 0 {
			if n > math.MaxUint32-1 {
				panic("string too long")
			}
			loadSym(b, _4, s.pkg, s.name)
			loadAddr(b, _1, op.dest)
			b.inst(asm.sw(_4, _1, 0))
			loadUint32(b, _4, uint32(n))
			b.inst(asm.sw(_4, _1, 4))
		} else {
			loadAddr(b, _1, op.dest)
			b.inst(asm.sw(_0, _1, 0))
			b.inst(asm.sw(_0, _1, 4))
		}
	} else {
		// other unary arith op
		switch op.op {
		case "-":
			loadRef(b, _4, op.b)
			b.inst(asm.sub(_4, _0, _4))
		case "!":
			loadRef(b, _4, op.b)
			b.inst(asm.sltu(_4, _0, _4)) // test non-zero first
			b.inst(asm.xori(_4, _4, 1))  // and flip
		case "?", "?f": // test if it is non-zero
			loadRef(b, _4, op.b)
			b.inst(asm.sltu(_4, _0, _4))
		case "^":
			loadRef(b, _4, op.b)
			b.inst(asm.nor(_4, _0, _4))
		case "&": // fetches the address of the block
			loadAddr(b, _4, op.b)
		case "<0":
			loadRef(b, _4, op.b)
			b.inst(asm.slt(_4, _4, _0))
		case "*":
			panic("op * is deprecated, please use NewAddrRef()")
		default:
			panic("unknown arith unary op: " + op.op)
		}

		saveRef(b, _4, op.dest, _1)
	}
}

func genOp(g *gener, b *Block, op op) {
	switch op := op.(type) {
	case *arithOp:
		genArithOp(g, b, op)
	case *callOp:
		genCallOp(g, b, op)
	case *comment:
		// do nothing
	default:
		panic("unknown op type")
	}
}

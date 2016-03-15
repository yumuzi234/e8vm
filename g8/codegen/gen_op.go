package codegen

import (
	"math"
)

var basicOpFuncs = map[string]func(
	dest, r1, r2 uint32,
) uint32{
	"+":   asm.add,
	"-":   asm.sub,
	"*":   asm.mul,
	"u*":  asm.mulu,
	"/":   asm.div,
	"u/":  asm.divu,
	"%":   asm.mod,
	"u%":  asm.modu,
	"&":   asm.and,
	"|":   asm.or,
	"^":   asm.xor,
	"nor": asm.nor,
}

func genArithOp(g *gener, b *Block, op *ArithOp) {
	if op.Dest == nil {
		panic("arith with no destination")
	}

	if op.Op == "0" {
		zeroRef(g, b, op.Dest)
		return
	}

	if op.A != nil {
		// binary arith op
		loadRef(b, _4, op.A)
		loadRef(b, _1, op.B)

		fn := basicOpFuncs[op.Op]
		if fn != nil {
			b.inst(fn(_4, _4, _1))
		} else {
			switch op.Op {
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
				panic("unknown arith op: " + op.Op)
			}
		}

		saveRef(b, _4, op.Dest, _1)
	} else if op.Op == "" {
		copyRef(g, b, op.Dest, op.B, false)
	} else if op.Op == "cast" {
		loadRef(b, _4, op.B)
		saveRef(b, _4, op.Dest, _1)
	} else if op.Op == "makeStr" {
		s := op.B.(*strConst)
		n := len(s.str)
		if n > 0 {
			if n > math.MaxInt32-1 {
				panic("string too long")
			}
			loadSym(b, _4, s.pkg, s.name)
			loadAddr(b, _1, op.Dest)
			b.inst(asm.sw(_4, _1, 0))
			loadUint32(b, _4, uint32(n))
			b.inst(asm.sw(_4, _1, 4))
		} else {
			loadAddr(b, _1, op.Dest)
			b.inst(asm.sw(_0, _1, 0))
			b.inst(asm.sw(_0, _1, 4))
		}
	} else if op.Op == "makeDat" {
		d := op.B.(*heapDat)
		n := d.n
		if n > 0 {
			if n > math.MaxInt32-1 {
				panic("dat too long")
			}
			loadSym(b, _4, d.pkg, d.name)
			loadAddr(b, _1, op.Dest)
			b.inst(asm.sw(_4, _1, 0))
			loadUint32(b, _4, uint32(n))
			b.inst(asm.sw(_4, _1, 4))
		} else {
			loadAddr(b, _1, op.Dest)
			b.inst(asm.sw(_0, _1, 0))
			b.inst(asm.sw(_0, _1, 4))
		}
	} else {
		// other unary arith op
		switch op.Op {
		case "-":
			loadRef(b, _4, op.B)
			b.inst(asm.sub(_4, _0, _4))
		case "!":
			loadRef(b, _4, op.B)
			b.inst(asm.sltu(_4, _0, _4)) // test non-zero first
			b.inst(asm.xori(_4, _4, 1))  // and flip
		case "?", "?f": // test if it is non-zero
			loadRef(b, _4, op.B)
			b.inst(asm.sltu(_4, _0, _4))
		case "^":
			loadRef(b, _4, op.B)
			b.inst(asm.nor(_4, _0, _4))
		case "&": // fetches the address of the block
			loadAddr(b, _4, op.B)
		case "<0":
			loadRef(b, _4, op.B)
			b.inst(asm.slt(_4, _4, _0))
		case "*":
			panic("op * is deprecated, please use NewAddrRef()")
		default:
			panic("unknown arith unary op: " + op.Op)
		}

		saveRef(b, _4, op.Dest, _1)
	}
}

func genOp(g *gener, b *Block, op Op) {
	switch op := op.(type) {
	case *ArithOp:
		genArithOp(g, b, op)
	case *CallOp:
		genCallOp(g, b, op)
	case *Comment:
		// do nothing
	default:
		panic("unknown op type")
	}
}

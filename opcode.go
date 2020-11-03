package dune

import (
	"fmt"
	"io"
)

type Opcode byte

const (
	op_ldk Opcode = iota // load constant A := RK(B)
	op_mov               // A := B
	op_mob               // A := B and C with true is B is not null or empty.
	op_add               // A := B + C
	op_sub               // A := B - C
	op_mul               // A := B * C
	op_div               // A := B / C
	op_mod               // A := B % C
	op_bor               // A := B | C
	op_and               // A := B & C
	op_xor               // A := B ^ C
	op_lsh               // A := B << C
	op_rsh               // A := B >> C
	op_inc               // A = A++. B up
	op_dec               // A = A--. B up
	op_unm               // A := -RK(B)
	op_not               // A := !RK(B)
	op_bnt               // A := ~B  bitwise not
	op_str               // Set register: A register number, B register value
	op_new               // create a new instance of a class: A class type, B retAddress, C argsAddress
	op_nes               // create a new instance of a class with a single arg: A class type, B retAddress, C argsAddress
	op_arr               // create a new array:  A array, B size
	op_map               // create a new map:  A map, B size
	op_key               // get the keys of a map or the indexes of an array: A := keys(B)
	op_val               // get the values of a map or array: A := values(B)
	op_len               // get the length of an array: A := len(B)
	op_enu               // get value from enum: A dest, B enum C key
	op_get               // get array index or map key: A dest, B source C index or key
	op_gto               // get optional chaining: A dest, B source C index or key. Reg0 stores the PC to jump if B is null.
	op_set               // set array index or map key: A array or Map, B index or key, C value
	op_spa               // spread last index of array A
	op_jmp               // jump A positions
	op_jpb               // jump back A positions
	op_ejp               // jump if A and B are equal C instructions.
	op_djp               // jump if A and B are different C instructions.
	op_tjp               // test if true or not null and jump: test A and jump B instructions. C=(0=jump if true, 1 jump if false)
	op_eql               // equality test
	op_neq               // inequality test
	op_seq               // strict equality test
	op_sne               // strict inequality test
	op_lst               // less than
	op_lse               // less or equal than
	op_cal               // call: A funcIndex, B retAddress, C argsAddress
	op_cco               // call optional chaining: A funcIndex, B retAddress, C argsAddress. Reg0 stores the PC to jump if B is null.
	op_cas               // call with single argument: A funcIndex, B retAddress, C argsAddress
	op_cso               // call optional with single argument: A funcIndex, B retAddress, C argsAddress. Reg0 stores the PC to jump if B is null.
	op_rnp               // Read native property: A := B
	op_ret               // return from a call: A dest
	op_clo               // create closure: A dest R(B value) funcIndex
	op_trw               // throw. A contains the error
	op_try               // try-catch: jump to A absolute pc, set the error to B. C: the 'finally' absolute pc.
	op_tre               // try-end: set the last try body as ended.
	op_cen               // catch-end: set the last catch body as ended. It is only emmited if there is no finally
	op_fen               // finally-end: set the last finally body as ended.
	op_trx               // try exit: a continue inside try/catch inside a loop for example
	op_del               // delete object property
)

const (
	vm_next = iota
	vm_continue
	vm_exit
)

func exec(i *Instruction, vm *VM) int {
	switch i.Opcode {
	case op_ldk:
		return exec_ldk(i, vm)

	case op_mov:
		return exec_mov(i, vm)

	case op_mob:
		return exec_mob(i, vm)

	case op_add:
		return exec_add(i, vm)

	case op_sub:
		return exec_sub(i, vm)

	case op_mul:
		return exec_mul(i, vm)

	case op_div:
		return exec_div(i, vm)

	case op_mod:
		return exec_mod(i, vm)

	case op_bor:
		return exec_bor(i, vm)

	case op_and:
		return exec_and(i, vm)

	case op_xor:
		return exec_xor(i, vm)

	case op_lsh:
		return exec_lsh(i, vm)

	case op_rsh:
		return exec_rsh(i, vm)

	case op_inc:
		return exec_inc(i, vm)

	case op_dec:
		return exec_dec(i, vm)

	case op_unm:
		return exec_unm(i, vm)

	case op_not:
		return exec_not(i, vm)

	case op_bnt:
		return exec_bnt(i, vm)

	case op_str:
		return exec_str(i, vm)

	case op_new:
		return exec_new(i, vm)

	case op_nes:
		return exec_nes(i, vm)

	case op_arr:
		return exec_arr(i, vm)

	case op_map:
		return exec_map(i, vm)

	case op_key:
		return exec_key(i, vm)

	case op_val:
		return exec_val(i, vm)

	case op_len:
		return exec_len(i, vm)

	case op_enu:
		return exec_enu(i, vm)

	case op_get:
		return exec_get(i, vm)

	case op_gto:
		return exec_gto(i, vm)

	case op_set:
		return exec_set(i, vm)

	case op_spa:
		return exec_spa(i, vm)

	case op_jmp:
		return exec_jmp(i, vm)

	case op_jpb:
		return exec_jpb(i, vm)

	case op_ejp:
		return exec_ejp(i, vm)

	case op_djp:
		return exec_djp(i, vm)

	case op_tjp:
		return exec_tjp(i, vm)

	case op_eql:
		return exec_eql(i, vm)

	case op_neq:
		return exec_neq(i, vm)

	case op_seq:
		return exec_seq(i, vm)

	case op_sne:
		return exec_sne(i, vm)

	case op_lst:
		return exec_lst(i, vm)

	case op_lse:
		return exec_lse(i, vm)

	case op_cal:
		return exec_cal(i, vm)

	case op_cco:
		return exec_cco(i, vm)

	case op_cas:
		return exec_cas(i, vm)

	case op_cso:
		return exec_cso(i, vm)

	case op_rnp:
		return exec_rnp(i, vm)

	case op_ret:
		return exec_ret(i, vm)

	case op_clo:
		return exec_clo(vm)

	case op_trw:
		return exec_trw(i, vm)

	case op_try:
		return exec_try(i, vm)

	case op_tre:
		return exec_tre(vm)

	case op_cen:
		return exec_cen(vm)

	case op_fen:
		return exec_fen(vm)

	case op_trx:
		return exec_trx(vm)

	case op_del:
		return exec_del(i, vm)

	default:
		panic(fmt.Sprintf("Invalid opcode: %v", i))
	}
}

func exec_mov(instr *Instruction, vm *VM) int {
	vm.set(instr.A, vm.get(instr.B))
	return vm_next
}

func exec_mob(instr *Instruction, vm *VM) int {
	// set A with B if it has instr.A value and C with true is set.
	bv := vm.get(instr.B)
	vm.set(instr.A, bv)
	switch bv.Type {
	case Bool:
		vm.set(instr.C, bv)

	case Int:
		if bv.ToInt() == 0 {
			vm.set(instr.C, FalseValue)
		} else {
			vm.set(instr.C, TrueValue)
		}

	case Float:
		if bv.ToFloat() == 0 {
			vm.set(instr.C, FalseValue)
		} else {
			vm.set(instr.C, TrueValue)
		}

	default:
		if bv.IsNilOrEmpty() {
			vm.set(instr.C, FalseValue)
		} else {
			vm.set(instr.C, TrueValue)
		}
	}

	return vm_next
}

func exec_ldk(instr *Instruction, vm *VM) int {
	k := vm.Program.Constants[instr.B.Value]
	vm.set(instr.A, k)
	return vm_next
}

func exec_add(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)
	switch lh.Type {
	case Int:
		switch rh.Type {
		case Float:
			vm.set(instr.A, NewFloat(lh.ToFloat()+rh.ToFloat()))
		case Int:
			vm.set(instr.A, NewInt64(lh.ToInt()+rh.ToInt()))
		case Rune:
			vm.set(instr.A, NewRune(lh.ToRune()+rh.ToRune()))
		case String:
			err := vm.AddAllocations(lh.Size())
			if err == nil {
				err = vm.AddAllocations(rh.Size())
			}
			if err != nil {
				if vm.handle(err) {
					return vm_continue
				} else {
					return vm_exit
				}
			}
			vm.set(instr.A, NewString(lh.ToString()+rh.ToString()))
		default:
			if vm.handle(vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type)) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Float:
		switch rh.Type {
		case Int, Float:
			vm.set(instr.A, NewFloat(lh.ToFloat()+rh.ToFloat()))
		case Rune:
			vm.set(instr.A, NewRune(lh.ToRune()+rh.ToRune()))
		case String:
			err := vm.AddAllocations(lh.Size())
			if err == nil {
				err = vm.AddAllocations(rh.Size())
			}
			if err != nil {
				if vm.handle(err) {
					return vm_continue
				} else {
					return vm_exit
				}
			}
			vm.set(instr.A, NewString(lh.ToString()+rh.ToString()))
		default:
			if vm.handle(vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type)) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Rune:
		switch rh.Type {
		case Rune, Int:
			vm.set(instr.A, NewRune(lh.ToRune()+rh.ToRune()))
		case String:
			err := vm.AddAllocations(lh.Size())
			if err == nil {
				err = vm.AddAllocations(rh.Size())
			}
			if err != nil {
				if vm.handle(err) {
					return vm_continue
				} else {
					return vm_exit
				}
			}
			vm.set(instr.A, NewString(lh.ToString()+rh.ToString()))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Bool:
		switch rh.Type {
		case String:
			vm.set(instr.A, NewString(lh.ToString()+rh.ToString()))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case String:
		switch rh.Type {
		case String, Int, Float, Bool, Rune:
			err := vm.AddAllocations(lh.Size())
			if err == nil {
				err = vm.AddAllocations(rh.Size())
			}
			if err != nil {
				if vm.handle((err)) {
					return vm_continue
				} else {
					return vm_exit
				}
			}
			vm.set(instr.A, NewString(lh.ToString()+rh.ToString()))
		case Null:
			vm.set(instr.A, lh)
		case Undefined:
			vm.set(instr.A, NewString(lh.ToString()+"undefined"))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Null:
		switch rh.Type {
		case Null, String, Int, Float:
			vm.set(instr.A, rh)
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Undefined:
		switch rh.Type {
		case Null, String, Int, Float:
			vm.set(instr.A, rh)
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	default:
		if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}

	return vm_next
}

func exec_sub(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)
	switch lh.Type {
	case Int:
		switch rh.Type {
		case Float:
			vm.set(instr.A, NewFloat(lh.ToFloat()-rh.ToFloat()))
		case Int:
			vm.set(instr.A, NewInt64(lh.ToInt()-rh.ToInt()))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Float:
		switch rh.Type {
		case Int, Float:
			vm.set(instr.A, NewFloat(lh.ToFloat()-rh.ToFloat()))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Rune:
		switch rh.Type {
		case Rune, Int:
			vm.set(instr.A, NewRune(lh.ToRune()-rh.ToRune()))
		case String:
			rs := rh.ToString()
			if len(rs) != 1 {
				if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
					return vm_continue
				} else {
					return vm_exit
				}
			}
			vm.set(instr.A, NewRune(lh.ToRune()-rune(rs[0])))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case String:
		ls := lh.ToString()
		if len(ls) != 1 {
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
		lr := rune(ls[0])
		switch rh.Type {
		case Rune, Int:
			vm.set(instr.A, NewRune(lr-rh.ToRune()))
		case String:
			rs := rh.ToString()
			if len(rs) != 1 {
				if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
					return vm_continue
				} else {
					return vm_exit
				}
			}
			vm.set(instr.A, NewRune(lr-rune(rs[0])))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	default:
		if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_mul(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)
	switch lh.Type {
	case Int:
		switch rh.Type {
		case Float:
			vm.set(instr.A, NewFloat(lh.ToFloat()*rh.ToFloat()))
		case Int:
			vm.set(instr.A, NewInt64(lh.ToInt()*rh.ToInt()))
		case Rune:
			vm.set(instr.A, NewRune(lh.ToRune()*rh.ToRune()))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Float:
		switch rh.Type {
		case Int, Float:
			vm.set(instr.A, NewFloat(lh.ToFloat()*rh.ToFloat()))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Rune:
		switch rh.Type {
		case Rune, Int:
			vm.set(instr.A, NewRune(lh.ToRune()*rh.ToRune()))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	default:
		if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}

	return vm_next
}

func exec_bor(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)
	switch lh.Type {
	case Int:
		switch rh.Type {
		case Int:
			vm.set(instr.A, NewInt64(lh.ToInt()|rh.ToInt()))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	default:
		if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_lsh(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)
	switch lh.Type {
	case Int:
		switch rh.Type {
		case Int:
			vm.set(instr.A, NewInt64(int64(lh.ToInt()<<uint64(rh.ToInt()))))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	default:
		if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_rsh(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)
	switch lh.Type {
	case Int:
		switch rh.Type {
		case Int:
			vm.set(instr.A, NewInt64(int64(lh.ToInt()>>uint64(rh.ToInt()))))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	default:
		if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_xor(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)
	switch lh.Type {
	case Int:
		switch rh.Type {
		case Int:
			vm.set(instr.A, NewInt(int(lh.ToInt())^int(rh.ToInt())))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	default:
		if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_and(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)
	switch lh.Type {
	case Int:
		switch rh.Type {
		case Int:
			vm.set(instr.A, NewInt(int(lh.ToInt())&int(rh.ToInt())))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	default:
		if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_div(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)

	if lh.Type == Rune || rh.Type == Rune {
		switch lh.Type {
		case Int, Rune:
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
		switch rh.Type {
		case Int, Rune:
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
		vm.set(instr.A, NewRune(lh.ToRune()/rh.ToRune()))
	} else {
		switch lh.Type {
		case Int, Float:
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}

		switch rh.Type {
		case Int, Float:
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
		rf := rh.ToFloat()
		if rf == 0 {
			if vm.handle((vm.NewError("Attempt to divide by zero"))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
		vm.set(instr.A, NewFloat(lh.ToFloat()/rf))
	}
	return vm_next
}

func exec_mod(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)
	switch lh.Type {
	case Float:
		if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	case Int:
		switch rh.Type {
		case Int:
			ri := rh.ToInt()
			if ri == 0 {
				if vm.handle((vm.NewError("Attempt to divide by zero"))) {
					return vm_continue
				} else {
					return vm_exit
				}
			}
			vm.set(instr.A, NewInt64(lh.ToInt()%ri))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Rune:
		switch rh.Type {
		case Rune, Int:
			ri := rh.ToRune()
			if ri == 0 {
				if vm.handle((vm.NewError("Attempt to divide by zero"))) {
					return vm_continue
				} else {
					return vm_exit
				}
			}
			vm.set(instr.A, NewRune(lh.ToRune()%ri))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	default:
		if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_str(instr *Instruction, vm *VM) int {
	if instr.A.Kind != AddrData {
		panic(fmt.Sprintf("Compiler error: invalid register number kind: %v", instr.A))
	}

	if instr.B.Kind != AddrData {
		panic(fmt.Sprintf("Compiler error: invalid register value kind: %v", instr.B))
	}

	switch instr.A.Value {
	case 0:
		vm.reg0 = instr.B.Value
	default:
		panic(fmt.Sprintf("Compiler error: invalid register number: %v", instr.A))
	}

	return vm_next
}

func exec_bnt(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	switch lh.Type {
	case Int:
		vm.set(instr.A, NewInt64(^lh.ToInt()))
	default:
		if vm.handle((vm.NewError("Invalid operation on %v", lh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_unm(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	switch lh.Type {
	case Int:
		vm.set(instr.A, NewInt64(lh.ToInt()*-1))
	case Float:
		vm.set(instr.A, NewFloat(lh.ToFloat()*-1))
	default:
		if vm.handle((vm.NewError("Invalid operation on %v", lh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_not(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	switch lh.Type {
	case Bool:
		vm.set(instr.A, NewBool(!lh.ToBool()))
	default:
		var empty bool
		if lh.Type == Int {
			// if the value is 0 treat it as null or empty like in javascript.
			empty = lh.ToInt() == 0
		} else {
			empty = lh.IsNilOrEmpty() // true if it has a value like in javascript
		}
		vm.set(instr.A, NewBool(empty))
	}
	return vm_next
}

func exec_inc(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.A)
	switch lh.Type {
	case Int:
		vm.set(instr.A, NewInt64(lh.ToInt()+1))
	case Float:
		vm.set(instr.A, NewFloat(lh.ToFloat()+1))
	default:
		if vm.handle((vm.NewError("Invalid operation on %v", lh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_dec(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.A)
	switch lh.Type {
	case Int:
		vm.set(instr.A, NewInt64(lh.ToInt()-1))
	case Float:
		vm.set(instr.A, NewFloat(lh.ToFloat()-1))
	default:
		if vm.handle((vm.NewError("Invalid operation on %v", lh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_lst(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)

	rh := vm.get(instr.C)
	switch lh.Type {
	case Int:
		switch rh.Type {
		case Float:
			vm.set(instr.A, NewBool(lh.ToFloat() < rh.ToFloat()))
		case Int:
			vm.set(instr.A, NewBool(lh.ToInt() < rh.ToInt()))
		case Null:
			vm.set(instr.A, NewBool(0 < rh.ToInt()))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Float:
		switch rh.Type {
		case Int, Float:
			vm.set(instr.A, NewBool(lh.ToFloat() < rh.ToFloat()))
		case Null:
			vm.set(instr.A, NewBool(0 < rh.ToInt()))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Rune:
		switch rh.Type {
		case Rune, Int:
			vm.set(instr.A, NewBool(lh.ToRune() < rh.ToRune()))
		case String:
			rs := rh.ToString()
			if len(rs) != 1 {
				if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
					return vm_continue
				} else {
					return vm_exit
				}
			}
			rr := rune(rs[0])
			vm.set(instr.A, NewBool(lh.ToRune() < rr))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case String:
		switch rh.Type {
		case String:
			vm.set(instr.A, NewBool(lh.ToString() < rh.ToString()))
		case Rune:
			ls := lh.ToString()
			if len(ls) != 1 {
				if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
					return vm_continue
				} else {
					return vm_exit
				}
			}
			lr := rune(ls[0])
			vm.set(instr.A, NewBool(lr < rh.ToRune()))
		case Null:
			vm.set(instr.A, FalseValue)
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Null:
		switch rh.Type {
		case Int, Float:
			vm.set(instr.A, NewBool(0 < rh.ToFloat()))
		case String, Rune:
			vm.set(instr.A, FalseValue)
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	default:
		if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_lse(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)
	switch lh.Type {
	case Int:
		switch rh.Type {
		case Float:
			vm.set(instr.A, NewBool(lh.ToFloat() <= rh.ToFloat()))
		case Int:
			vm.set(instr.A, NewBool(lh.ToInt() <= rh.ToInt()))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Float:
		switch rh.Type {
		case Int, Float:
			vm.set(instr.A, NewBool(lh.ToFloat() <= rh.ToFloat()))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case Rune:
		switch rh.Type {
		case Rune, Int:
			vm.set(instr.A, NewBool(lh.ToRune() <= rh.ToRune()))
		case String:
			rs := rh.ToString()
			if len(rs) != 1 {
				if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
					return vm_continue
				} else {
					return vm_exit
				}
			}
			rr := rune(rs[0])
			vm.set(instr.A, NewBool(lh.ToRune() <= rr))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	case String:
		switch rh.Type {
		case String:
			vm.set(instr.A, NewBool(lh.ToString() <= rh.ToString()))
		case Rune:
			ls := lh.ToString()
			if len(ls) != 1 {
				if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
					return vm_continue
				} else {
					return vm_exit
				}
			}
			lr := rune(ls[0])
			vm.set(instr.A, NewBool(lr <= rh.ToRune()))
		default:
			if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}

	default:
		if vm.handle((vm.NewError("Invalid operation on %v and %v", lh.Type, rh.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_eql(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)
	vm.set(instr.A, NewBool(lh.Equals(rh)))
	return vm_next
}

func exec_neq(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)
	vm.set(instr.A, NewBool(!lh.Equals(rh)))
	return vm_next
}

func exec_seq(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)
	vm.set(instr.A, NewBool(lh.StrictEquals(rh)))
	return vm_next
}

func exec_sne(instr *Instruction, vm *VM) int {
	lh := vm.get(instr.B)
	rh := vm.get(instr.C)
	vm.set(instr.A, NewBool(!lh.StrictEquals(rh)))
	return vm_next
}

func exec_arr(instr *Instruction, vm *VM) int {
	vm.set(instr.A, NewArray(int(instr.B.Value)))
	return vm_next
}

func exec_map(instr *Instruction, vm *VM) int {
	vm.set(instr.A, NewMap(int(instr.B.Value)))
	return vm_next
}

func exec_enu(instr *Instruction, vm *VM) int {
	enum := vm.Program.Enums[int(instr.B.Value)]
	value := enum.Values[int(instr.C.Value)]
	k := vm.Program.Constants[value.KIndex]
	vm.set(instr.A, k)
	return vm_next
}

func exec_get(instr *Instruction, vm *VM) int {
	if _, err := vm.getFromObject(instr, true); err != nil {
		if vm.handle((err)) {
			return vm_continue
		} else {
			vm.Error = err
			return vm_exit
		}
	}

	// set value in an array or map: A array, B index, C value
	return vm_next
}

func exec_gto(instr *Instruction, vm *VM) int {
	ok, err := vm.getFromObject(instr, false)
	if err != nil {
		if vm.handle((err)) {
			return vm_continue
		} else {
			vm.Error = err
			return vm_exit
		}
	}

	if !ok {
		vm.incPC(int(vm.reg0))
		return vm_continue
	}

	// set value in an array or map: A array, B index, C value
	return vm_next
}

func exec_set(instr *Instruction, vm *VM) int {
	if err := vm.setToObject(instr); err != nil {
		if vm.handle((err)) {
			return vm_continue
		} else {
			vm.Error = err
			return vm_exit
		}
	}
	return vm_next
}

func exec_spa(instr *Instruction, vm *VM) int {
	v := vm.get(instr.A)
	if v.Type != Array {
		if vm.handle((vm.NewError("Expected array, got %v", v.TypeName()))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}

	va := v.ToArrayObject().Array

	ln := len(va)
	if ln == 0 {
		return vm_next
	}

	last := va[ln-1]
	switch last.Type {
	case Null, Undefined:
		vm.set(instr.A, NewArrayValues(va[:ln-1]))
	case Array:
		n := append(va[:ln-1], last.ToArrayObject().Array...)
		vm.set(instr.A, NewArrayValues(n))
	default:
		if vm.handle((vm.NewError("Expected array, got %v", last.TypeName()))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}

	return vm_next
}

func exec_key(instr *Instruction, vm *VM) int {
	// gets the keys of a map or the indexes of an array: A := keys(B)
	bv := vm.get(instr.B)
	switch bv.Type {

	case Null:
		// allow to iterate if not initialize (set an empty array)
		vm.set(instr.A, NewArray(0))

	case Array:
		s := bv.ToArray()
		ln := len(s)
		values := make([]Value, ln)
		for i := 0; i < ln; i++ {
			values[i] = NewInt(i)
		}
		vm.set(instr.A, NewArrayValues(values))

	case Map:
		m := bv.ToMap()
		m.RLock()
		s := m.Map
		values := make([]Value, len(s))
		i := 0
		for k := range s {
			values[i] = k
			i++
		}
		m.RUnlock()
		vm.set(instr.A, NewArrayValues(values))

	case Enum:
		i := bv.ToEnum()
		enum := vm.Program.Enums[i]
		ln := len(enum.Values)
		values := make([]Value, ln)
		for i := 0; i < ln; i++ {
			values[i] = NewInt(i)
		}
		vm.set(instr.A, NewArrayValues(values))

	case Object:
		obj := bv.ToObject()
		if n, ok := obj.(IterableByIndex); ok {
			ln := n.Len()
			values := make([]Value, ln)
			for i := 0; i < ln; i++ {
				values[i] = NewInt(i)
			}
			vm.set(instr.A, NewArrayValues(values))
		} else {
			if vm.handle((vm.NewError("Expected a key or index enumerable, got %v", bv.TypeName()))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}

	default:
		if vm.handle((vm.NewError("Expected a enumerable, got %v", bv.TypeName()))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_val(instr *Instruction, vm *VM) int {
	// gets the values of a map or array: A := values(B)
	bv := vm.get(instr.B)
	switch bv.Type {

	case Null, Undefined:
		// allow to iterate if not initialize (set an empty array)
		vm.set(instr.A, NewArray(0))

	case Array:
		// copiar los valores para que si se modifican dentro de un loop no afecten a la iteración
		s := bv.ToArray()
		values := make([]Value, len(s))
		copy(values, s)
		vm.set(instr.A, NewArrayValues(values))
	case Bytes:
		s := bv.ToBytes()
		values := make([]Value, len(s))
		for i, v := range s {
			values[i] = NewInt(int(v))
		}
		vm.set(instr.A, NewArrayValues(values))
	case Map:
		m := bv.ToMap()
		m.RLock()
		s := m.Map
		values := make([]Value, len(s))
		i := 0
		for _, v := range s {
			values[i] = v
			i++
		}
		m.RUnlock()
		vm.set(instr.A, NewArrayValues(values))
	case Object:
		obj := bv.ToObject()
		if enum, ok := obj.(Enumerable); ok {
			vals, err := enum.Values()
			if err != nil {
				if vm.handle((vm.NewError("Enumerable error: %v", err))) {
					return vm_continue
				} else {
					return vm_exit
				}
			} else {
				vm.set(instr.A, NewArrayValues(vals))
			}
		} else if vm.handle((vm.NewError("Expected a enumerable, got %v", bv.String()))) {
			return vm_continue
		} else {
			return vm_exit
		}

	default:
		if vm.handle((vm.NewError("Expected a enumerable, got %v", bv.String()))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_len(instr *Instruction, vm *VM) int {
	bv := vm.get(instr.B)
	switch bv.Type {
	case Array:
		vm.set(instr.A, NewInt(len(bv.ToArray())))
	case Map:
		m := bv.ToMap()
		m.RLock()
		vm.set(instr.A, NewInt(len(m.Map)))
		m.RUnlock()
	case Object:
		if col, ok := bv.ToObject().(IterableByIndex); ok {
			vm.set(instr.A, NewInt(col.Len()))
		} else {
			if vm.handle((vm.NewError("The value is not a collection %v", bv.Type))) {
				return vm_continue
			} else {
				return vm_exit
			}
		}
	default:
		if vm.handle((vm.NewError("The value is not a collection %v", bv.Type))) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

// Read native property A := B
func exec_rnp(instr *Instruction, vm *VM) int {
	n := vm.get(instr.B)

	i := n.ToNativeFunction()

	if err := vm.callNativeFunc(i, nil, instr.A, NullValue); err != nil {
		if vm.handle(vm.WrapError(err)) {
			return vm_continue
		} else {
			return vm_exit
		}
	}
	return vm_next
}

func exec_trw(instr *Instruction, vm *VM) int {
	v := vm.get(instr.A)

	var err error
	if v.Type == Object {
		if e, ok := v.ToObject().(Error); ok {
			// don't alter the stack trace if it is a re throw
			if e.IsRethrow {
				err = e
			} else {
				err = vm.WrapError(e)
			}
		}
	}

	if err == nil {
		err = vm.NewError(v.String())
	}

	// check if is inside a catch to discard it.
	l := len(vm.tryCatchs) - 1
	if l < 0 {
		// run finalizers before exiting
		currentFrame := vm.callStack[vm.fp]
		vm.runFinalizables(currentFrame)

		// an unhandled error
		vm.Error = err
		return vm_exit
	}

	// check if we are in the finally block.
	// remove the current try if we are and let the next one handle it.
	try := vm.tryCatchs[l]
	if try.finallyExecuted {
		vm.tryCatchs = vm.tryCatchs[:l]
	}

	if vm.handle((err)) {
		return vm_continue
	} else {
		return vm_exit
	}
}

func exec_try(instr *Instruction, vm *VM) int {
	//  jump to A absolute pc, set the error to B. C: the 'finally' absolute pc.
	var catchPC int
	if instr.A.Kind == AddrVoid {
		catchPC = -1
	} else {
		catchPC = int(instr.A.Value)
	}

	try := &tryCatch{
		catchPC:  catchPC,
		retPC:    -1,
		errorReg: instr.B,
		fp:       vm.fp,
	}

	// set the finally pc if provided
	if instr.C.Kind == AddrData {
		try.finallyPC = int(instr.C.Value)
	} else {
		try.finallyPC = -1
	}

	vm.tryCatchs = append(vm.tryCatchs, try)
	return vm_next
}

func exec_trx(vm *VM) int {
	i := len(vm.tryCatchs) - 1

	if i >= 0 {
		try := vm.tryCatchs[i]

		// if there is no finally just remove it
		if try.finallyPC == -1 {
			vm.tryCatchs = vm.tryCatchs[:i]
			return vm_next
		}

		// make it continue to the next instruction after the finally ends
		try.retPC = vm.callStack[vm.fp].pc + 1

		// advance to the finally part
		vm.setPC(try.finallyPC)

		return vm_continue
	}

	return vm_continue
}

func exec_tre(vm *VM) int {
	l := len(vm.tryCatchs) - 1

	try := vm.tryCatchs[l]
	if try.finallyPC == -1 {
		// if there is no finally, discard it
		vm.tryCatchs = vm.tryCatchs[:l]
	}
	return vm_next
}

func exec_cen(vm *VM) int {
	l := len(vm.tryCatchs) - 1

	// don't need to check finally because cen is only emmited if there is no finally
	vm.tryCatchs = vm.tryCatchs[:l]
	return vm_next
}

func exec_fen(vm *VM) int {
	l := len(vm.tryCatchs) - 1
	try := vm.tryCatchs[l]
	vm.tryCatchs = vm.tryCatchs[:l]

	// if the error was unhandled because there was no catch block
	// the handle it now that the finally has been processed.
	if try.err != nil {
		vm.incPC(1)
		if vm.handle(try.err) {
			return vm_continue
		} else {
			return vm_exit
		}
	}

	if try.retPC != -1 {
		vm.setPC(try.retPC)
		return vm_continue
	}

	return vm_next
}

func exec_ejp(instr *Instruction, vm *VM) int {
	// jump if A and B are equal C instructions.

	lh := vm.get(instr.A)
	rh := vm.get(instr.B)

	if lh.Equals(rh) {
		vm.incPC(int(instr.C.Value))
	}
	return vm_next
}

func exec_djp(instr *Instruction, vm *VM) int {
	// jump if A and B are different C instructions.

	lh := vm.get(instr.A)
	rh := vm.get(instr.B)

	if !lh.Equals(rh) {
		vm.incPC(int(instr.C.Value))
	}
	return vm_next
}

func exec_tjp(instr *Instruction, vm *VM) int {
	// test if true or not null and jump: test A and jump B instructions. C=(0=jump if true, 1 jump if false)

	av := vm.get(instr.A)
	var expr bool
	switch av.Type {
	case Bool:
		expr = av.ToBool()
	case Int:
		// if the value is 0 treat it as null or empty like in javascript.
		expr = av.ToInt() != 0
	case Float:
		// if the value is 0 treat it as null or empty like in javascript.
		expr = av.ToFloat() != 0
	default:
		expr = !av.IsNilOrEmpty() // true if it has a value like in javascript
	}

	cv := instr.C.Value

	switch jumpType(cv) {
	case jumpIfFalse:
		if expr {
			vm.incPC(int(instr.B.Value))
		}
	case jumpIfTrue:
		if !expr {
			vm.incPC(int(instr.B.Value))
		}
	case jumpIfNotNull:
		if !av.IsNil() {
			vm.incPC(int(instr.B.Value))
		}
	}

	return vm_next
}

func exec_jmp(instr *Instruction, vm *VM) int {
	vm.incPC(int(instr.A.Value))
	return vm_next
}

func exec_jpb(instr *Instruction, vm *VM) int {
	vm.incPC(int(instr.A.Value) * -1)
	return vm_continue
}

func exec_clo(vm *VM) int {
	// R(A) dest R(B value) funcIndex
	instr := vm.instruction()
	funcIndex := instr.B.Value

	// copy  closures carried from parent functions
	frame := vm.callStack[vm.fp]
	f := vm.Program.Functions[frame.funcIndex]
	fLen := len(f.Closures)
	frLen := len(frame.closures)

	// mark it so it is not reused
	frame.inClosure = true

	c := &Closure{
		FuncIndex: int(funcIndex),
		closures:  make([]*closureRegister, fLen+frLen),
	}

	copy(c.closures, frame.closures)

	// copy closures defined in this function.
	for i, r := range f.Closures {
		c.closures[frLen+i] = &closureRegister{register: r, values: frame.values}
	}

	vm.set(instr.A, NewObject(c))
	return vm_next
}

func exec_new(instr *Instruction, vm *VM) int {
	// A class index, B retAddress, C argsAddress

	var args []Value
	if instr.C != Void {
		args = vm.get(instr.C).ToArrayObject().Array
	}

	i := newInstance(instr.A, vm)

	v := NewObject(i)
	vm.set(instr.B, v)

	f, ok := i.Function("constructor", vm.Program)
	if ok {
		return vm.callProgramFunc(f, Void, args, true, v, nil)
	}
	return vm_next
}

func exec_nes(instr *Instruction, vm *VM) int {
	// A class index, B retAddress, C argsAddress

	args := []Value{vm.get(instr.C)}

	i := newInstance(instr.A, vm)

	v := NewObject(i)
	vm.set(instr.B, v)

	f, ok := i.Function("constructor", vm.Program)
	if ok {
		return vm.callProgramFunc(f, Void, args, true, v, nil)
	}

	return vm_next
}

func exec_cal(instr *Instruction, vm *VM) int {
	// A funcIndex, B retAddress, C argsAddress

	var args []Value
	if instr.C != Void {
		args = vm.get(instr.C).ToArrayObject().Array
	}

	return vm.call(instr.A, instr.B, args, false)
}

func exec_cco(instr *Instruction, vm *VM) int {
	// A funcIndex, B retAddress, C argsAddress

	var args []Value
	if instr.C != Void {
		args = vm.get(instr.C).ToArrayObject().Array
	}

	return vm.call(instr.A, instr.B, args, true)
}

func exec_cas(instr *Instruction, vm *VM) int {
	// A funcIndex, B retAddress, C argsAddress
	args := []Value{vm.get(instr.C)}
	return vm.call(instr.A, instr.B, args, false)
}

func exec_cso(instr *Instruction, vm *VM) int {
	// A funcIndex, B retAddress, C argsAddress
	args := []Value{vm.get(instr.C)}
	return vm.call(instr.A, instr.B, args, true)
}

func exec_ret(instr *Instruction, vm *VM) int {
	currentFrame := vm.callStack[vm.fp]

	// check if we are inside a try-finally
	if vm.returnFromFinally() {
		return vm_continue
	}

	// run finalizers for all functions except the global func
	// which is handled in the main run loop
	if vm.fp > 0 {
		vm.runFinalizables(currentFrame)
	}

	var retValue Value
	if instr.A != Void {
		retValue = vm.get(instr.A)
	}

	if vm.fp == 0 {
		// returning from main: exit
		vm.Error = io.EOF
		vm.RetValue = retValue
		return vm_exit
	}

	// pop one frame
	vm.callStack = vm.callStack[:vm.fp]
	vm.fp--

	// if the frame can be reused, clear its memory and add it to the cache.
	// This makes a huge impact in memory hungry programs.
	if !currentFrame.inClosure {
		currentFrame.finalizables = nil
		currentFrame.closures = nil
		for i := range currentFrame.values {
			currentFrame.values[i] = UndefinedValue
		}
		vm.frameCache = append(vm.frameCache, currentFrame)
	}

	prevFrame := vm.callStack[vm.fp]

	// set the return value
	if prevFrame.retAddress != Void {
		vm.set(prevFrame.retAddress, retValue)
	}

	if currentFrame.exit {
		vm.RetValue = retValue
		return vm_exit
	}

	return vm_continue
}

func exec_del(instr *Instruction, vm *VM) int {
	obj := vm.get(instr.A)
	if obj.Type != Map {
		return vm_next
	}

	property := vm.get(instr.B)

	m := obj.ToMap()
	m.Lock()
	delete(m.Map, property)
	m.Unlock()
	return vm_next
}

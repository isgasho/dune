// Code generated by "stringer -type=Opcode,AddressKind"; DO NOT EDIT.

package dune

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[op_ldk-0]
	_ = x[op_mov-1]
	_ = x[op_mob-2]
	_ = x[op_add-3]
	_ = x[op_sub-4]
	_ = x[op_mul-5]
	_ = x[op_div-6]
	_ = x[op_mod-7]
	_ = x[op_bor-8]
	_ = x[op_and-9]
	_ = x[op_xor-10]
	_ = x[op_lsh-11]
	_ = x[op_rsh-12]
	_ = x[op_inc-13]
	_ = x[op_dec-14]
	_ = x[op_unm-15]
	_ = x[op_not-16]
	_ = x[op_bnt-17]
	_ = x[op_str-18]
	_ = x[op_new-19]
	_ = x[op_nes-20]
	_ = x[op_arr-21]
	_ = x[op_map-22]
	_ = x[op_key-23]
	_ = x[op_val-24]
	_ = x[op_len-25]
	_ = x[op_enu-26]
	_ = x[op_get-27]
	_ = x[op_gto-28]
	_ = x[op_set-29]
	_ = x[op_spa-30]
	_ = x[op_jmp-31]
	_ = x[op_jpb-32]
	_ = x[op_ejp-33]
	_ = x[op_djp-34]
	_ = x[op_tjp-35]
	_ = x[op_eql-36]
	_ = x[op_neq-37]
	_ = x[op_seq-38]
	_ = x[op_sne-39]
	_ = x[op_lst-40]
	_ = x[op_lse-41]
	_ = x[op_cal-42]
	_ = x[op_cco-43]
	_ = x[op_cas-44]
	_ = x[op_cso-45]
	_ = x[op_rnp-46]
	_ = x[op_ret-47]
	_ = x[op_clo-48]
	_ = x[op_trw-49]
	_ = x[op_try-50]
	_ = x[op_tre-51]
	_ = x[op_cen-52]
	_ = x[op_fen-53]
	_ = x[op_trx-54]
	_ = x[op_del-55]
}

const _Opcode_name = "op_ldkop_movop_mobop_addop_subop_mulop_divop_modop_borop_andop_xorop_lshop_rshop_incop_decop_unmop_notop_bntop_strop_newop_nesop_arrop_mapop_keyop_valop_lenop_enuop_getop_gtoop_setop_spaop_jmpop_jpbop_ejpop_djpop_tjpop_eqlop_neqop_seqop_sneop_lstop_lseop_calop_ccoop_casop_csoop_rnpop_retop_cloop_trwop_tryop_treop_cenop_fenop_trxop_del"

var _Opcode_index = [...]uint16{0, 6, 12, 18, 24, 30, 36, 42, 48, 54, 60, 66, 72, 78, 84, 90, 96, 102, 108, 114, 120, 126, 132, 138, 144, 150, 156, 162, 168, 174, 180, 186, 192, 198, 204, 210, 216, 222, 228, 234, 240, 246, 252, 258, 264, 270, 276, 282, 288, 294, 300, 306, 312, 318, 324, 330, 336}

func (i Opcode) String() string {
	if i >= Opcode(len(_Opcode_index)-1) {
		return "Opcode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Opcode_name[_Opcode_index[i]:_Opcode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AddrVoid-0]
	_ = x[AddrLocal-1]
	_ = x[AddrGlobal-2]
	_ = x[AddrConstant-3]
	_ = x[AddrClosure-4]
	_ = x[AddrEnum-5]
	_ = x[AddrFunc-6]
	_ = x[AddrNativeFunc-7]
	_ = x[AddrClass-8]
	_ = x[AddrData-9]
	_ = x[AddrUnresolved-10]
}

const _AddressKind_name = "AddrVoidAddrLocalAddrGlobalAddrConstantAddrClosureAddrEnumAddrFuncAddrNativeFuncAddrClassAddrDataAddrUnresolved"

var _AddressKind_index = [...]uint8{0, 8, 17, 27, 39, 50, 58, 66, 80, 89, 97, 111}

func (i AddressKind) String() string {
	if i >= AddressKind(len(_AddressKind_index)-1) {
		return "AddressKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _AddressKind_name[_AddressKind_index[i]:_AddressKind_index[i+1]]
}

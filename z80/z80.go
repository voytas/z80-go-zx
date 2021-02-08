package z80

type word uint16

type CPU struct {
	// Program Counter - 16-bit address of the current instruction being fetched from memory
	PC   word
	IN   func(a, n byte) byte
	mem  Memory
	reg  *registers
	halt bool
	t    byte
}

func NewCPU(mem Memory) *CPU {
	cpu := &CPU{}
	cpu.mem = mem
	cpu.Reset()
	return cpu
}

func (cpu *CPU) readByte() byte {
	b := cpu.mem.read(cpu.PC)
	cpu.PC += 1
	return b
}

func (cpu *CPU) readWord() word {
	w := word(cpu.mem.read(cpu.PC)) + word(cpu.mem.read(cpu.PC+1))<<8
	cpu.PC += 2
	return w
}

func (cpu *CPU) wait() {}

func (cpu *CPU) Reset() {
	cpu.PC = 0
	cpu.reg = newRegisters()
	cpu.halt = false
}

func (cpu *CPU) Run() {
	var prefix byte
	for {
		opcode := cpu.readByte()

		if prefix == prefix_ix || prefix == prefix_iy {
			cpu.t = 4 // add extra 4 states for prefixed op
		} else {
			cpu.t = 0
		}

		switch opcode {
		case nop:
			cpu.t += 4
		case halt:
			cpu.t += 4
			cpu.wait()
			cpu.halt = true
			return
		case di:
			// TODO: Implement
		case ei:
			// TODO: Implement
		case cpl:
			cpu.reg.A = ^cpu.reg.A
			cpu.reg.F |= f_H | f_N
			cpu.t += 4
		case scf:
			cpu.reg.F &= ^(f_H | f_N)
			cpu.reg.F |= f_C
			cpu.t += 4
		case ccf:
			cpu.reg.F &= ^(f_N)
			cpu.reg.F ^= f_H | f_C
			cpu.t += 4
		case rlca:
			cpu.reg.F &= ^(f_H | f_N | f_C)
			a7 := cpu.reg.A >> 7
			cpu.reg.F |= a7
			cpu.reg.A = cpu.reg.A<<1 | a7
			cpu.t += 4
		case rrca:
			cpu.reg.F &= ^(f_H | f_N | f_C)
			a7 := cpu.reg.A & 0x01
			cpu.reg.F |= a7
			cpu.reg.A = cpu.reg.A>>1 | a7<<7
			cpu.t += 4
		case rla:
			fc := cpu.reg.F & f_C
			cpu.reg.F &= ^(f_H | f_N | f_C)
			a7 := cpu.reg.A >> 7
			cpu.reg.A = cpu.reg.A<<1 | fc
			cpu.reg.F |= a7
			cpu.t += 4
		case rra:
			fc := cpu.reg.F & f_C
			cpu.reg.F &= ^(f_H | f_N | f_C)
			a0 := cpu.reg.A & 0x01
			cpu.reg.A = cpu.reg.A>>1 | fc<<7
			cpu.reg.F |= a0
			cpu.t += 4
		case daa:
			cpu.reg.F &= ^(f_S | f_Z | f_P)
			a := cpu.reg.A
			if a&0xF > 9 || cpu.reg.F&f_H != 0 {
				if cpu.reg.F&f_N > 0 {
					a -= 0x06
				} else {
					a += 0x06
				}
			}
			if a > 0x99 || cpu.reg.F&f_C != 0 {
				if cpu.reg.F&f_N > 0 {
					a -= 0x60
				} else {
					a += 0x60
				}
			}
			cpu.reg.F |= a & f_S
			if a == 0 {
				cpu.reg.F |= f_Z
			}
			if cpu.reg.F&f_N > 0 {
				if cpu.reg.F&f_H > 0 && cpu.reg.A&0xF < 6 {
					cpu.reg.F |= f_H
				} else {
					cpu.reg.F &= ^f_H
				}
			} else {
				if cpu.reg.A&0xF > 9 {
					cpu.reg.F |= f_H
				} else {
					cpu.reg.F &= ^f_H
				}
			}
			cpu.reg.F |= parity[a]
			if cpu.reg.A > 0x99 {
				cpu.reg.F |= f_C
			}
			cpu.reg.A = a
			cpu.t += 4
		case ex_af_af:
			cpu.reg.A, cpu.reg.A_ = cpu.reg.A_, cpu.reg.A
			cpu.reg.F, cpu.reg.F_ = cpu.reg.F_, cpu.reg.F
			cpu.t += 4
		case exx:
			cpu.reg.B, cpu.reg.B_, cpu.reg.C, cpu.reg.C_ = cpu.reg.B_, cpu.reg.B, cpu.reg.C_, cpu.reg.C
			cpu.reg.D, cpu.reg.D_, cpu.reg.E, cpu.reg.E_ = cpu.reg.D_, cpu.reg.D, cpu.reg.E_, cpu.reg.E
			cpu.reg.H, cpu.reg.H_, cpu.reg.L, cpu.reg.L_ = cpu.reg.H_, cpu.reg.H, cpu.reg.L_, cpu.reg.L
			cpu.t += 4
		case ex_de_hl:
			// TODO: Implement
		case ex_sp_hl:
			// TODO: Implement
		case add_a_n, add_a_a, add_a_b, add_a_c, add_a_d, add_a_e, add_a_h, add_a_l, add_a_hl:
			var n byte
			if opcode == add_a_n {
				n = cpu.readByte()
				cpu.t += 7
			} else if opcode == add_a_hl {
				n = cpu.mem.read(cpu.getHL(prefix))
				cpu.t += 7
			} else {
				n = cpu.reg.getReg(opcode&0b00000111, prefix)
				cpu.t += 4
			}
			cpu.reg.F = f_NONE
			sum := cpu.reg.A + n
			if sum > 0x7F {
				cpu.reg.F |= f_S
			} else if sum == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= (cpu.reg.A ^ n ^ sum) & f_H
			if (cpu.reg.A^n)&0x80 == 0 && (cpu.reg.A^sum)&0x80 > 0 {
				cpu.reg.F |= f_P
			}
			if sum < cpu.reg.A {
				cpu.reg.F |= f_C
			}
			cpu.reg.A = sum
		case adc_a_n, adc_a_a, adc_a_b, adc_a_c, adc_a_d, adc_a_e, adc_a_h, adc_a_l, adc_a_hl:
			var n byte
			if opcode == adc_a_n {
				n = cpu.readByte()
				cpu.t = 7
			} else if opcode == adc_a_hl {
				n = cpu.mem.read(cpu.reg.getHL())
				cpu.t += 7
			} else {
				n = *cpu.reg.getR(opcode & 0b00000111)
				cpu.t += 4
			}
			cf := cpu.reg.F & f_C
			cpu.reg.F = f_NONE
			sum_w := word(cpu.reg.A) + word(n) + word(cf)
			sum_b := byte(sum_w)
			cpu.reg.F |= f_S & sum_b
			if sum_b == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= (cpu.reg.A ^ n ^ sum_b) & f_H
			if (cpu.reg.A^n)&0x80 == 0 && (cpu.reg.A^sum_b)&0x80 > 0 {
				cpu.reg.F |= f_P
			}
			if sum_w > 0xff {
				cpu.reg.F |= f_C
			}
			cpu.reg.A = sum_b
		case add_hl_bc, add_hl_de, add_hl_hl, add_hl_sp:
			hl := cpu.reg.getHL()
			var nn word
			switch opcode {
			case add_hl_bc:
				nn = cpu.reg.getBC()
			case add_hl_de:
				nn = cpu.reg.getDE()
			case add_hl_hl:
				nn = hl
			case add_hl_sp:
				nn = cpu.reg.SP
			}
			sum := hl + nn
			cpu.reg.setHL(sum)
			cpu.reg.F &= ^(f_H | f_N | f_C)
			if sum < hl {
				cpu.reg.F |= f_C
			}
			cpu.reg.F |= byte((hl^nn^sum)>>8) & f_H
			cpu.t += 11
		case sub_n, sub_a, sub_b, sub_c, sub_d, sub_e, sub_h, sub_l, sub_hl:
			a := cpu.reg.A
			var n byte
			if opcode == sub_n {
				n = cpu.readByte()
				cpu.t += 7
			} else if opcode == sub_hl {
				n = cpu.mem.read(cpu.reg.getHL())
				cpu.t += 7
			} else {
				n = *cpu.reg.getR(opcode & 0b00000111)
				cpu.t += 4
			}
			cpu.reg.A -= n
			cpu.reg.F = f_N
			cpu.reg.F |= f_S & cpu.reg.A
			if cpu.reg.A == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= byte(a^n^cpu.reg.A) & f_H
			if (a^n)&0x80 > 0 && (a^cpu.reg.A)&0x80 > 0 {
				cpu.reg.F |= f_P
			}
			if cpu.reg.A > a {
				cpu.reg.F |= f_C
			}
		case cp_a, cp_b, cp_c, cp_d, cp_e, cp_h, cp_l, cp_hl, cp_n:
			var n byte
			if opcode == cp_hl {
				n = cpu.mem.read(cpu.reg.getHL())
				cpu.t += 7
			} else if opcode == cp_n {
				// TODO: Implement
			} else {
				n = *cpu.reg.getR(opcode & 0b00000111)
				cpu.t += 4
			}
			test := cpu.reg.A - n
			cpu.reg.F = f_N
			cpu.reg.F |= f_S & test
			if test == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= byte(cpu.reg.A^n^test) & f_H
			if (cpu.reg.A^n)&0x80 > 0 && (cpu.reg.A^test)&0x80 > 0 {
				cpu.reg.F |= f_P
			}
			if test > cpu.reg.A {
				cpu.reg.F |= f_C
			}
		case sbc_a_a, sbc_a_b, sbc_a_c, sbc_a_d, sbc_a_e, sbc_a_h, sbc_a_l, sbc_a_hl, sbc_a_n:
			var n byte
			if opcode == sbc_a_hl {
				n = cpu.mem.read(cpu.reg.getHL())
				cpu.t += 7
			} else if opcode == sbc_a_n {
				// TODO: Implement
			} else {
				n = *cpu.reg.getR(opcode & 0b00000111)
				cpu.t += 4
			}
			cf := cpu.reg.F & f_C
			cpu.reg.F = f_N
			sub_w := word(cpu.reg.A) - word(n) - word(cf)
			sub_b := byte(sub_w)
			cpu.reg.F |= f_S & sub_b
			if sub_b == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= byte(cpu.reg.A^n^sub_b) & f_H
			if (cpu.reg.A^n)&0x80 > 0 && (sub_b^cpu.reg.A)&0x80 > 0 {
				cpu.reg.F |= f_P
			}
			if sub_w > 0xff {
				cpu.reg.F |= f_C
			}
			cpu.reg.A = sub_b
		case and_a, and_b, and_c, and_d, and_e, and_h, and_l, and_hl, and_n:
			var n byte
			if opcode == and_hl {
				n = cpu.mem.read(cpu.reg.getHL())
				cpu.t += 7
			} else if opcode == and_n {
				// TODO: Implement
			} else {
				n = *cpu.reg.getR(opcode & 0b00000111)
				cpu.t += 4
			}
			cpu.reg.F = f_H
			cpu.reg.A &= n
			cpu.reg.F |= f_S & cpu.reg.A
			if cpu.reg.A == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= parity[cpu.reg.A]
		case or_a, or_b, or_c, or_d, or_e, or_h, or_l, or_hl, or_n:
			var n byte
			if opcode == or_hl {
				n = cpu.mem.read(cpu.reg.getHL())
				cpu.t += 7
			} else if opcode == or_n {
				// TODO: Implement
			} else {
				n = *cpu.reg.getR(opcode & 0b00000111)
				cpu.t += 4
			}
			cpu.reg.F = f_NONE
			cpu.reg.A |= n
			cpu.reg.F |= f_S & cpu.reg.A
			if cpu.reg.A == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= parity[cpu.reg.A]
		case xor_a, xor_b, xor_c, xor_d, xor_e, xor_h, xor_l, xor_hl, xor_n:
			var n byte
			if opcode == xor_hl {
				n = cpu.mem.read(cpu.reg.getHL())
				cpu.t += 7
			} else if opcode == xor_n {
				// TODO: Implement
			} else {
				n = *cpu.reg.getR(opcode & 0b00000111)
				cpu.t += 4
			}
			cpu.reg.F = f_NONE
			cpu.reg.A ^= n
			cpu.reg.F |= f_S & cpu.reg.A
			if cpu.reg.A == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= parity[cpu.reg.A]
		case ld_a_n, ld_b_n, ld_c_n, ld_d_n, ld_e_n, ld_h_n, ld_l_n:
			r := cpu.reg.getR(opcode & 0b00111000 >> 3)
			*r = cpu.readByte()
			cpu.t += 7
		case
			ld_a_a, ld_a_b, ld_a_c, ld_a_d, ld_a_e, ld_a_h, ld_a_l,
			ld_b_a, ld_b_b, ld_b_c, ld_b_d, ld_b_e, ld_b_h, ld_b_l,
			ld_c_a, ld_c_b, ld_c_c, ld_c_d, ld_c_e, ld_c_h, ld_c_l,
			ld_d_a, ld_d_b, ld_d_c, ld_d_d, ld_d_e, ld_d_h, ld_d_l,
			ld_e_a, ld_e_b, ld_e_c, ld_e_d, ld_e_e, ld_e_h, ld_e_l,
			ld_h_a, ld_h_b, ld_h_c, ld_h_d, ld_h_e, ld_h_h, ld_h_l,
			ld_l_a, ld_l_b, ld_l_c, ld_l_d, ld_l_e, ld_l_h, ld_l_l:
			rs := cpu.reg.getR(opcode & 0b00000111)
			rd := cpu.reg.getR(opcode & 0b00111000 >> 3)
			*rd = *rs
			cpu.t += 4
		case ld_bc_nn:
			cpu.reg.C, cpu.reg.B = cpu.readByte(), cpu.readByte()
			cpu.t += 10
		case ld_de_nn:
			cpu.reg.E, cpu.reg.D = cpu.readByte(), cpu.readByte()
			cpu.t += 10
		case ld_hl_nn:
			switch prefix {
			case prefix_ix:
				cpu.reg.IX = cpu.readWord()
			case prefix_iy:
				cpu.reg.IY = cpu.readWord()
			default:
				cpu.reg.L, cpu.reg.H = cpu.readByte(), cpu.readByte()
			}
			cpu.t += 10
		case ld_sp_nn:
			cpu.reg.SP = cpu.readWord()
			cpu.t += 10
		case ld_sp_hl:
			// TODO: Implement
		case ld_hl_mm:
			w := cpu.readWord()
			cpu.reg.L = cpu.mem.read(w)
			cpu.reg.H = cpu.mem.read(w + 1)
			cpu.t += 16
		case ld_mm_hl:
			w := cpu.readWord()
			cpu.mem.write(w, cpu.reg.L)
			cpu.mem.write(w+1, cpu.reg.H)
			cpu.t += 16
		case ld_mhl_n:
			cpu.mem.write(cpu.getHL(prefix), cpu.readByte())
			// t is wrong for IX & IY
			cpu.t += 10
		case ld_mm_a:
			w := cpu.readWord()
			cpu.mem.write(w, cpu.reg.A)
			cpu.t += 13
		case ld_a_mm:
			w := cpu.readWord()
			cpu.reg.A = cpu.mem.read(w)
			cpu.t += 13
		case ld_bc_a:
			cpu.mem.write(cpu.reg.getBC(), cpu.reg.A)
			cpu.t += 7
		case ld_de_a:
			cpu.mem.write(cpu.reg.getDE(), cpu.reg.A)
			cpu.t += 7
		case ld_a_bc:
			cpu.reg.A = cpu.mem.read(cpu.reg.getBC())
			cpu.t += 7
		case ld_a_de:
			cpu.reg.A = cpu.mem.read(cpu.reg.getDE())
			cpu.t += 7
		case ld_a_hl, ld_b_hl, ld_c_hl, ld_d_hl, ld_e_hl, ld_h_hl, ld_l_hl:
			cpu.reg.setReg(opcode&0b00111000>>3, prefix_none, cpu.mem.read(cpu.getHL(prefix)))
			cpu.t += 7
		case ld_hl_a, ld_hl_b, ld_hl_c, ld_hl_d, ld_hl_e, ld_hl_h, ld_hl_l:
			cpu.mem.write(cpu.getHL(prefix), cpu.reg.getReg(opcode&0b00000111, prefix_none))
			cpu.t += 7
		case inc_a, inc_b, inc_c, inc_d, inc_e, inc_h, inc_l:
			r := cpu.reg.getR(opcode & 0b00111000 >> 3)
			cpu.reg.F &= ^(f_S | f_Z | f_H | f_P | f_N)
			if *r == 0x7F {
				cpu.reg.F |= f_P
			}
			if *r&0x0F == 0x0F {
				cpu.reg.F |= f_H
			}
			*r += 1
			if *r > 0x7F {
				cpu.reg.F |= f_S
			}
			if *r == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.t += 4
		case inc_bc:
			cpu.reg.setBC(cpu.reg.getBC() + 1)
			cpu.t += 6
		case inc_de:
			cpu.reg.setDE(cpu.reg.getDE() + 1)
			cpu.t += 6
		case inc_hl:
			cpu.reg.setHL(cpu.reg.getHL() + 1)
			cpu.t += 6
		case inc_sp:
			cpu.reg.SP += 1
			cpu.t += 6
		case inc_mhl:
			mm := cpu.reg.getHL()
			b := cpu.mem.read(mm)
			cpu.reg.F &= ^(f_S | f_Z | f_P | f_N)
			if b == 0x7F {
				cpu.reg.F |= f_P
			}
			if b&0x0F == 0x0F {
				cpu.reg.F |= f_H
			}
			b += 1
			if b == 0x00 {
				cpu.reg.F |= f_Z
			}
			if b > 0x7F {
				cpu.reg.F |= f_S
			}
			cpu.mem.write(mm, b)
			cpu.t += 11
		case dec_a, dec_b, dec_c, dec_d, dec_e, dec_h, dec_l:
			r := cpu.reg.getR(opcode & 0b00111000 >> 3)
			cpu.reg.F = cpu.reg.F & ^(f_S|f_Z|f_H|f_P) | f_N
			if *r == 0x80 {
				cpu.reg.F |= f_P
			}
			if *r&0x0F == 0 {
				cpu.reg.F |= f_H
			}
			*r -= 1
			if *r > 0x7F {
				cpu.reg.F |= f_S
			}
			if *r == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.t += 4
		case dec_bc:
			cpu.reg.setBC(cpu.reg.getBC() - 1)
			cpu.t += 6
		case dec_de:
			cpu.reg.setDE(cpu.reg.getDE() - 1)
			cpu.t += 6
		case dec_hl:
			cpu.reg.setHL(cpu.reg.getHL() - 1)
			cpu.t += 6
		case dec_sp:
			cpu.reg.SP -= 1
			cpu.t += 6
		case dec_mhl:
			mm := cpu.reg.getHL()
			b := cpu.mem.read(mm)
			cpu.reg.F &= ^(f_S | f_Z | f_P)
			cpu.reg.F |= f_N
			if b == 0x80 {
				cpu.reg.F |= f_P
			}
			if b&0x0F == 0 {
				cpu.reg.F |= f_H
			}
			b -= 1
			if b == 0x00 {
				cpu.reg.F |= f_Z
			}
			if b > 0x7F {
				cpu.reg.F |= f_S
			}
			cpu.mem.write(mm, b)
			cpu.t += 11
		case jr_o:
			o := cpu.readByte()
			if o&0x80 == 0 {
				cpu.PC += word(o)
			} else {
				cpu.PC -= word(^o + 1)
			}
			cpu.t += 12
		case jr_z_o:
			o := cpu.readByte()
			if cpu.reg.F&f_Z == f_Z {
				if o&0x80 == 0 {
					cpu.PC += word(o)
				} else {
					cpu.PC -= word(^o + 1)
				}
				cpu.t += 12
			} else {
				cpu.t += 7
			}
		case jr_nz_o:
			o := cpu.readByte()
			if cpu.reg.F&f_Z == 0 {
				if o&0x80 == 0 {
					cpu.PC += word(o)
				} else {
					cpu.PC -= word(^o + 1)
				}
				cpu.t += 12
			} else {
				cpu.t += 7
			}
		case jr_c:
			o := cpu.readByte()
			if cpu.reg.F&f_C == f_C {
				if o&0x80 == 0 {
					cpu.PC += word(o)
				} else {
					cpu.PC -= word(^o + 1)
				}
				cpu.t += 12
			} else {
				cpu.t += 7
			}
		case jr_nc_o:
			o := cpu.readByte()
			if cpu.reg.F&f_C == 0 {
				if o&0x80 == 0 {
					cpu.PC += word(o)
				} else {
					cpu.PC -= word(^o + 1)
				}
				cpu.t += 12
			} else {
				cpu.t += 7
			}
		case djnz:
			o := cpu.readByte()
			cpu.reg.B -= 1
			if cpu.reg.B == 0 {
				cpu.t += 8
			} else {
				if o&0x80 == 0 {
					cpu.PC += word(o)
				} else {
					cpu.PC -= word(^o + 1)
				}
				cpu.t += 13
			}
		case jp_nn:
			cpu.PC = cpu.readWord()
			cpu.t += 10
		case jp_c_nn, jp_m_nn, jp_nc_nn, jp_nz_nn, jp_p_nn, jp_pe_nn, jp_po_nn, jp_z_nn:
			if cpu.shouldJump(opcode) {
				cpu.PC = cpu.readWord()
			}
			cpu.t += 10
		case jp_hl:
			// TODO: Implement
		case call_nn, call_c_nn, call_m_nn, call_nc_nn, call_nz_nn, call_p_nn, call_pe_nn, call_po_nn, call_z_nn:
			if cpu.shouldJump(opcode) {
				pc := cpu.readWord()
				cpu.reg.SP -= 1
				cpu.mem.write(cpu.reg.SP, byte(cpu.PC>>8))
				cpu.reg.SP -= 1
				cpu.mem.write(cpu.reg.SP, byte(cpu.PC))
				cpu.PC = pc
				cpu.t += 17
			} else {
				cpu.PC += 2
				cpu.t += 10
			}
		case ret:
			cpu.PC = word(cpu.mem.read(cpu.reg.SP+1))<<8 | word(cpu.mem.read(cpu.reg.SP))
			cpu.reg.SP += 2
			cpu.t += 10
		case ret_c, ret_m, ret_nc, ret_nz, ret_p, ret_pe, ret_po, ret_z:
			if cpu.shouldJump(opcode) {
				cpu.PC = word(cpu.mem.read(cpu.reg.SP+1))<<8 | word(cpu.mem.read(cpu.reg.SP))
				cpu.reg.SP += 2
				cpu.t += 11
			} else {
				cpu.t += 5
			}
		case rst_00h, rst_08h, rst_10h, rst_18h, rst_20h, rst_28h, rst_30h, rst_38h:
			cpu.reg.SP -= 1
			cpu.mem.write(cpu.reg.SP, byte(cpu.PC>>8))
			cpu.reg.SP -= 1
			cpu.mem.write(cpu.reg.SP, byte(cpu.PC))
			cpu.PC = word(8 * ((opcode & 0b00111000) >> 3))
			cpu.t += 11
		case push_af:
			cpu.reg.SP -= 1
			cpu.mem.write(cpu.reg.SP, cpu.reg.A)
			cpu.reg.SP -= 1
			cpu.mem.write(cpu.reg.SP, cpu.reg.F)
			cpu.t += 11
		case push_bc:
			cpu.reg.SP -= 1
			cpu.mem.write(cpu.reg.SP, cpu.reg.B)
			cpu.reg.SP -= 1
			cpu.mem.write(cpu.reg.SP, cpu.reg.C)
			cpu.t += 11
		case push_de:
			cpu.reg.SP -= 1
			cpu.mem.write(cpu.reg.SP, cpu.reg.D)
			cpu.reg.SP -= 1
			cpu.mem.write(cpu.reg.SP, cpu.reg.E)
			cpu.t += 11
		case push_hl:
			cpu.reg.SP -= 1
			cpu.mem.write(cpu.reg.SP, cpu.reg.H)
			cpu.reg.SP -= 1
			cpu.mem.write(cpu.reg.SP, cpu.reg.L)
			cpu.t += 11
		case pop_af:
			cpu.reg.A, cpu.reg.F = cpu.mem.read(cpu.reg.SP+1), cpu.mem.read(cpu.reg.SP)
			cpu.reg.SP += 2
			cpu.t += 10
		case pop_bc:
			cpu.reg.B, cpu.reg.C = cpu.mem.read(cpu.reg.SP+1), cpu.mem.read(cpu.reg.SP)
			cpu.reg.SP += 2
			cpu.t += 10
		case pop_de:
			cpu.reg.D, cpu.reg.E = cpu.mem.read(cpu.reg.SP+1), cpu.mem.read(cpu.reg.SP)
			cpu.reg.SP += 2
			cpu.t += 10
		case pop_hl:
			cpu.reg.H, cpu.reg.L = cpu.mem.read(cpu.reg.SP+1), cpu.mem.read(cpu.reg.SP)
			cpu.reg.SP += 2
			cpu.t += 10
		case in_a_n:
			n := cpu.readByte()
			if cpu.IN != nil {
				cpu.reg.A = cpu.IN(cpu.reg.A, n)
			}
			cpu.t += 11
		case out_n_a:
			// TODO:
		case prefix_cb:
			cpu.prefix_cb(cpu.readByte())
		case prefix_ix:
			if prefix == prefix_ix || prefix == prefix_iy || prefix == prefix_ED {
				// NOP
				cpu.t += 4
			} else {
				prefix = prefix_ix
				continue
			}
		case prefix_iy:
			if prefix == prefix_ix || prefix == prefix_iy || prefix == prefix_ED {
				// NOP
				cpu.t += 4
			} else {
				prefix = prefix_iy
				continue
			}
		}

		prefix = prefix_none // reset ix or iy prefix if
		cpu.wait()
	}
}

func (cpu *CPU) shouldJump(opcode byte) bool {
	if opcode == call_nn {
		return true
	}

	switch opcode & 0b00111000 {
	case 0b00000000: // Non-Zero (NZ)
		return cpu.reg.F&f_Z == 0
	case 0b00001000: // Zero (Z)
		return cpu.reg.F&f_Z != 0
	case 0b00010000: // Non Carry (NC)
		return cpu.reg.F&f_C == 0
	case 0b00011000: // Carry (C)
		return cpu.reg.F&f_C != 0
	case 0b00100000: // Parity Odd (PO)
		return cpu.reg.F&f_P == 0
	case 0b00101000: // Parity Even (PE)
		return cpu.reg.F&f_P != 0
	case 0b00110000: // Sign Positive (P)
		return cpu.reg.F&f_S == 0
	case 0b00111000: // Sign Negative (M)
		return cpu.reg.F&f_S != 0
	}

	panic("Invalid opcode")
}

func (cpu *CPU) prefix_cb(opcode byte) {
	var v byte
	var hl word

	reg := opcode & 0b00000111
	if reg == r_HL {
		hl = cpu.reg.getHL()
		v = cpu.mem.read(hl)
		cpu.t += 15
	} else {
		v = *cpu.reg.regs8[reg]
		cpu.t += 8
	}

	var cy byte
	write := func() {
		if reg == r_HL {
			cpu.mem.write(hl, v)
		} else {
			*cpu.reg.regs8[reg] = v
		}
		cpu.reg.F = f_NONE
		cpu.reg.F |= f_S & v
		if v == 0 {
			cpu.reg.F |= f_Z
		}
		cpu.reg.F |= parity[v] | cy
	}

	switch opcode & 0b11111000 {
	case rlc_r:
		cy = v >> 7
		v = v<<1 | cy
		write()
	case rrc_r:
		cy = v & f_C
		v = v>>1 | cy<<7
		write()
	case rl_r:
		cy = v >> 7
		v = v<<1 | f_C&cpu.reg.F
		write()
	case rr_r:
		cy = v & f_C
		v = v>>1 | f_C&cpu.reg.F<<7
		write()
	case sla_r:
		cy = v >> 7
		v = v << 1
		write()
	case sra_r:
		cy = v & f_C
		v = v&0x80 | v>>1
		write()
	case sll_r:
		cy = v >> 7
		v = v<<1 | 0x01
		write()
	case srl_r:
		cy = v & f_C
		v = v >> 1
		write()
	default:
		bit := (opcode & 0b00111000) >> 3
		switch opcode & 0b11000000 {
		case bit_b:
			if reg == r_HL {
				cpu.t += 12
			}
			cpu.reg.F &= ^(f_Z | f_N)
			cpu.reg.F |= f_H
			if v&bit_mask[bit] == 0 {
				cpu.reg.F |= f_Z
			}
		case res_b:
			v &= ^bit_mask[bit]
			write()
		case set_b:
			v |= bit_mask[bit]
			write()
		}
	}
}

// Returns value of HL / (IX + d) / (IY + d) register. The prefix parameter specifies
// whether to use IX or IY register instead of HL. t-states are updated with extra
// cycles in case of IX and IY.
func (cpu *CPU) getHL(prefix byte) word {
	var hl word
	switch prefix {
	case prefix_ix:
		o := cpu.readByte()
		if o&0x80 == 0 {
			hl = cpu.reg.IX + word(o)
		} else {
			hl = cpu.reg.IX - word(^o+1)
		}
		cpu.t += 8
	case prefix_iy:
		o := cpu.readByte()
		if o&0x80 == 0 {
			hl = cpu.reg.IY + word(o)
		} else {
			hl = cpu.reg.IY - word(^o+1)
		}
		cpu.t += 8
	default:
		hl = cpu.reg.getHL()
	}

	return hl
}

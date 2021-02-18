package z80

import (
	"fmt"

	"github.com/voytas/z80-go-zx/z80/memory"
)

type CPU struct {
	PC               uint16
	IN               func(a, n byte) byte
	OUT              func(a, n byte)
	mem              memory.Memory
	reg              *registers
	t                byte
	halt, iff1, iff2 bool
}

func NewCPU(mem memory.Memory) *CPU {
	cpu := &CPU{}
	cpu.mem = mem
	cpu.Reset()
	return cpu
}

func (cpu *CPU) readByte() byte {
	b := cpu.mem.Read(cpu.PC)
	cpu.PC += 1
	return b
}

func (cpu *CPU) readWord() uint16 {
	w := uint16(cpu.mem.Read(cpu.PC)) | uint16(cpu.mem.Read(cpu.PC+1))<<8
	cpu.PC += 2
	return w
}

func (cpu *CPU) wait() {
}

func (cpu *CPU) Reset() {
	cpu.PC = 0
	cpu.reg = newRegisters()
	cpu.halt = false
	cpu.iff1, cpu.iff2 = true, true
}

func (cpu *CPU) Run() {
	for {
		opcode := cpu.readByte()

		//cpu.debug(opcode)

		// Get the t-state for the current instruction
		if cpu.reg.prefix == useIX || cpu.reg.prefix == useIY {
			t := t_states_ixy[opcode]
			if t != 0 {
				cpu.t = t
			} else {
				cpu.t = 4 + t_states[opcode]
			}
		} else {
			cpu.t = t_states[opcode]
		}

		switch opcode {
		case nop:
		case halt:
			cpu.reg.prefix = noPrefix
			cpu.wait()
			cpu.halt = true
			return
		case di:
			cpu.iff1, cpu.iff2 = false, false
		case ei:
			cpu.iff1, cpu.iff2 = true, true
		case cpl:
			cpu.reg.A = ^cpu.reg.A
			cpu.reg.F |= f_H | f_N
		case scf:
			cpu.reg.F &= ^(f_H | f_N)
			cpu.reg.F |= f_C
		case ccf:
			cpu.reg.F &= ^(f_N)
			cpu.reg.F ^= f_H | f_C
		case rlca:
			cpu.reg.F &= ^(f_H | f_N | f_C)
			a7 := cpu.reg.A >> 7
			cpu.reg.F |= a7
			cpu.reg.A = cpu.reg.A<<1 | a7
		case rrca:
			cpu.reg.F &= ^(f_H | f_N | f_C)
			a7 := cpu.reg.A & 0x01
			cpu.reg.F |= a7
			cpu.reg.A = cpu.reg.A>>1 | a7<<7
		case rla:
			fc := cpu.reg.F & f_C
			cpu.reg.F &= ^(f_H | f_N | f_C)
			a7 := cpu.reg.A >> 7
			cpu.reg.A = cpu.reg.A<<1 | fc
			cpu.reg.F |= a7
		case rra:
			fc := cpu.reg.F & f_C
			cpu.reg.F &= ^(f_H | f_N | f_C)
			a0 := cpu.reg.A & 0x01
			cpu.reg.A = cpu.reg.A>>1 | fc<<7
			cpu.reg.F |= a0
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
		case ex_af_af:
			cpu.reg.A, cpu.reg.A_ = cpu.reg.A_, cpu.reg.A
			cpu.reg.F, cpu.reg.F_ = cpu.reg.F_, cpu.reg.F
		case exx:
			cpu.reg.B, cpu.reg.B_, cpu.reg.C, cpu.reg.C_ = cpu.reg.B_, cpu.reg.B, cpu.reg.C_, cpu.reg.C
			cpu.reg.D, cpu.reg.D_, cpu.reg.E, cpu.reg.E_ = cpu.reg.D_, cpu.reg.D, cpu.reg.E_, cpu.reg.E
			cpu.reg.H, cpu.reg.H_, cpu.reg.L, cpu.reg.L_ = cpu.reg.H_, cpu.reg.H, cpu.reg.L_, cpu.reg.L
		case ex_de_hl:
			cpu.reg.D, cpu.reg.E, cpu.reg.H, cpu.reg.L = cpu.reg.H, cpu.reg.L, cpu.reg.D, cpu.reg.E
		case ex_sp_hl:
			h, l := *cpu.reg.getReg(r_H), *cpu.reg.getReg(r_L)
			cpu.reg.setHLb(cpu.mem.Read(cpu.reg.SP+1), cpu.mem.Read(cpu.reg.SP))
			cpu.mem.Write(cpu.reg.SP, l)
			cpu.mem.Write(cpu.reg.SP+1, h)
		case add_a_n, add_a_a, add_a_b, add_a_c, add_a_d, add_a_e, add_a_h, add_a_l, add_a_hl:
			a := cpu.reg.A
			var n byte
			switch opcode {
			case add_a_n:
				n = cpu.readByte()
			case add_a_hl:
				n = cpu.mem.Read(cpu.getHL(true))
			default:
				n = *cpu.reg.getReg(opcode & 0b00000111)
			}
			cpu.reg.A += n
			cpu.reg.F = f_NONE
			if cpu.reg.A > 0x7F {
				cpu.reg.F |= f_S
			} else if cpu.reg.A == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= (a ^ n ^ cpu.reg.A) & f_H
			if (a^n)&0x80 == 0 && (a^cpu.reg.A)&0x80 != 0 {
				cpu.reg.F |= f_P
			}
			if cpu.reg.A < a {
				cpu.reg.F |= f_C
			}
		case adc_a_n, adc_a_a, adc_a_b, adc_a_c, adc_a_d, adc_a_e, adc_a_h, adc_a_l, adc_a_hl:
			var n byte
			switch opcode {
			case adc_a_n:
				n = cpu.readByte()
			case adc_a_hl:
				n = cpu.mem.Read(cpu.getHL(true))
			default:
				n = *cpu.reg.getReg(opcode & 0b00000111)
			}
			cf := cpu.reg.F & f_C
			cpu.reg.F = f_NONE
			sum_w := uint16(cpu.reg.A) + uint16(n) + uint16(cf)
			sum_b := byte(sum_w)
			cpu.reg.F |= f_S & sum_b
			if sum_b == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= (cpu.reg.A ^ n ^ sum_b) & f_H
			if (cpu.reg.A^n)&0x80 == 0 && (cpu.reg.A^sum_b)&0x80 != 0 {
				cpu.reg.F |= f_P
			}
			if sum_w > 0xff {
				cpu.reg.F |= f_C
			}
			cpu.reg.A = sum_b
		case add_hl_bc, add_hl_de, add_hl_hl, add_hl_sp:
			hl := cpu.getHL(false)
			var nn uint16
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
			cpu.reg.setHLw(sum)
			cpu.reg.F &= ^(f_H | f_N | f_C)
			if sum < hl {
				cpu.reg.F |= f_C
			}
			cpu.reg.F |= byte((hl^nn^sum)>>8) & f_H
		case sub_a, sub_b, sub_c, sub_d, sub_e, sub_h, sub_l, sub_hl, sub_n:
			a := cpu.reg.A
			var n byte
			switch opcode {
			case sub_n:
				n = cpu.readByte()
			case sub_hl:
				n = cpu.mem.Read(cpu.getHL(true))
			default:
				n = *cpu.reg.getReg(opcode & 0b00000111)
			}
			cpu.reg.A -= n
			cpu.reg.F = f_N
			cpu.reg.F |= f_S & cpu.reg.A
			if cpu.reg.A == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= (a ^ n ^ cpu.reg.A) & f_H
			if (a^n)&0x80 != 0 && (a^cpu.reg.A)&0x80 != 0 {
				cpu.reg.F |= f_P
			}
			if cpu.reg.A > a {
				cpu.reg.F |= f_C
			}
		case cp_a, cp_b, cp_c, cp_d, cp_e, cp_h, cp_l, cp_hl, cp_n:
			var n byte
			switch opcode {
			case cp_n:
				n = cpu.readByte()
			case cp_hl:
				n = cpu.mem.Read(cpu.getHL(true))
			default:
				n = *cpu.reg.getReg(opcode & 0b00000111)
			}
			test := cpu.reg.A - n
			cpu.reg.F = f_N
			cpu.reg.F |= f_S & test
			if test == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= byte(cpu.reg.A^n^test) & f_H
			if (cpu.reg.A^n)&0x80 != 0 && (cpu.reg.A^test)&0x80 != 0 {
				cpu.reg.F |= f_P
			}
			if test > cpu.reg.A {
				cpu.reg.F |= f_C
			}
		case sbc_a_a, sbc_a_b, sbc_a_c, sbc_a_d, sbc_a_e, sbc_a_h, sbc_a_l, sbc_a_hl, sbc_a_n:
			var n byte
			switch opcode {
			case sbc_a_n:
				n = cpu.readByte()
			case sbc_a_hl:
				n = cpu.mem.Read(cpu.getHL(true))
			default:
				n = *cpu.reg.getReg(opcode & 0b00000111)
			}
			cf := cpu.reg.F & f_C
			cpu.reg.F = f_N
			sub_w := uint16(cpu.reg.A) - uint16(n) - uint16(cf)
			sub_b := byte(sub_w)
			cpu.reg.F |= f_S & sub_b
			if sub_b == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= byte(cpu.reg.A^n^sub_b) & f_H
			if (cpu.reg.A^n)&0x80 != 0 && (sub_b^cpu.reg.A)&0x80 != 0 {
				cpu.reg.F |= f_P
			}
			if sub_w > 0xff {
				cpu.reg.F |= f_C
			}
			cpu.reg.A = sub_b
		case and_a, and_b, and_c, and_d, and_e, and_h, and_l, and_hl, and_n:
			var n byte
			switch opcode {
			case and_n:
				n = cpu.readByte()
			case and_hl:
				n = cpu.mem.Read(cpu.getHL(true))
			default:
				n = *cpu.reg.getReg(opcode & 0b00000111)
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
			switch opcode {
			case or_n:
				n = cpu.readByte()
			case or_hl:
				n = cpu.mem.Read(cpu.getHL(true))
			default:
				n = *cpu.reg.getReg(opcode & 0b00000111)
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
			switch opcode {
			case xor_n:
				n = cpu.readByte()
			case xor_hl:
				n = cpu.mem.Read(cpu.getHL(true))
			default:
				n = *cpu.reg.getReg(opcode & 0b00000111)
			}
			cpu.reg.F = f_NONE
			cpu.reg.A ^= n
			cpu.reg.F |= f_S & cpu.reg.A
			if cpu.reg.A == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= parity[cpu.reg.A]
		case ld_a_n, ld_b_n, ld_c_n, ld_d_n, ld_e_n, ld_h_n, ld_l_n:
			r := cpu.reg.getReg(opcode & 0b00111000 >> 3)
			*r = cpu.readByte()
		case
			ld_a_a, ld_a_b, ld_a_c, ld_a_d, ld_a_e, ld_a_h, ld_a_l,
			ld_b_a, ld_b_b, ld_b_c, ld_b_d, ld_b_e, ld_b_h, ld_b_l,
			ld_c_a, ld_c_b, ld_c_c, ld_c_d, ld_c_e, ld_c_h, ld_c_l,
			ld_d_a, ld_d_b, ld_d_c, ld_d_d, ld_d_e, ld_d_h, ld_d_l,
			ld_e_a, ld_e_b, ld_e_c, ld_e_d, ld_e_e, ld_e_h, ld_e_l,
			ld_h_a, ld_h_b, ld_h_c, ld_h_d, ld_h_e, ld_h_h, ld_h_l,
			ld_l_a, ld_l_b, ld_l_c, ld_l_d, ld_l_e, ld_l_h, ld_l_l:
			rs := cpu.reg.getReg(opcode & 0b00000111)
			rd := cpu.reg.getReg(opcode & 0b00111000 >> 3)
			*rd = *rs
		case ld_bc_nn:
			cpu.reg.C, cpu.reg.B = cpu.readByte(), cpu.readByte()
		case ld_de_nn:
			cpu.reg.E, cpu.reg.D = cpu.readByte(), cpu.readByte()
		case ld_hl_nn:
			h, l := cpu.reg.getReg(r_H), cpu.reg.getReg(r_L)
			*l, *h = cpu.readByte(), cpu.readByte()
		case ld_sp_nn:
			cpu.reg.SP = cpu.readWord()
		case ld_sp_hl:
			cpu.reg.SP = cpu.reg.getHL()
		case ld_hl_mm:
			addr := cpu.readWord()
			h, l := cpu.reg.getReg(r_H), cpu.reg.getReg(r_L)
			*l = cpu.mem.Read(addr)
			*h = cpu.mem.Read(addr + 1)
		case ld_mm_hl:
			w := cpu.readWord()
			h, l := cpu.reg.getReg(r_H), cpu.reg.getReg(r_L)
			cpu.mem.Write(w, *l)
			cpu.mem.Write(w+1, *h)
		case ld_mhl_n:
			cpu.mem.Write(cpu.getHL(true), cpu.readByte())
		case ld_mm_a:
			w := cpu.readWord()
			cpu.mem.Write(w, cpu.reg.A)
		case ld_a_mm:
			w := cpu.readWord()
			cpu.reg.A = cpu.mem.Read(w)
		case ld_bc_a:
			cpu.mem.Write(cpu.reg.getBC(), cpu.reg.A)
		case ld_de_a:
			cpu.mem.Write(cpu.reg.getDE(), cpu.reg.A)
		case ld_a_bc:
			cpu.reg.A = cpu.mem.Read(cpu.reg.getBC())
		case ld_a_de:
			cpu.reg.A = cpu.mem.Read(cpu.reg.getDE())
		case ld_a_hl, ld_b_hl, ld_c_hl, ld_d_hl, ld_e_hl, ld_h_hl, ld_l_hl:
			*cpu.reg.get[opcode&0b00111000>>3] = cpu.mem.Read(cpu.getHL(true))
		case ld_hl_a, ld_hl_b, ld_hl_c, ld_hl_d, ld_hl_e, ld_hl_h, ld_hl_l:
			cpu.mem.Write(cpu.getHL(true), *cpu.reg.get[opcode&0b00000111])
		case inc_a, inc_b, inc_c, inc_d, inc_e, inc_h, inc_l:
			r := cpu.reg.getReg(opcode & 0b00111000 >> 3)
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
		case inc_bc:
			cpu.reg.setBC(cpu.reg.getBC() + 1)
		case inc_de:
			cpu.reg.setDE(cpu.reg.getDE() + 1)
		case inc_hl:
			cpu.reg.setHLw(cpu.reg.getHL() + 1)
		case inc_sp:
			cpu.reg.SP += 1
		case inc_mhl:
			mm := cpu.getHL(true)
			b := cpu.mem.Read(mm)
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
			cpu.mem.Write(mm, b)
		case dec_a, dec_b, dec_c, dec_d, dec_e, dec_h, dec_l:
			r := cpu.reg.getReg(opcode & 0b00111000 >> 3)
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
		case dec_bc:
			cpu.reg.setBC(cpu.reg.getBC() - 1)
		case dec_de:
			cpu.reg.setDE(cpu.reg.getDE() - 1)
		case dec_hl:
			cpu.reg.setHLw(cpu.reg.getHL() - 1)
		case dec_sp:
			cpu.reg.SP -= 1
		case dec_mhl:
			mm := cpu.getHL(true)
			b := cpu.mem.Read(mm)
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
			cpu.mem.Write(mm, b)
		case jr_o:
			o := cpu.readByte()
			if o&0x80 == 0 {
				cpu.PC += uint16(o)
			} else {
				cpu.PC -= uint16(^o + 1)
			}
		case jr_z_o:
			o := cpu.readByte()
			if cpu.reg.F&f_Z == f_Z {
				if o&0x80 == 0 {
					cpu.PC += uint16(o)
				} else {
					cpu.PC -= uint16(^o + 1)
				}
				cpu.t += 5
			}
		case jr_nz_o:
			o := cpu.readByte()
			if cpu.reg.F&f_Z == 0 {
				if o&0x80 == 0 {
					cpu.PC += uint16(o)
				} else {
					cpu.PC -= uint16(^o + 1)
				}
				cpu.t += 5
			}
		case jr_c:
			o := cpu.readByte()
			if cpu.reg.F&f_C == f_C {
				if o&0x80 == 0 {
					cpu.PC += uint16(o)
				} else {
					cpu.PC -= uint16(^o + 1)
				}
				cpu.t += 5
			}
		case jr_nc_o:
			o := cpu.readByte()
			if cpu.reg.F&f_C == 0 {
				if o&0x80 == 0 {
					cpu.PC += uint16(o)
				} else {
					cpu.PC -= uint16(^o + 1)
				}
				cpu.t += 5
			}
		case djnz:
			o := cpu.readByte()
			cpu.reg.B -= 1
			if cpu.reg.B != 0 {
				if o&0x80 == 0 {
					cpu.PC += uint16(o)
				} else {
					cpu.PC -= uint16(^o + 1)
				}
				cpu.t += 5
			}
		case jp_nn:
			cpu.PC = cpu.readWord()
		case jp_c_nn, jp_m_nn, jp_nc_nn, jp_nz_nn, jp_p_nn, jp_pe_nn, jp_po_nn, jp_z_nn:
			if cpu.shouldJump(opcode) {
				cpu.PC = cpu.readWord()
			} else {
				cpu.PC += 2
			}
		case jp_hl:
			cpu.PC = cpu.reg.getHL()
		case call_nn, call_c_nn, call_m_nn, call_nc_nn, call_nz_nn, call_p_nn, call_pe_nn, call_po_nn, call_z_nn:
			if cpu.shouldJump(opcode) {
				pc := cpu.readWord()
				cpu.reg.SP -= 1
				cpu.mem.Write(cpu.reg.SP, byte(cpu.PC>>8))
				cpu.reg.SP -= 1
				cpu.mem.Write(cpu.reg.SP, byte(cpu.PC))
				cpu.PC = pc
				cpu.t += 7
			} else {
				cpu.PC += 2
			}
		case ret:
			cpu.PC = uint16(cpu.mem.Read(cpu.reg.SP+1))<<8 | uint16(cpu.mem.Read(cpu.reg.SP))
			cpu.reg.SP += 2
		case ret_c, ret_m, ret_nc, ret_nz, ret_p, ret_pe, ret_po, ret_z:
			if cpu.shouldJump(opcode) {
				cpu.PC = uint16(cpu.mem.Read(cpu.reg.SP+1))<<8 | uint16(cpu.mem.Read(cpu.reg.SP))
				cpu.reg.SP += 2
				cpu.t += 6
			}
		case rst_00h, rst_08h, rst_10h, rst_18h, rst_20h, rst_28h, rst_30h, rst_38h:
			cpu.reg.SP -= 1
			cpu.mem.Write(cpu.reg.SP, byte(cpu.PC>>8))
			cpu.reg.SP -= 1
			cpu.mem.Write(cpu.reg.SP, byte(cpu.PC))
			cpu.PC = uint16(8 * ((opcode & 0b00111000) >> 3))
		case push_af:
			cpu.reg.SP -= 1
			cpu.mem.Write(cpu.reg.SP, cpu.reg.A)
			cpu.reg.SP -= 1
			cpu.mem.Write(cpu.reg.SP, cpu.reg.F)
		case push_bc:
			cpu.reg.SP -= 1
			cpu.mem.Write(cpu.reg.SP, cpu.reg.B)
			cpu.reg.SP -= 1
			cpu.mem.Write(cpu.reg.SP, cpu.reg.C)
		case push_de:
			cpu.reg.SP -= 1
			cpu.mem.Write(cpu.reg.SP, cpu.reg.D)
			cpu.reg.SP -= 1
			cpu.mem.Write(cpu.reg.SP, cpu.reg.E)
		case push_hl:
			cpu.reg.SP -= 1
			cpu.mem.Write(cpu.reg.SP, *cpu.reg.getReg(r_H))
			cpu.reg.SP -= 1
			cpu.mem.Write(cpu.reg.SP, *cpu.reg.getReg(r_L))
		case pop_af:
			cpu.reg.A, cpu.reg.F = cpu.mem.Read(cpu.reg.SP+1), cpu.mem.Read(cpu.reg.SP)
			cpu.reg.SP += 2
		case pop_bc:
			cpu.reg.B, cpu.reg.C = cpu.mem.Read(cpu.reg.SP+1), cpu.mem.Read(cpu.reg.SP)
			cpu.reg.SP += 2
		case pop_de:
			cpu.reg.D, cpu.reg.E = cpu.mem.Read(cpu.reg.SP+1), cpu.mem.Read(cpu.reg.SP)
			cpu.reg.SP += 2
		case pop_hl:
			*cpu.reg.getReg(r_H), *cpu.reg.getReg(r_L) = cpu.mem.Read(cpu.reg.SP+1), cpu.mem.Read(cpu.reg.SP)
			cpu.reg.SP += 2
		case in_a_n:
			n := cpu.readByte()
			if cpu.IN != nil {
				cpu.reg.A = cpu.IN(cpu.reg.A, n)
			}
		case out_n_a:
			n := cpu.readByte()
			if cpu.OUT != nil {
				cpu.OUT(cpu.reg.A, n)
			}
		case prefix_cb:
			cpu.prefixCB(cpu.readByte())
		case prefix_ed:
			cpu.prefixED(cpu.readByte())
		case useIX:
			if cpu.reg.prefix != noPrefix {
				cpu.wait()
			}
			cpu.reg.prefix = useIX
			continue
		case useIY:
			if cpu.reg.prefix != noPrefix {
				cpu.wait()
			}
			cpu.reg.prefix = useIY
			continue
		}

		cpu.reg.prefix = noPrefix
		cpu.wait()
	}
}

// Handles opcodes with CB prefix
func (cpu *CPU) prefixCB(opcode byte) {
	var v byte
	var hl uint16

	reg := opcode & 0b00000111
	if reg == 0b110 {
		hl = cpu.reg.getHL()
		v = cpu.mem.Read(hl)
		cpu.t += 15 // the only exception is bit operation that takes 12 t-states
	} else {
		v = *cpu.reg.get[reg]
		cpu.t += 8
	}

	var cy byte
	write := func() {
		if reg == 0b110 {
			cpu.mem.Write(hl, v)
		} else {
			*cpu.reg.get[reg] = v
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
			if reg == 0b110 {
				cpu.t -= 3 // for bit operation it is 12 t-states, not usual 15
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

func (cpu *CPU) prefixED(opcode byte) {
	t := t_states[opcode]
	if t != 0 {
		cpu.t += t
	} else {
		cpu.t += 2 * t_states[nop]
	}

	switch opcode {
	case neg, 0x54, 0x64, 0x74, 0x4C, 0x5C, 0x6C, 0x7C:
		a := cpu.reg.A
		cpu.reg.A = ^a + 1
		cpu.reg.F = f_N
		cpu.reg.F |= f_S & cpu.reg.A
		if cpu.reg.A == 0 {
			cpu.reg.F |= f_Z
		}
		cpu.reg.F |= byte(a^cpu.reg.A) & f_H
		if cpu.reg.A == 0x80 {
			cpu.reg.F |= f_P
		}
		if a != 0 {
			cpu.reg.F |= f_C
		}
	case adc_hl_bc, adc_hl_de, adc_hl_hl, adc_hl_sp:
		hl := cpu.reg.getHL()
		nn := cpu.reg.getReg16(opcode & 0b00110000 >> 4)
		sum := hl + nn + uint16(cpu.reg.F&f_C)
		cpu.reg.F = f_NONE
		if sum > 0x7FFF {
			cpu.reg.F |= f_S
		}
		if sum == 0 {
			cpu.reg.F |= f_Z
		}
		cpu.reg.F |= byte((hl^nn^sum)>>8) & f_H
		if (hl^nn)&0x8000 == 0 && (hl^sum)&0x8000 != 0 {
			cpu.reg.F |= f_P
		}
		if sum < hl {
			cpu.reg.F |= f_C
		}
		cpu.reg.setHLw(sum)
	case sbc_hl_bc, sbc_hl_de, sbc_hl_hl, sbc_hl_sp:
		hl := cpu.reg.getHL()
		nn := cpu.reg.getReg16(opcode & 0b00110000 >> 4)
		sub := hl - nn - uint16(cpu.reg.F&f_C)
		cpu.reg.F = f_N
		if sub > 0x7FFF {
			cpu.reg.F |= f_S
		}
		if sub == 0 {
			cpu.reg.F |= f_Z
		}
		cpu.reg.F |= byte((hl^nn^sub)>>8) & f_H
		if (hl^nn)&0x8000 != 0 && (hl^sub)&0x8000 != 0 {
			cpu.reg.F |= f_P
		}
		if sub > hl {
			cpu.reg.F |= f_C
		}
		cpu.reg.setHLw(sub)
	case rld:
		// TODO: Implement
	case rrd:
		// TODO: Implement
	case in_a_c, in_b_c, in_c_c, in_d_c, in_e_c, in_f_c, in_h_c, in_l_c:
		// TODO: Implement
	case out_c_a, out_c_b, out_c_c, out_c_d, out_c_e, out_c_f, out_c_h, out_c_l:
		// TODO: Implement
	case im0, im1, im2:
		// TODO: Implement
	case retn, 0x55, 0x65, 0x75, 0x5D, 0x6D:
		// TODO: Implement
	case reti, 0x7D:
		// TODO: Implement
	case ld_mm_bc, ld_mm_hl2, ld_mm_de, ld_mm_sp:
		// TODO: Implement
	case ld_bc_mm, ld_de_mm, ld_hl_mm2, ld_sp_mm:
		// TODO: Implement
	case ld_a_r:
		// TODO: Implement
	case ld_r_a:
		// TODO: Implement
	case ld_a_i:
		// TODO: Implement
	case ld_i_a:
		// TODO: Implement
	case ldi:
		// TODO: Implement
	case ldir:
		// TODO: Implement
	case cpi:
		// TODO: Implement
	case cpir:
		// TODO: Implement
	case ini:
		// TODO: Implement
	case inir:
		// TODO: Implement
	case outi:
		// TODO: Implement
	case otir:
		// TODO: Implement
	case ldd:
		// TODO: Implement
	case lddr:
		// TODO: Implement
	case cpd:
		// TODO: Implement
	case cpdr:
		// TODO: Implement
	case ind:
		// TODO: Implement
	case indr:
		// TODO: Implement
	case outd:
		// TODO: Implement
	case otdr:
		// TODO: Implement
	default:
		// NOP
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

	panic(fmt.Sprintf("Invalid opcode %v", opcode))
}

// Returns value of HL / (IX + d) / (IY + d) register. The prefix
// determines whether to use IX or IY register instead of HL.
// Offset argument specifies whether IX / IY should include offset.
func (cpu *CPU) getHL(offset bool) uint16 {
	hl := cpu.reg.getHL()
	if offset && (cpu.reg.prefix == useIX || cpu.reg.prefix == useIY) {
		o := cpu.readByte()
		if o&0x80 == 0 {
			return hl + uint16(o)
		} else {
			return hl - uint16(^o+1)
		}
	}
	return hl
}

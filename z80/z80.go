package z80

import (
	"fmt"

	"github.com/voytas/z80-go-zx/z80/memory"
)

// Represents emulated Z80 CPU
type CPU struct {
	IN               func(hi, lo byte) byte  // callback function to execute on IN instruction
	OUT              func(hi, lo, data byte) // callback function to execute on OUT instruction
	mem              memory.Memory           // memory
	reg              *registers              // registers
	t                byte                    // t-states
	halt, iff1, iff2 bool                    // states of halt, iff1 and iff2
	im               byte                    // interrupt mode (im0, im1 or in2)
}

func NewCPU(mem memory.Memory) *CPU {
	cpu := &CPU{}
	cpu.mem = mem
	cpu.IN = func(_, _ byte) byte { return 0xFF }
	cpu.OUT = func(hi, lo, data byte) {}
	cpu.Reset()
	return cpu
}

func (cpu *CPU) readByte() byte {
	b := cpu.mem.Read(cpu.reg.PC)
	cpu.reg.PC += 1
	return b
}

func (cpu *CPU) readWord() uint16 {
	w := uint16(cpu.mem.Read(cpu.reg.PC)) | uint16(cpu.mem.Read(cpu.reg.PC+1))<<8
	cpu.reg.PC += 2
	return w
}

func (cpu *CPU) wait() {
}

func (cpu *CPU) Reset() {
	cpu.reg = newRegisters()
	cpu.reg.PC, cpu.reg.SP = 0, 0xFFFF
	cpu.halt = false
	cpu.iff1, cpu.iff2 = false, false
	cpu.reg.A, cpu.reg.F, cpu.reg.I = 0xFF, 0xFF, 0x00
}

func (cpu *CPU) Run() {
	for {
		opcode := cpu.readByte()

		//cpu.debug(opcode)

		// Get the t-state for the current instruction
		if cpu.reg.prefix == noPrefix {
			cpu.t = tStates[opcode]
		} else {
			t := tStatesIXY[opcode]
			if t != 0 {
				cpu.t = t
			} else {
				cpu.t = 4 + tStates[opcode]
			}
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
		case rlca:
			a7 := cpu.reg.A >> 7
			cpu.reg.A = cpu.reg.A<<1 | a7
			cpu.reg.F = cpu.reg.F&(fS|fZ|fP) | cpu.reg.A&(fY|fX) | a7
		case rrca:
			a0 := cpu.reg.A & 0x01
			cpu.reg.A = cpu.reg.A>>1 | a0<<7
			cpu.reg.F = cpu.reg.F&(fS|fZ|fP) | cpu.reg.A&(fY|fX) | a0
		case rla:
			a7 := cpu.reg.A >> 7
			cpu.reg.A = cpu.reg.A<<1 | cpu.reg.F&fC
			cpu.reg.F = cpu.reg.F&(fS|fZ|fP) | cpu.reg.A&(fY|fX) | a7
		case rra:
			a0 := cpu.reg.A & 0x01
			cpu.reg.A = cpu.reg.A>>1 | cpu.reg.F&fC<<7
			cpu.reg.F = cpu.reg.F&(fS|fZ|fP) | cpu.reg.A&(fY|fX) | a0
		case cpl:
			cpu.reg.A = ^cpu.reg.A
			cpu.reg.F = cpu.reg.F&(fS|fZ|fP|fC) | fH | fN | cpu.reg.A&(fY|fX)
		case scf:
			cpu.reg.F = cpu.reg.F&(fS|fZ|fP) | fC | cpu.reg.A&(fY|fX)
		case ccf:
			cpu.reg.F = (cpu.reg.F&(fS|fZ|fP|fC) | cpu.reg.F&fC<<4 | cpu.reg.A&(fY|fX)) ^ fC
		case daa:
			cf := cpu.reg.F & fC
			hf := cpu.reg.F & fH
			nf := cpu.reg.F & fN
			lb := cpu.reg.A & 0x0F
			diff := byte(0)
			cpu.reg.F &= fN
			if cf != 0 || cpu.reg.A > 0x99 {
				diff = 0x60
				cpu.reg.F |= fC
			}
			if hf != 0 || lb > 0x09 {
				diff += 0x06
			}
			if nf == 0 {
				cpu.reg.A += diff
			} else {
				cpu.reg.A -= diff
			}
			cpu.reg.F |= parity[cpu.reg.A] | cpu.reg.A&(fS|fY|fX)
			if cpu.reg.A == 0 {
				cpu.reg.F |= fZ
			}
			if nf == 0 && lb > 0x09 || nf != 0 && hf != 0 && lb < 0x06 {
				cpu.reg.F |= fH
			}
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
			h, l := cpu.reg.r(rH), cpu.reg.r(rL)
			x, y := cpu.mem.Read(cpu.reg.SP+1), cpu.mem.Read(cpu.reg.SP)
			cpu.mem.Write(cpu.reg.SP, *l)
			cpu.mem.Write(cpu.reg.SP+1, *h)
			*h, *l = x, y
		case add_a_n, add_a_a, add_a_b, add_a_c, add_a_d, add_a_e, add_a_h, add_a_l, add_a_hl:
			a := cpu.reg.A
			var n byte
			switch opcode {
			case add_a_n:
				n = cpu.readByte()
			case add_a_hl:
				n = cpu.mem.Read(cpu.getHL())
			default:
				n = *cpu.reg.r(opcode & 0b00000111)
			}
			cpu.reg.A += n
			cpu.reg.F = (fS | fY | fX) & cpu.reg.A
			if cpu.reg.A == 0 {
				cpu.reg.F |= fZ
			}
			cpu.reg.F |= (a ^ n ^ cpu.reg.A) & fH
			if (a^n)&0x80 == 0 && (a^cpu.reg.A)&0x80 != 0 {
				cpu.reg.F |= fP
			}
			if cpu.reg.A < a {
				cpu.reg.F |= fC
			}
		case adc_a_n, adc_a_a, adc_a_b, adc_a_c, adc_a_d, adc_a_e, adc_a_h, adc_a_l, adc_a_hl:
			var n byte
			switch opcode {
			case adc_a_n:
				n = cpu.readByte()
			case adc_a_hl:
				n = cpu.mem.Read(cpu.getHL())
			default:
				n = *cpu.reg.r(opcode & 0b00000111)
			}
			cf := cpu.reg.F & fC
			sum_w := uint16(cpu.reg.A) + uint16(n) + uint16(cf)
			sum_b := byte(sum_w)
			cpu.reg.F = (fS | fY | fX) & sum_b
			if sum_b == 0 {
				cpu.reg.F |= fZ
			}
			cpu.reg.F |= (cpu.reg.A ^ n ^ sum_b) & fH
			if (cpu.reg.A^n)&0x80 == 0 && (cpu.reg.A^sum_b)&0x80 != 0 {
				cpu.reg.F |= fP
			}
			if sum_w > 0xff {
				cpu.reg.F |= fC
			}
			cpu.reg.A = sum_b
		case add_hl_bc, add_hl_de, add_hl_hl, add_hl_sp:
			hl := cpu.reg.HL()
			var nn uint16
			switch opcode {
			case add_hl_bc:
				nn = cpu.reg.BC()
			case add_hl_de:
				nn = cpu.reg.DE()
			case add_hl_hl:
				nn = hl
			case add_hl_sp:
				nn = cpu.reg.SP
			}
			sum := hl + nn
			cpu.reg.setHL(sum)
			cpu.reg.F = cpu.reg.F & ^(fH|fN|fC) | byte((hl^nn^sum)>>8)&fH | byte(sum>>8)&(fY|fX)
			if sum < hl {
				cpu.reg.F |= fC
			}
		case sub_a, sub_b, sub_c, sub_d, sub_e, sub_h, sub_l, sub_hl, sub_n:
			a := cpu.reg.A
			var n byte
			switch opcode {
			case sub_n:
				n = cpu.readByte()
			case sub_hl:
				n = cpu.mem.Read(cpu.getHL())
			default:
				n = *cpu.reg.r(opcode & 0b00000111)
			}
			cpu.reg.A -= n
			cpu.reg.F = (fS|fY|fX)&cpu.reg.A | fN | (a^n^cpu.reg.A)&fH
			if cpu.reg.A == 0 {
				cpu.reg.F |= fZ
			}
			if (a^n)&0x80 != 0 && (a^cpu.reg.A)&0x80 != 0 {
				cpu.reg.F |= fP
			}
			if cpu.reg.A > a {
				cpu.reg.F |= fC
			}
		case cp_a, cp_b, cp_c, cp_d, cp_e, cp_h, cp_l, cp_hl, cp_n:
			var n byte
			switch opcode {
			case cp_n:
				n = cpu.readByte()
			case cp_hl:
				n = cpu.mem.Read(cpu.getHL())
			default:
				n = *cpu.reg.r(opcode & 0b00000111)
			}
			test := cpu.reg.A - n
			cpu.reg.F = fN | fS&test | n&(fY|fX) | byte(cpu.reg.A^n^test)&fH
			if test == 0 {
				cpu.reg.F |= fZ
			}
			if (cpu.reg.A^n)&0x80 != 0 && (cpu.reg.A^test)&0x80 != 0 {
				cpu.reg.F |= fP
			}
			if test > cpu.reg.A {
				cpu.reg.F |= fC
			}
		case sbc_a_a, sbc_a_b, sbc_a_c, sbc_a_d, sbc_a_e, sbc_a_h, sbc_a_l, sbc_a_hl, sbc_a_n:
			var n byte
			switch opcode {
			case sbc_a_n:
				n = cpu.readByte()
			case sbc_a_hl:
				n = cpu.mem.Read(cpu.getHL())
			default:
				n = *cpu.reg.r(opcode & 0b00000111)
			}
			cf := cpu.reg.F & fC
			sub_w := uint16(cpu.reg.A) - uint16(n) - uint16(cf)
			sub_b := byte(sub_w)
			cpu.reg.F = (fS|fY|fX)&sub_b | fN | byte(cpu.reg.A^n^sub_b)&fH
			if sub_b == 0 {
				cpu.reg.F |= fZ
			}
			if (cpu.reg.A^n)&0x80 != 0 && (sub_b^cpu.reg.A)&0x80 != 0 {
				cpu.reg.F |= fP
			}
			if sub_w > 0xff {
				cpu.reg.F |= fC
			}
			cpu.reg.A = sub_b
		case and_a, and_b, and_c, and_d, and_e, and_h, and_l, and_hl, and_n:
			var n byte
			switch opcode {
			case and_n:
				n = cpu.readByte()
			case and_hl:
				n = cpu.mem.Read(cpu.getHL())
			default:
				n = *cpu.reg.r(opcode & 0b00000111)
			}
			cpu.reg.A &= n
			cpu.reg.F = (fS|fY|fX)&cpu.reg.A | fH | parity[cpu.reg.A]
			if cpu.reg.A == 0 {
				cpu.reg.F |= fZ
			}
		case or_a, or_b, or_c, or_d, or_e, or_h, or_l, or_hl, or_n:
			var n byte
			switch opcode {
			case or_n:
				n = cpu.readByte()
			case or_hl:
				n = cpu.mem.Read(cpu.getHL())
			default:
				n = *cpu.reg.r(opcode & 0b00000111)
			}
			cpu.reg.A |= n
			cpu.reg.F = (fS|fY|fX)&cpu.reg.A | parity[cpu.reg.A]
			if cpu.reg.A == 0 {
				cpu.reg.F |= fZ
			}
		case xor_a, xor_b, xor_c, xor_d, xor_e, xor_h, xor_l, xor_hl, xor_n:
			var n byte
			switch opcode {
			case xor_n:
				n = cpu.readByte()
			case xor_hl:
				n = cpu.mem.Read(cpu.getHL())
			default:
				n = *cpu.reg.r(opcode & 0b00000111)
			}
			cpu.reg.A ^= n
			cpu.reg.F = (fS|fY|fX)&cpu.reg.A | parity[cpu.reg.A]
			if cpu.reg.A == 0 {
				cpu.reg.F |= fZ
			}
		case ld_a_n, ld_b_n, ld_c_n, ld_d_n, ld_e_n, ld_h_n, ld_l_n:
			r := cpu.reg.r(opcode & 0b00111000 >> 3)
			*r = cpu.readByte()
		case
			ld_a_a, ld_a_b, ld_a_c, ld_a_d, ld_a_e, ld_a_h, ld_a_l,
			ld_b_a, ld_b_b, ld_b_c, ld_b_d, ld_b_e, ld_b_h, ld_b_l,
			ld_c_a, ld_c_b, ld_c_c, ld_c_d, ld_c_e, ld_c_h, ld_c_l,
			ld_d_a, ld_d_b, ld_d_c, ld_d_d, ld_d_e, ld_d_h, ld_d_l,
			ld_e_a, ld_e_b, ld_e_c, ld_e_d, ld_e_e, ld_e_h, ld_e_l,
			ld_h_a, ld_h_b, ld_h_c, ld_h_d, ld_h_e, ld_h_h, ld_h_l,
			ld_l_a, ld_l_b, ld_l_c, ld_l_d, ld_l_e, ld_l_h, ld_l_l:
			rs := cpu.reg.r(opcode & 0b00000111)
			rd := cpu.reg.r(opcode & 0b00111000 >> 3)
			*rd = *rs
		case ld_bc_nn:
			cpu.reg.C, cpu.reg.B = cpu.readByte(), cpu.readByte()
		case ld_de_nn:
			cpu.reg.E, cpu.reg.D = cpu.readByte(), cpu.readByte()
		case ld_hl_nn:
			h, l := cpu.reg.r(rH), cpu.reg.r(rL)
			*l, *h = cpu.readByte(), cpu.readByte()
		case ld_sp_nn:
			cpu.reg.SP = cpu.readWord()
		case ld_sp_hl:
			cpu.reg.SP = cpu.reg.HL()
		case ld_hl_mm:
			addr := cpu.readWord()
			h, l := cpu.reg.r(rH), cpu.reg.r(rL)
			*l = cpu.mem.Read(addr)
			*h = cpu.mem.Read(addr + 1)
		case ld_mm_hl:
			addr := cpu.readWord()
			h, l := cpu.reg.r(rH), cpu.reg.r(rL)
			cpu.mem.Write(addr, *l)
			cpu.mem.Write(addr+1, *h)
		case ld_mhl_n:
			cpu.mem.Write(cpu.getHL(), cpu.readByte())
		case ld_mm_a:
			cpu.mem.Write(cpu.readWord(), cpu.reg.A)
		case ld_a_mm:
			cpu.reg.A = cpu.mem.Read(cpu.readWord())
		case ld_bc_a:
			cpu.mem.Write(cpu.reg.BC(), cpu.reg.A)
		case ld_de_a:
			cpu.mem.Write(cpu.reg.DE(), cpu.reg.A)
		case ld_a_bc:
			cpu.reg.A = cpu.mem.Read(cpu.reg.BC())
		case ld_a_de:
			cpu.reg.A = cpu.mem.Read(cpu.reg.DE())
		case ld_a_hl, ld_b_hl, ld_c_hl, ld_d_hl, ld_e_hl, ld_h_hl, ld_l_hl:
			*cpu.reg.raw[opcode&0b00111000>>3] = cpu.mem.Read(cpu.getHL())
		case ld_hl_a, ld_hl_b, ld_hl_c, ld_hl_d, ld_hl_e, ld_hl_h, ld_hl_l:
			cpu.mem.Write(cpu.getHL(), *cpu.reg.raw[opcode&0b00000111])
		case inc_a, inc_b, inc_c, inc_d, inc_e, inc_h, inc_l:
			r := cpu.reg.r(opcode & 0b00111000 >> 3)
			cpu.reg.F &= fC
			if *r == 0x7F {
				cpu.reg.F |= fP
			}
			if *r&0x0F == 0x0F {
				cpu.reg.F |= fH
			}
			*r += 1
			cpu.reg.F |= *r & (fS | fY | fX)
			if *r == 0 {
				cpu.reg.F |= fZ
			}
		case inc_bc:
			cpu.reg.setBC(cpu.reg.BC() + 1)
		case inc_de:
			cpu.reg.setDE(cpu.reg.DE() + 1)
		case inc_hl:
			cpu.reg.setHL(cpu.reg.HL() + 1)
		case inc_sp:
			cpu.reg.SP += 1
		case inc_mhl:
			addr := cpu.getHL()
			b := cpu.mem.Read(addr)
			cpu.reg.F &= fC
			if b == 0x7F {
				cpu.reg.F |= fP
			}
			if b&0x0F == 0x0F {
				cpu.reg.F |= fH
			}
			b += 1
			if b == 0x00 {
				cpu.reg.F |= fZ
			}
			cpu.reg.F |= b & (fS | fY | fX)
			cpu.mem.Write(addr, b)
		case dec_a, dec_b, dec_c, dec_d, dec_e, dec_h, dec_l:
			r := cpu.reg.r(opcode & 0b00111000 >> 3)
			cpu.reg.F = cpu.reg.F&fC | fN
			if *r == 0x80 {
				cpu.reg.F |= fP
			}
			if *r&0x0F == 0 {
				cpu.reg.F |= fH
			}
			*r -= 1
			cpu.reg.F |= *r & (fS | fY | fX)
			if *r == 0 {
				cpu.reg.F |= fZ
			}
		case dec_bc:
			cpu.reg.setBC(cpu.reg.BC() - 1)
		case dec_de:
			cpu.reg.setDE(cpu.reg.DE() - 1)
		case dec_hl:
			cpu.reg.setHL(cpu.reg.HL() - 1)
		case dec_sp:
			cpu.reg.SP -= 1
		case dec_mhl:
			addr := cpu.getHL()
			b := cpu.mem.Read(addr)
			cpu.reg.F = cpu.reg.F&fC | fN
			if b == 0x80 {
				cpu.reg.F |= fP
			}
			if b&0x0F == 0 {
				cpu.reg.F |= fH
			}
			b -= 1
			if b == 0x00 {
				cpu.reg.F |= fZ
			}
			cpu.reg.F |= b & (fS | fY | fX)
			cpu.mem.Write(addr, b)
		case jr_o:
			o := cpu.readByte()
			if o&0x80 == 0 {
				cpu.reg.PC += uint16(o)
			} else {
				cpu.reg.PC -= uint16(^o + 1)
			}
		case jr_z_o:
			o := cpu.readByte()
			if cpu.reg.F&fZ == fZ {
				if o&0x80 == 0 {
					cpu.reg.PC += uint16(o)
				} else {
					cpu.reg.PC -= uint16(^o + 1)
				}
				cpu.t += 5
			}
		case jr_nz_o:
			o := cpu.readByte()
			if cpu.reg.F&fZ == 0 {
				if o&0x80 == 0 {
					cpu.reg.PC += uint16(o)
				} else {
					cpu.reg.PC -= uint16(^o + 1)
				}
				cpu.t += 5
			}
		case jr_c:
			o := cpu.readByte()
			if cpu.reg.F&fC == fC {
				if o&0x80 == 0 {
					cpu.reg.PC += uint16(o)
				} else {
					cpu.reg.PC -= uint16(^o + 1)
				}
				cpu.t += 5
			}
		case jr_nc_o:
			o := cpu.readByte()
			if cpu.reg.F&fC == 0 {
				if o&0x80 == 0 {
					cpu.reg.PC += uint16(o)
				} else {
					cpu.reg.PC -= uint16(^o + 1)
				}
				cpu.t += 5
			}
		case djnz:
			o := cpu.readByte()
			cpu.reg.B -= 1
			if cpu.reg.B != 0 {
				if o&0x80 == 0 {
					cpu.reg.PC += uint16(o)
				} else {
					cpu.reg.PC -= uint16(^o + 1)
				}
				cpu.t += 5
			}
		case jp_nn:
			cpu.reg.PC = cpu.readWord()
		case jp_c_nn, jp_m_nn, jp_nc_nn, jp_nz_nn, jp_p_nn, jp_pe_nn, jp_po_nn, jp_z_nn:
			if cpu.shouldJump(opcode) {
				cpu.reg.PC = cpu.readWord()
			} else {
				cpu.reg.PC += 2
			}
		case jp_hl:
			cpu.reg.PC = cpu.reg.HL()
		case call_nn, call_c_nn, call_m_nn, call_nc_nn, call_nz_nn, call_p_nn, call_pe_nn, call_po_nn, call_z_nn:
			if cpu.shouldJump(opcode) {
				pc := cpu.readWord()
				cpu.reg.SP -= 1
				cpu.mem.Write(cpu.reg.SP, byte(cpu.reg.PC>>8))
				cpu.reg.SP -= 1
				cpu.mem.Write(cpu.reg.SP, byte(cpu.reg.PC))
				cpu.reg.PC = pc
				cpu.t += 7
			} else {
				cpu.reg.PC += 2
			}
		case ret:
			cpu.reg.PC = uint16(cpu.mem.Read(cpu.reg.SP+1))<<8 | uint16(cpu.mem.Read(cpu.reg.SP))
			cpu.reg.SP += 2
		case ret_c, ret_m, ret_nc, ret_nz, ret_p, ret_pe, ret_po, ret_z:
			if cpu.shouldJump(opcode) {
				cpu.reg.PC = uint16(cpu.mem.Read(cpu.reg.SP+1))<<8 | uint16(cpu.mem.Read(cpu.reg.SP))
				cpu.reg.SP += 2
				cpu.t += 6
			}
		case rst_00h, rst_08h, rst_10h, rst_18h, rst_20h, rst_28h, rst_30h, rst_38h:
			cpu.reg.SP -= 1
			cpu.mem.Write(cpu.reg.SP, byte(cpu.reg.PC>>8))
			cpu.reg.SP -= 1
			cpu.mem.Write(cpu.reg.SP, byte(cpu.reg.PC))
			cpu.reg.PC = uint16(8 * ((opcode & 0b00111000) >> 3))
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
			cpu.mem.Write(cpu.reg.SP, *cpu.reg.r(rH))
			cpu.reg.SP -= 1
			cpu.mem.Write(cpu.reg.SP, *cpu.reg.r(rL))
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
			*cpu.reg.r(rH), *cpu.reg.r(rL) = cpu.mem.Read(cpu.reg.SP+1), cpu.mem.Read(cpu.reg.SP)
			cpu.reg.SP += 2
		case in_a_n:
			cpu.reg.A = cpu.IN(cpu.reg.A, cpu.readByte())
		case out_n_a:
			cpu.OUT(cpu.reg.A, cpu.readByte(), cpu.reg.A)
		case prefix_cb:
			cpu.prefixCB()
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

func (cpu *CPU) shouldJump(opcode byte) bool {
	if opcode == call_nn {
		return true
	}

	switch opcode & 0b00111000 {
	case 0b00000000: // Non-Zero (NZ)
		return cpu.reg.F&fZ == 0
	case 0b00001000: // Zero (Z)
		return cpu.reg.F&fZ != 0
	case 0b00010000: // Non Carry (NC)
		return cpu.reg.F&fC == 0
	case 0b00011000: // Carry (C)
		return cpu.reg.F&fC != 0
	case 0b00100000: // Parity Odd (PO)
		return cpu.reg.F&fP == 0
	case 0b00101000: // Parity Even (PE)
		return cpu.reg.F&fP != 0
	case 0b00110000: // Sign Positive (P)
		return cpu.reg.F&fS == 0
	case 0b00111000: // Sign Negative (M)
		return cpu.reg.F&fS != 0
	}

	panic(fmt.Sprintf("Invalid opcode %v", opcode))
}

// Returns value of HL / (IX + d) / (IY + d) register. The current prefix
// determines whether to use IX or IY register instead of HL.
func (cpu *CPU) getHL() uint16 {
	if cpu.reg.prefix != noPrefix {
		return cpu.getHLOffset(cpu.readByte())
	}
	return cpu.reg.HL()
}

// For IX or IY add offset to register value, otherwise return HL.
func (cpu *CPU) getHLOffset(offset byte) uint16 {
	hl := cpu.reg.HL()
	if offset == 0 {
		return hl
	} else if offset&0x80 == 0 {
		return hl + uint16(offset)
	} else {
		return hl - uint16(^offset+1)
	}
}

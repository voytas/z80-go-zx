package z80

type word uint16

type CPU struct {
	PC   word
	mem  Memory
	r    *registers
	halt bool
	IN   func(a, n byte) byte
}

func NewCPU(mem Memory) *CPU {
	cpu := &CPU{}
	cpu.mem = mem
	cpu.Reset()
	return cpu
}

func (c *CPU) readByte() byte {
	b := c.mem.read(c.PC)
	c.PC += 1
	return b
}

func (c *CPU) readWord() word {
	w := word(c.mem.read(c.PC)) + word(c.mem.read(c.PC+1))<<8
	c.PC += 2
	return w
}

func (c *CPU) wait(t byte) {}

func (c *CPU) Reset() {
	c.PC = 0
	c.r = newRegisters()
	c.halt = false
}

func (c *CPU) Run() {
	for {
		opcode := c.readByte()

		var t byte // T states
		switch opcode {
		case nop:
			t = 4
		case halt:
			t = 4
			c.wait(t)
			c.halt = true
			return
		case cpl:
			c.r.A = ^c.r.A
			c.r.F |= f_H | f_N
			t = 4
		case scf:
			c.r.F &= ^(f_H | f_N)
			c.r.F |= f_C
			t = 4
		case ccf:
			c.r.F &= ^(f_N)
			c.r.F ^= f_H | f_C
			t = 4
		case rlca:
			c.r.F &= ^(f_H | f_N | f_C)
			a7 := c.r.A >> 7
			c.r.F |= a7
			c.r.A = c.r.A<<1 | a7
			t = 4
		case rrca:
			c.r.F &= ^(f_H | f_N | f_C)
			a7 := c.r.A & 0x01
			c.r.F |= a7
			c.r.A = c.r.A>>1 | a7<<7
			t = 4
		case rla:
			fc := c.r.F & f_C
			c.r.F &= ^(f_H | f_N | f_C)
			a7 := c.r.A >> 7
			c.r.A = c.r.A<<1 | fc
			c.r.F |= a7
			t = 4
		case rra:
			fc := c.r.F & f_C
			c.r.F &= ^(f_H | f_N | f_C)
			a0 := c.r.A & 0x01
			c.r.A = c.r.A>>1 | fc<<7
			c.r.F |= a0
			t = 4
		case daa:
			c.r.F &= ^(f_S | f_Z | f_P)
			a := c.r.A
			if a&0xF > 9 || c.r.F&f_H != 0 {
				if c.r.F&f_N > 0 {
					a -= 0x06
				} else {
					a += 0x06
				}
			}
			if a > 0x99 || c.r.F&f_C != 0 {
				if c.r.F&f_N > 0 {
					a -= 0x60
				} else {
					a += 0x60
				}
			}
			c.r.F |= a & f_S
			if a == 0 {
				c.r.F |= f_Z
			}
			if c.r.F&f_N > 0 {
				if c.r.F&f_H > 0 && c.r.A&0xF < 6 {
					c.r.F |= f_H
				} else {
					c.r.F &= ^f_H
				}
			} else {
				if c.r.A&0xF > 9 {
					c.r.F |= f_H
				} else {
					c.r.F &= ^f_H
				}
			}
			c.r.F |= parity[a]
			if c.r.A > 0x99 {
				c.r.F |= f_C
			}
			c.r.A = a
			t = 4
		case ex_af_af:
			c.r.A, c.r.A_ = c.r.A_, c.r.A
			c.r.F, c.r.F_ = c.r.F_, c.r.F
			t = 4
		case exx:
			c.r.B, c.r.B_, c.r.C, c.r.C_ = c.r.B_, c.r.B, c.r.C_, c.r.C
			c.r.D, c.r.D_, c.r.E, c.r.E_ = c.r.D_, c.r.D, c.r.E_, c.r.E
			c.r.H, c.r.H_, c.r.L, c.r.L_ = c.r.H_, c.r.H, c.r.L_, c.r.L
			t = 4
		case add_a_n, add_a_a, add_a_b, add_a_c, add_a_d, add_a_e, add_a_h, add_a_l, add_a_hl:
			var n byte
			if opcode == add_a_n {
				n = c.readByte()
				t = 7
			} else if opcode == add_a_hl {
				n = c.mem.read(c.r.getHL())
				t = 7
			} else {
				n = *c.r.getR(opcode & 0b00000111)
				t = 4
			}
			c.r.F = f_NONE
			sum := c.r.A + n
			if sum > 0x7F {
				c.r.F |= f_S
			} else if sum == 0 {
				c.r.F |= f_Z
			}
			c.r.F |= (c.r.A ^ n ^ sum) & f_H
			if (c.r.A^n)&0x80 == 0 && (c.r.A^sum)&0x80 > 0 {
				c.r.F |= f_P
			}
			if sum < c.r.A {
				c.r.F |= f_C
			}
			c.r.A = sum
		case adc_a_n, adc_a_a, adc_a_b, adc_a_c, adc_a_d, adc_a_e, adc_a_h, adc_a_l, adc_a_hl:
			var n byte
			if opcode == adc_a_n {
				n = c.readByte()
				t = 7
			} else if opcode == adc_a_hl {
				n = c.mem.read(c.r.getHL())
				t = 7
			} else {
				n = *c.r.getR(opcode & 0b00000111)
				t = 4
			}
			cf := c.r.F & f_C
			c.r.F = f_NONE
			sum_w := word(c.r.A) + word(n) + word(cf)
			sum_b := byte(sum_w)
			c.r.F |= f_S & sum_b
			if sum_b == 0 {
				c.r.F |= f_Z
			}
			c.r.F |= (c.r.A ^ n ^ sum_b) & f_H
			if (c.r.A^n)&0x80 == 0 && (c.r.A^sum_b)&0x80 > 0 {
				c.r.F |= f_P
			}
			if sum_w > 0xff {
				c.r.F |= f_C
			}
			c.r.A = sum_b
		case add_hl_bc, add_hl_de, add_hl_hl, add_hl_sp:
			hl := c.r.getHL()
			var nn word
			switch opcode {
			case add_hl_bc:
				nn = c.r.getBC()
			case add_hl_de:
				nn = c.r.getDE()
			case add_hl_hl:
				nn = hl
			case add_hl_sp:
				nn = c.r.SP
			}
			sum := hl + nn
			c.r.setHL(sum)
			c.r.F &= ^(f_H | f_N | f_C)
			if sum < hl {
				c.r.F |= f_C
			}
			c.r.F |= byte((hl^nn^sum)>>8) & f_H
			t = 11
		case sub_n, sub_a, sub_b, sub_c, sub_d, sub_e, sub_h, sub_l, sub_hl:
			a := c.r.A
			var n byte
			if opcode == sub_n {
				n = c.readByte()
				t = 7
			} else if opcode == sub_hl {
				n = c.mem.read(c.r.getHL())
				t = 7
			} else {
				n = *c.r.getR(opcode & 0b00000111)
				t = 4
			}
			c.r.A -= n
			c.r.F = f_N
			c.r.F |= f_S & c.r.A
			if c.r.A == 0 {
				c.r.F |= f_Z
			}
			c.r.F |= byte(a^n^c.r.A) & f_H
			if (a^n)&0x80 > 0 && (a^c.r.A)&0x80 > 0 {
				c.r.F |= f_P
			}
			if c.r.A > a {
				c.r.F |= f_C
			}
		case cp_a, cp_b, cp_c, cp_d, cp_e, cp_h, cp_l, cp_hl:
			var n byte
			if opcode == cp_hl {
				n = c.mem.read(c.r.getHL())
				t = 7
			} else {
				n = *c.r.getR(opcode & 0b00000111)
				t = 4
			}
			test := c.r.A - n
			c.r.F = f_N
			c.r.F |= f_S & test
			if test == 0 {
				c.r.F |= f_Z
			}
			c.r.F |= byte(c.r.A^n^test) & f_H
			if (c.r.A^n)&0x80 > 0 && (c.r.A^test)&0x80 > 0 {
				c.r.F |= f_P
			}
			if test > c.r.A {
				c.r.F |= f_C
			}
		case sbc_a_a, sbc_a_b, sbc_a_c, sbc_a_d, sbc_a_e, sbc_a_h, sbc_a_l, sbc_a_hl:
			var n byte
			if opcode == sbc_a_hl {
				n = c.mem.read(c.r.getHL())
				t = 7
			} else {
				n = *c.r.getR(opcode & 0b00000111)
				t = 4
			}
			cf := c.r.F & f_C
			c.r.F = f_N
			sub_w := word(c.r.A) - word(n) - word(cf)
			sub_b := byte(sub_w)
			c.r.F |= f_S & sub_b
			if sub_b == 0 {
				c.r.F |= f_Z
			}
			c.r.F |= byte(c.r.A^n^sub_b) & f_H
			if (c.r.A^n)&0x80 > 0 && (sub_b^c.r.A)&0x80 > 0 {
				c.r.F |= f_P
			}
			if sub_w > 0xff {
				c.r.F |= f_C
			}
			c.r.A = sub_b
		case and_a, and_b, and_c, and_d, and_e, and_h, and_l, and_hl:
			var n byte
			if opcode == and_hl {
				n = c.mem.read(c.r.getHL())
				t = 7
			} else {
				n = *c.r.getR(opcode & 0b00000111)
				t = 4
			}
			c.r.F = f_H
			c.r.A &= n
			c.r.F |= f_S & c.r.A
			if c.r.A == 0 {
				c.r.F |= f_Z
			}
			c.r.F |= parity[c.r.A]
		case or_a, or_b, or_c, or_d, or_e, or_h, or_l, or_hl:
			var n byte
			if opcode == or_hl {
				n = c.mem.read(c.r.getHL())
				t = 7
			} else {
				t = 4
				n = *c.r.getR(opcode & 0b00000111)
			}
			c.r.F = f_NONE
			c.r.A |= n
			c.r.F |= f_S & c.r.A
			if c.r.A == 0 {
				c.r.F |= f_Z
			}
			c.r.F |= parity[c.r.A]
		case xor_a, xor_b, xor_c, xor_d, xor_e, xor_h, XOR_L, xor_hl:
			var n byte
			if opcode == xor_hl {
				n = c.mem.read(c.r.getHL())
				t = 7
			} else {
				n = *c.r.getR(opcode & 0b00000111)
				t = 4
			}
			c.r.F = f_NONE
			c.r.A ^= n
			c.r.F |= f_S & c.r.A
			if c.r.A == 0 {
				c.r.F |= f_Z
			}
			c.r.F |= parity[c.r.A]
		case ld_a_n, ld_b_n, ld_c_n, ld_d_n, ld_e_n, ld_h_n, ld_l_n:
			r := c.r.getR(opcode & 0b00111000 >> 3)
			*r = c.readByte()
			t = 7
		case
			ld_a_a, ld_a_b, ld_a_c, ld_a_d, ld_a_e, ld_a_h, ld_a_l,
			ld_b_a, ld_b_b, ld_b_c, ld_b_d, ld_b_e, ld_b_h, ld_b_l,
			ld_c_a, ld_c_b, ld_c_c, ld_c_d, ld_c_e, ld_c_h, ld_c_l,
			ld_d_a, ld_d_b, ld_d_c, ld_d_d, ld_d_e, ld_d_h, ld_d_l,
			ld_e_a, ld_e_b, ld_e_c, ld_e_d, ld_e_e, ld_e_h, ld_e_l,
			ld_h_a, ld_h_b, ld_h_c, ld_h_d, ld_h_e, ld_h_h, ld_h_l,
			ld_l_a, ld_l_b, ld_l_c, ld_l_d, ld_l_e, ld_l_h, ld_l_l:
			rs := c.r.getR(opcode & 0b00000111)
			rd := c.r.getR(opcode & 0b00111000 >> 3)
			*rd = *rs
			t = 4
		case ld_bc_nn:
			c.r.C, c.r.B = c.readByte(), c.readByte()
			t = 10
		case ld_de_nn:
			c.r.E, c.r.D = c.readByte(), c.readByte()
			t = 10
		case ld_hl_nn:
			c.r.L, c.r.H = c.readByte(), c.readByte()
			t = 10
		case ld_sp_nn:
			c.r.SP = word(c.readByte()) | word(c.readByte())<<8
			t = 10
		case ld_hl_mm:
			w := c.readWord()
			c.r.L = c.mem.read(w)
			c.r.H = c.mem.read(w + 1)
			t = 16
		case ld_mm_hl:
			w := c.readWord()
			c.mem.write(w, c.r.L)
			c.mem.write(w+1, c.r.H)
			t = 16
		case ld_mhl_n:
			c.mem.write(c.r.getHL(), c.readByte())
			t = 10
		case ld_mm_a:
			w := c.readWord()
			c.mem.write(w, c.r.A)
			t = 13
		case ld_a_mm:
			w := c.readWord()
			c.r.A = c.mem.read(w)
			t = 13
		case ld_bc_a:
			c.mem.write(c.r.getBC(), c.r.A)
			t = 7
		case ld_de_a:
			c.mem.write(c.r.getDE(), c.r.A)
			t = 7
		case ld_a_bc:
			c.r.A = c.mem.read(c.r.getBC())
			t = 7
		case ld_a_de:
			c.r.A = c.mem.read(c.r.getDE())
			t = 7
		case ld_a_hl, ld_b_hl, ld_c_hl, ld_d_hl, ld_e_hl, ld_h_hl, ld_l_hl:
			r := c.r.getR(opcode & 0b00111000 >> 3)
			*r = c.mem.read(c.r.getHL())
			t = 7
		case ld_hl_a, ld_hl_b, ld_hl_c, ld_hl_d, ld_hl_e, ld_hl_h, ld_hl_l:
			r := c.r.getR(opcode & 0b00000111)
			c.mem.write(c.r.getHL(), *r)
			t = 7
		case inc_a, inc_b, inc_c, inc_d, inc_e, inc_h, inc_l:
			r := c.r.getR(opcode & 0b00111000 >> 3)
			c.r.F &= ^(f_S | f_Z | f_H | f_P | f_N)
			if *r == 0x7F {
				c.r.F |= f_P
			}
			if *r&0x0F == 0x0F {
				c.r.F |= f_H
			}
			*r += 1
			if *r > 0x7F {
				c.r.F |= f_S
			}
			if *r == 0 {
				c.r.F |= f_Z
			}
			t = 4
		case inc_bc:
			c.r.setBC(c.r.getBC() + 1)
			t = 6
		case inc_de:
			c.r.setDE(c.r.getDE() + 1)
			t = 6
		case inc_hl:
			c.r.setHL(c.r.getHL() + 1)
			t = 6
		case inc_sp:
			c.r.SP += 1
			t = 6
		case inc_mhl:
			mm := c.r.getHL()
			b := c.mem.read(mm)
			c.r.F &= ^(f_S | f_Z | f_P | f_N)
			if b == 0x7F {
				c.r.F |= f_P
			}
			if b&0x0F == 0x0F {
				c.r.F |= f_H
			}
			b += 1
			if b == 0x00 {
				c.r.F |= f_Z
			}
			if b > 0x7F {
				c.r.F |= f_S
			}
			c.mem.write(mm, b)
			t = 11
		case dec_a, dec_b, dec_c, dec_d, dec_e, dec_h, dec_l:
			r := c.r.getR(opcode & 0b00111000 >> 3)
			c.r.F = c.r.F & ^(f_S|f_Z|f_H|f_P) | f_N
			if *r == 0x80 {
				c.r.F |= f_P
			}
			if *r&0x0F == 0 {
				c.r.F |= f_H
			}
			*r -= 1
			if *r > 0x7F {
				c.r.F |= f_S
			}
			if *r == 0 {
				c.r.F |= f_Z
			}
			t = 4
		case dec_bc:
			c.r.setBC(c.r.getBC() - 1)
			t = 6
		case dec_de:
			c.r.setDE(c.r.getDE() - 1)
			t = 6
		case dec_hl:
			c.r.setHL(c.r.getHL() - 1)
			t = 6
		case dec_sp:
			c.r.SP -= 1
			t = 6
		case dec_mhl:
			mm := c.r.getHL()
			b := c.mem.read(mm)
			c.r.F &= ^(f_S | f_Z | f_P)
			c.r.F |= f_N
			if b == 0x80 {
				c.r.F |= f_P
			}
			if b&0x0F == 0 {
				c.r.F |= f_H
			}
			b -= 1
			if b == 0x00 {
				c.r.F |= f_Z
			}
			if b > 0x7F {
				c.r.F |= f_S
			}
			c.mem.write(mm, b)
			t = 11
		case jr_o:
			o := c.readByte()
			if o&0x80 == 0 {
				c.PC += word(o)
			} else {
				c.PC -= word(^o + 1)
			}
			t = 12
		case jr_z_o:
			o := c.readByte()
			if c.r.F&f_Z == f_Z {
				if o&0x80 == 0 {
					c.PC += word(o)
				} else {
					c.PC -= word(^o + 1)
				}
				t = 12
			} else {
				t = 7
			}
		case jr_nz_o:
			o := c.readByte()
			if c.r.F&f_Z == 0 {
				if o&0x80 == 0 {
					c.PC += word(o)
				} else {
					c.PC -= word(^o + 1)
				}
				t = 12
			} else {
				t = 7
			}
		case jr_c:
			o := c.readByte()
			if c.r.F&f_C == f_C {
				if o&0x80 == 0 {
					c.PC += word(o)
				} else {
					c.PC -= word(^o + 1)
				}
				t = 12
			} else {
				t = 7
			}
		case jr_nc_o:
			o := c.readByte()
			if c.r.F&f_C == 0 {
				if o&0x80 == 0 {
					c.PC += word(o)
				} else {
					c.PC -= word(^o + 1)
				}
				t = 12
			} else {
				t = 7
			}
		case djnz:
			o := c.readByte()
			c.r.B -= 1
			if c.r.B == 0 {
				t = 8
			} else {
				if o&0x80 == 0 {
					c.PC += word(o)
				} else {
					c.PC -= word(^o + 1)
				}
				t = 13
			}
		case jp_nn:
			c.PC = word(c.readByte()) | word(c.readByte())<<8
			t = 10
		case jp_c_nn, jp_m_nn, jp_nc_nn, jp_nz_nn, jp_p_nn, jp_pe_nn, jp_po_nn, jp_z_nn:
			if c.shouldJump(opcode) {
				c.PC = word(c.readByte()) | word(c.readByte())<<8
			}
			t = 10
		case call_nn, call_c_nn, CALL_M_nn, call_nc_nn, call_nz_nn, call_p_nn, call_pe_nn, call_po_nn, call_z_nn:
			if c.shouldJump(opcode) {
				pc := word(c.readByte()) | word(c.readByte())<<8
				c.r.SP -= 1
				c.mem.write(c.r.SP, byte(c.PC>>8))
				c.r.SP -= 1
				c.mem.write(c.r.SP, byte(c.PC))
				c.PC = pc
				t = 17
			} else {
				c.PC += 2
				t = 10
			}
		case ret:
			c.PC = word(c.mem.read(c.r.SP+1))<<8 | word(c.mem.read(c.r.SP))
			c.r.SP += 2
			t = 10
		case ret_c, ret_m, ret_nc, ret_nz, ret_p, ret_pe, ret_po, ret_z:
			if c.shouldJump(opcode) {
				c.PC = word(c.mem.read(c.r.SP+1))<<8 | word(c.mem.read(c.r.SP))
				c.r.SP += 2
				t = 11
			} else {
				t = 5
			}
		case rst_00h, rst_08h, rst_10h, rst_18h, rst_20h, rst_28h, rst_30h, rst_38h:
			c.r.SP -= 1
			c.mem.write(c.r.SP, byte(c.PC>>8))
			c.r.SP -= 1
			c.mem.write(c.r.SP, byte(c.PC))
			c.PC = word(8 * ((opcode & 0b00111000) >> 3))
			t = 11
		case push_af:
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.A)
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.F)
			t = 11
		case push_bc:
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.B)
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.C)
			t = 11
		case push_de:
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.D)
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.E)
			t = 11
		case push_hl:
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.H)
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.L)
			t = 11
		case pop_af:
			c.r.A, c.r.F = c.mem.read(c.r.SP+1), c.mem.read(c.r.SP)
			c.r.SP += 2
			t = 10
		case pop_bc:
			c.r.B, c.r.C = c.mem.read(c.r.SP+1), c.mem.read(c.r.SP)
			c.r.SP += 2
			t = 10
		case pop_de:
			c.r.D, c.r.E = c.mem.read(c.r.SP+1), c.mem.read(c.r.SP)
			c.r.SP += 2
			t = 10
		case pop_hl:
			c.r.H, c.r.L = c.mem.read(c.r.SP+1), c.mem.read(c.r.SP)
			c.r.SP += 2
			t = 10
		case in_a_n:
			n := c.readByte()
			if c.IN != nil {
				c.r.A = c.IN(c.r.A, n)
			}
			t = 11
		case prefix_cb:
			c.cb(c.readByte(), &t)
		case prefix_ix:
			c.ix(c.readByte())
		case prefix_iy:
			c.iy(c.readByte())
		}

		c.wait(t)
	}
}

func (c *CPU) shouldJump(opcode byte) bool {
	if opcode == call_nn {
		return true
	}

	switch opcode & 0b00111000 {
	case 0b00000000: // Non-Zero (NZ)
		return c.r.F&f_Z == 0
	case 0b00001000: // Zero (Z)
		return c.r.F&f_Z != 0
	case 0b00010000: // Non Carry (NC)
		return c.r.F&f_C == 0
	case 0b00011000: // Carry (C)
		return c.r.F&f_C != 0
	case 0b00100000: // Parity Odd (PO)
		return c.r.F&f_P == 0
	case 0b00101000: // Parity Even (PE)
		return c.r.F&f_P != 0
	case 0b00110000: // Sign Positive (P)
		return c.r.F&f_S == 0
	case 0b00111000: // Sign Negative (M)
		return c.r.F&f_S != 0
	}

	panic("Invalid opcode")
}

func (c *CPU) cb(opcode byte, t *byte) {
	var v byte
	var hl word

	reg := opcode & 0b00000111
	if reg == r_HL {
		hl = c.r.getHL()
		v = c.mem.read(hl)
		*t = 15
	} else {
		v = *c.r.regs8[reg]
		*t = 8
	}

	var cy byte
	write := func() {
		if reg == r_HL {
			c.mem.write(hl, v)
		} else {
			*c.r.regs8[reg] = v
		}
		c.r.F = f_NONE
		c.r.F |= f_S & v
		if v == 0 {
			c.r.F |= f_Z
		}
		c.r.F |= parity[v] | cy
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
		v = v<<1 | f_C&c.r.F
		write()
	case rr_r:
		cy = v & f_C
		v = v>>1 | f_C&c.r.F<<7
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
				*t = 12
			}
			c.r.F &= ^(f_Z | f_N)
			c.r.F |= f_H
			if v&bit_mask[bit] == 0 {
				c.r.F |= f_Z
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

func (c *CPU) ix(opcode byte) {
	if opcode == prefix_ix || opcode == prefix_ED || opcode == prefix_iy {
		//t = 4
	}
}

func (c *CPU) iy(opcode byte) {

}

package z80

type word uint16

type CPU struct {
	PC   word
	mem  Memory
	r    *registers
	halt bool
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
		case NOP:
			t = 4
		case HALT:
			t = 4
			c.wait(t)
			c.halt = true
			return
		case CPL:
			c.r.A = ^c.r.A
			c.r.F |= f_H | f_N
			t = 4
		case SCF:
			c.r.F &= ^(f_H | f_N)
			c.r.F |= f_C
			t = 4
		case CCF:
			c.r.F &= ^(f_N)
			c.r.F ^= f_H | f_C
			t = 4
		case RLCA:
			c.r.F &= ^(f_H | f_N | f_C)
			a7 := c.r.A >> 7
			c.r.F |= a7
			c.r.A = c.r.A<<1 | a7
			t = 4
		case RRCA:
			c.r.F &= ^(f_H | f_N | f_C)
			a7 := c.r.A & 0x01
			c.r.F |= a7
			c.r.A = c.r.A>>1 | a7<<7
			t = 4
		case RLA:
			fc := c.r.F & f_C
			c.r.F &= ^(f_H | f_N | f_C)
			a7 := c.r.A >> 7
			c.r.A = c.r.A<<1 | fc
			c.r.F |= a7
			t = 4
		case RRA:
			fc := c.r.F & f_C
			c.r.F &= ^(f_H | f_N | f_C)
			a0 := c.r.A & 0x01
			c.r.A = c.r.A>>1 | fc<<7
			c.r.F |= a0
			t = 4
		case DAA:
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
		case EX_AF_AF:
			c.r.A, c.r.A_ = c.r.A_, c.r.A
			c.r.F, c.r.F_ = c.r.F_, c.r.F
			t = 4
		case ADD_A_n, ADD_A_A, ADD_A_B, ADD_A_C, ADD_A_D, ADD_A_E, ADD_A_H, ADD_A_L, ADD_A_HL:
			var n byte
			if opcode == ADD_A_n {
				n = c.readByte()
				t = 7
			} else if opcode == ADD_A_HL {
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
		case ADC_A_A, ADC_A_B, ADC_A_C, ADC_A_D, ADC_A_E, ADC_A_H, ADC_A_L, ADC_A_HL:
			var n byte
			if opcode == ADC_A_HL {
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
		case ADD_HL_BC, ADD_HL_DE, ADD_HL_HL, ADD_HL_SP:
			hl := c.r.getHL()
			var nn word
			switch opcode {
			case ADD_HL_BC:
				nn = c.r.getBC()
			case ADD_HL_DE:
				nn = c.r.getDE()
			case ADD_HL_HL:
				nn = hl
			case ADD_HL_SP:
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
		case SUB_A, SUB_B, SUB_C, SUB_D, SUB_E, SUB_H, SUB_L, SUB_HL:
			a := c.r.A
			var n byte
			if opcode == SUB_HL {
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
		case CP_A, CP_B, CP_C, CP_D, CP_E, CP_H, CP_L, CP_HL:
			var n byte
			if opcode == CP_HL {
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
		case SBC_A_A, SBC_A_B, SBC_A_C, SBC_A_D, SBC_A_E, SBC_A_H, SBC_A_L, SBC_A_HL:
			var n byte
			if opcode == SBC_A_HL {
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
		case AND_A, AND_B, AND_C, AND_D, AND_E, AND_H, AND_L, AND_HL:
			var n byte
			if opcode == AND_HL {
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
		case OR_A, OR_B, OR_C, OR_D, OR_E, OR_H, OR_L, OR_HL:
			var n byte
			if opcode == OR_HL {
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
		case XOR_A, XOR_B, XOR_C, XOR_D, XOR_E, XOR_H, XOR_L, XOR_HL:
			var n byte
			if opcode == XOR_HL {
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
		case LD_A_n, LD_B_n, LD_C_n, LD_D_n, LD_E_n, LD_H_n, LD_L_n:
			r := c.r.getR(opcode & 0b00111000 >> 3)
			*r = c.readByte()
			t = 7
		case
			LD_A_A, LD_A_B, LD_A_C, LD_A_D, LD_A_E, LD_A_H, LD_A_L,
			LD_B_A, LD_B_B, LD_B_C, LD_B_D, LD_B_E, LD_B_H, LD_B_L,
			LD_C_A, LD_C_B, LD_C_C, LD_C_D, LD_C_E, LD_C_H, LD_C_L,
			LD_D_A, LD_D_B, LD_D_C, LD_D_D, LD_D_E, LD_D_H, LD_D_L,
			LD_E_A, LD_E_B, LD_E_C, LD_E_D, LD_E_E, LD_E_H, LD_E_L,
			LD_H_A, LD_H_B, LD_H_C, LD_H_D, LD_H_E, LD_H_H, LD_H_L,
			LD_L_A, LD_L_B, LD_L_C, LD_L_D, LD_L_E, LD_L_H, LD_L_L:
			rs := c.r.getR(opcode & 0b00000111)
			rd := c.r.getR(opcode & 0b00111000 >> 3)
			*rd = *rs
			t = 4
		case LD_BC_nn:
			c.r.C, c.r.B = c.readByte(), c.readByte()
			t = 10
		case LD_DE_nn:
			c.r.E, c.r.D = c.readByte(), c.readByte()
			t = 10
		case LD_HL_nn:
			c.r.L, c.r.H = c.readByte(), c.readByte()
			t = 10
		case LD_SP_nn:
			c.r.SP = word(c.readByte()) | word(c.readByte())<<8
			t = 10
		case LD_HL_mm:
			w := c.readWord()
			c.r.L = c.mem.read(w)
			c.r.H = c.mem.read(w + 1)
			t = 16
		case LD_mm_HL:
			w := c.readWord()
			c.mem.write(w, c.r.L)
			c.mem.write(w+1, c.r.H)
			t = 16
		case LD_mHL_n:
			c.mem.write(c.r.getHL(), c.readByte())
			t = 10
		case LD_mm_A:
			w := c.readWord()
			c.mem.write(w, c.r.A)
			t = 13
		case LD_A_mm:
			w := c.readWord()
			c.r.A = c.mem.read(w)
			t = 13
		case LD_BC_A:
			c.mem.write(c.r.getBC(), c.r.A)
			t = 7
		case LD_DE_A:
			c.mem.write(c.r.getDE(), c.r.A)
			t = 7
		case LD_A_BC:
			c.r.A = c.mem.read(c.r.getBC())
			t = 7
		case LD_A_DE:
			c.r.A = c.mem.read(c.r.getDE())
			t = 7
		case LD_A_HL, LD_B_HL, LD_C_HL, LD_D_HL, LD_E_HL, LD_H_HL, LD_L_HL:
			r := c.r.getR(opcode & 0b00111000 >> 3)
			*r = c.mem.read(c.r.getHL())
			t = 7
		case LD_HL_A, LD_HL_B, LD_HL_C, LD_HL_D, LD_HL_E, LD_HL_H, LD_HL_L:
			r := c.r.getR(opcode & 0b00000111)
			c.mem.write(c.r.getHL(), *r)
			t = 7
		case INC_A, INC_B, INC_C, INC_D, INC_E, INC_H, INC_L:
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
		case INC_BC:
			c.r.setBC(c.r.getBC() + 1)
			t = 6
		case INC_DE:
			c.r.setDE(c.r.getDE() + 1)
			t = 6
		case INC_HL:
			c.r.setHL(c.r.getHL() + 1)
			t = 6
		case INC_SP:
			c.r.SP += 1
			t = 6
		case INC_mHL:
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
		case DEC_A, DEC_B, DEC_C, DEC_D, DEC_E, DEC_H, DEC_L:
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
		case DEC_BC:
			c.r.setBC(c.r.getBC() - 1)
			t = 6
		case DEC_DE:
			c.r.setDE(c.r.getDE() - 1)
			t = 6
		case DEC_HL:
			c.r.setHL(c.r.getHL() - 1)
			t = 6
		case DEC_SP:
			c.r.SP -= 1
			t = 6
		case DEC_mHL:
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
		case JR:
			o := c.readByte()
			if o&0x80 == 0 {
				c.PC += word(o)
			} else {
				c.PC -= word(^o + 1)
			}
			t = 12
		case JR_Z:
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
		case JR_NZ:
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
		case JR_C:
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
		case JR_NC:
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
		case DJNZ:
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
		case JP:
			c.PC = word(c.readByte()) | word(c.readByte())<<8
			t = 10
		case JP_C, JP_M, JP_NC, JP_NZ, JP_P, JP_PE, JP_PO, JP_Z:
			if c.shouldJump(opcode) {
				c.PC = word(c.readByte()) | word(c.readByte())<<8
			}
			t = 10
		case CALL_C, CALL_M, CALL_NC, CALL_NZ, CALL_P, CALL_PE, CALL_PO, CALL_Z:
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
		case RET:
			c.PC = word(c.mem.read(c.r.SP+1))<<8 | word(c.mem.read(c.r.SP))
			c.r.SP += 2
			t = 10
		case RET_C, RET_M, RET_NC, RET_NZ, RET_P, RET_PE, RET_PO, RET_Z:
			if c.shouldJump(opcode) {
				c.PC = word(c.mem.read(c.r.SP+1))<<8 | word(c.mem.read(c.r.SP))
				c.r.SP += 2
				t = 11
			} else {
				t = 5
			}
		case RST_00h, RST_08h, RST_10h, RST_18h, RST_20h, RST_28h, RST_30h, RST_38h:
			c.r.SP -= 1
			c.mem.write(c.r.SP, byte(c.PC>>8))
			c.r.SP -= 1
			c.mem.write(c.r.SP, byte(c.PC))
			c.PC = word(8 * ((opcode & 0b00111000) >> 3))
			t = 11
		case PUSH_AF:
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.A)
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.F)
			t = 11
		case PUSH_BC:
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.B)
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.C)
			t = 11
		case PUSH_DE:
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.D)
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.E)
			t = 11
		case PUSH_HL:
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.H)
			c.r.SP -= 1
			c.mem.write(c.r.SP, c.r.L)
			t = 11
		case POP_AF:
			c.r.A, c.r.F = c.mem.read(c.r.SP+1), c.mem.read(c.r.SP)
			c.r.SP += 2
			t = 10
		case POP_BC:
			c.r.B, c.r.C = c.mem.read(c.r.SP+1), c.mem.read(c.r.SP)
			c.r.SP += 2
			t = 10
		case POP_DE:
			c.r.D, c.r.E = c.mem.read(c.r.SP+1), c.mem.read(c.r.SP)
			c.r.SP += 2
			t = 10
		case POP_HL:
			c.r.H, c.r.L = c.mem.read(c.r.SP+1), c.mem.read(c.r.SP)
			c.r.SP += 2
			t = 10
		case __CB__:
			c.cb(c.readByte(), &t)
		}

		c.wait(t)
	}
}

func (c *CPU) shouldJump(opcode byte) bool {
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

	write := func(v byte) {
		if reg == r_HL {
			c.mem.write(hl, v)
		} else {
			*c.r.regs8[reg] = v
		}
	}

	switch opcode & 0b11111000 {
	case RLC_r:
		cy := v >> 7
		v = v<<1 | cy
		c.r.F = f_NONE
		c.r.F |= f_S & v
		if v == 0 {
			c.r.F |= f_Z
		}
		c.r.F |= parity[v] | cy
		write(v)
	case RRC_r:
		cy := v & f_C
		v = v>>1 | cy<<7
		c.r.F = f_NONE
		c.r.F |= f_S & v
		if v == 0 {
			c.r.F |= f_Z
		}
		c.r.F |= parity[v] | cy
		write(v)
	case RL_r:
	case RR_r:
	case SLA_r:
	case SRA_r:
	case SLL_r:
	case SRL_r:
	default:
		bit := (opcode & 0b00111000) >> 3
		switch opcode & 0b11000000 {
		case BIT_b:
			if reg == r_HL {
				*t = 12
			}
			c.r.F &= ^(f_Z | f_N)
			c.r.F |= f_H
			if v&bit_mask[bit] == 0 {
				c.r.F |= f_Z
			}
		case RES_b:
			write(v & ^bit_mask[bit])
		case SET_b:
			write(v | bit_mask[bit])
		}
	}
}

package z80

type word uint16

type CPU struct {
	PC   word
	mem  *Memory
	r    *registers
	halt bool
}

type Memory struct {
	Cells    []byte
	RAMStart word
}

func NewCPU(mem *Memory) *CPU {
	cpu := &CPU{}
	cpu.mem = mem
	cpu.Reset()
	return cpu
}

func (c *CPU) readByte() byte {
	b := c.mem.Cells[c.PC]
	c.PC += 1
	return b
}

func (c *CPU) readAddr(addr word) byte {
	return c.mem.Cells[addr]
}

func (c *CPU) writeByte(addr word, b byte) {
	c.mem.Cells[addr] = b
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
		case EX_AF_AF:
			c.r.A, c.r.A_ = c.r.A_, c.r.A
			c.r.F, c.r.F_ = c.r.F_, c.r.F
			t = 4
		case ADD_HL_BC, ADD_HL_DE, ADD_HL_HL, ADD_HL_SP:
			hl := c.r.getRR(r_HL)
			rr := c.r.getRR(opcode & 0b00110000 >> 4)
			sum := hl + rr
			c.r.setRRn(r_HL, sum)
			c.r.F &= ^(f_H | f_N | f_C)
			if sum < hl {
				c.r.F |= f_C
			}
			c.r.F |= byte((hl^rr^sum)>>8) & f_H
			t = 11
		case LD_A_n, LD_B_n, LD_C_n, LD_D_n, LD_E_n, LD_H_n, LD_L_n:
			r := c.r.getR(opcode & 0b00111000 >> 3)
			*r = c.readByte()
			t = 7
		case LD_BC_nn, LD_DE_nn, LD_HL_nn, LD_SP_nn:
			c.r.setRRnn(opcode&0b00110000>>4, c.readByte(), c.readByte())
			t = 10
		case LD_BC_A:
			c.writeByte(c.r.getRR(r_BC), c.r.A)
			t = 7
		case LD_A_BC:
			c.r.A = c.readAddr(c.r.getRR(r_BC))
			t = 7
		case INC_A, INC_B, INC_C, INC_D, INC_E, INC_H, INC_L:
			r := c.r.getR(opcode & 0b00111000 >> 3)
			c.r.F &= ^(f_S | f_Z | f_H | f_PV | f_N)
			if *r == 0x7F {
				c.r.F |= f_PV
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
		case DEC_BC, DEC_DE, DEC_HL, DEC_SP:
			reg := opcode & 0b00110000 >> 4
			c.r.setRRn(reg, c.r.getRR(reg)-1)
			t = 6
		case JR:
			o := c.readByte()
			if o&0x80 == 0 {
				c.PC += word(o)
			} else {
				c.PC -= word(^o + 1)
			}
			t = 12
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
		}

		c.wait(t)
	}
}

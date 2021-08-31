package z80

import (
	"fmt"

	"github.com/voytas/z80-go-zx/z80/memory"
)

type IOBus interface {
	Read(hi, lo byte) byte
	Write(hi, lo, data byte)
}

// Represents emulated Z80 Z80
type Z80 struct {
	IOBus            IOBus
	mem              memory.Memory // memory
	Reg              *registers    // registers
	halt, iff1, iff2 bool          // states of halt, iff1 and iff2
	im               byte          // interrupt mode (im0, im1 or in2)
	TC               *TCounter     // T states counter
	Trap             func()        // traps to execute on PC address
}

// Creates a new instance of the Z80 emulator.
func NewZ80(mem memory.Memory) *Z80 {
	z80 := &Z80{}
	z80.mem = mem
	z80.Reset()
	return z80
}

// Fetches the opcode and increments PC afterwards. The cost is 4T.
func (z80 *Z80) fetch() byte {
	b := z80.mem.Read(z80.Reg.PC)
	z80.Reg.PC += 1
	z80.TC.Add(4)
	return b
}

// Reads 8 bit value from memory location specified by the current PC value
// and increments PC afterwards. The cost is 3T.
func (z80 *Z80) nextByte() byte {
	b := z80.mem.Read(z80.Reg.PC)
	z80.Reg.PC += 1
	z80.TC.Add(3)
	return b
}

// Reads 16 bit value from memory location specified by the current PC value
// and increments PC afterwards. The cost is 2 * 3T.
func (z80 *Z80) nextWord() uint16 {
	return uint16(z80.nextByte()) | uint16(z80.nextByte())<<8
}

// Reads 8 bit value from the memory address. Does not affect PC. The cost is 3T.
func (z80 *Z80) read(addr uint16) byte {
	b := z80.mem.Read(addr)
	z80.TC.Add(3)
	return b
}

// Writes 8 bit value to the memory address. The cost is 3T.
func (z80 *Z80) write(addr uint16, value byte) {
	z80.mem.Write(addr, value)
	z80.TC.Add(3)
}

// Reads 8 bit value from the bus (IN port). The cost is 4T.
func (z80 *Z80) readBus(hi, lo byte) byte {
	b := byte(0xFF)
	if z80.IOBus != nil {
		b = z80.IOBus.Read(hi, lo)
	} else {
		z80.TC.Add(4)
	}
	return b
}

// Writes 8 bit value to the bus (OUT port). The cost is 4T.
func (z80 *Z80) writeBus(hi, lo, data byte) {
	if z80.IOBus != nil {
		z80.IOBus.Write(hi, lo, data)
	} else {
		z80.TC.Add(4)
	}
}

// Writes PC to stack. The cost is 2 * 3T.
func (z80 *Z80) pushPC() {
	z80.Reg.SP -= 1
	z80.write(z80.Reg.SP, byte(z80.Reg.PC>>8))
	z80.Reg.SP -= 1
	z80.write(z80.Reg.SP, byte(z80.Reg.PC))
}

// Emulates maskable interrupt (INT)
func (z80 *Z80) INT(data byte) {
	z80.halt = false
	if !z80.iff1 {
		return
	}
	z80.iff1, z80.iff2 = false, false
	switch z80.im {
	case 0, 1:
		z80.pushPC()
		z80.Reg.PC = 0x38 // RST 38h
		z80.TC.Add(7)
	case 2:
		z80.pushPC()
		addr := uint16(z80.Reg.I)<<8 + uint16(data)
		z80.Reg.PC = uint16(z80.read(addr+1))<<8 | uint16(z80.read(addr))
		z80.TC.Add(7)
	}
	z80.Reg.IncR()
}

// Emulates non-maskable interrupt (NMI)
func (z80 *Z80) NMI() {
	z80.halt = false
	z80.iff2, z80.iff1 = z80.iff1, false
	z80.pushPC()
	z80.Reg.PC = 0x66
	z80.TC.Add(5)
	z80.Reg.IncR()
}

func (z80 *Z80) Reset() {
	z80.Reg = newRegisters()
	z80.Reg.PC, z80.Reg.SP = 0, 0xFFFF
	z80.halt = false
	z80.iff1, z80.iff2 = false, false
	z80.Reg.A, z80.Reg.F = 0xFF, 0xFF
	z80.Reg.I, z80.Reg.R = 0x00, 0x00
	z80.halt = false
	z80.TC = &TCounter{}
}

// Executes the instructions until maximum number of T states is reached.
// (limit equal to 0 specifies unlimited number of T states to execute)
func (z80 *Z80) Run(limit int) {
	z80.TC.limit(limit)
	for !(z80.Reg.prefix == noPrefix && z80.TC.done()) {
		if z80.Trap != nil {
			z80.Trap()
		}

		var opcode byte
		if z80.halt {
			z80.TC.halt()
			break
		} else {
			opcode = z80.fetch()
		}

		z80.Reg.IncR()

		switch opcode {
		case nop:
		case halt:
			z80.Reg.prefix = noPrefix
			z80.halt = true
		case di:
			z80.iff1, z80.iff2 = false, false
		case ei:
			z80.iff1, z80.iff2 = true, true
		case rlca:
			a7 := z80.Reg.A >> 7
			z80.Reg.A = z80.Reg.A<<1 | a7
			z80.Reg.F = z80.Reg.F&(FS|FZ|FP) | z80.Reg.A&(FY|FX) | a7
		case rrca:
			a0 := z80.Reg.A & 0x01
			z80.Reg.A = z80.Reg.A>>1 | a0<<7
			z80.Reg.F = z80.Reg.F&(FS|FZ|FP) | z80.Reg.A&(FY|FX) | a0
		case rla:
			a7 := z80.Reg.A >> 7
			z80.Reg.A = z80.Reg.A<<1 | z80.Reg.F&FC
			z80.Reg.F = z80.Reg.F&(FS|FZ|FP) | z80.Reg.A&(FY|FX) | a7
		case rra:
			a0 := z80.Reg.A & 0x01
			z80.Reg.A = z80.Reg.A>>1 | z80.Reg.F&FC<<7
			z80.Reg.F = z80.Reg.F&(FS|FZ|FP) | z80.Reg.A&(FY|FX) | a0
		case cpl:
			z80.Reg.A = ^z80.Reg.A
			z80.Reg.F = z80.Reg.F&(FS|FZ|FP|FC) | FH | FN | z80.Reg.A&(FY|FX)
		case scf:
			z80.Reg.F = z80.Reg.F&(FS|FZ|FP) | FC | z80.Reg.A&(FY|FX)
		case ccf:
			z80.Reg.F = (z80.Reg.F&(FS|FZ|FP|FC) | z80.Reg.F&FC<<4 | z80.Reg.A&(FY|FX)) ^ FC
		case daa:
			cf := z80.Reg.F & FC
			hf := z80.Reg.F & FH
			nf := z80.Reg.F & FN
			lb := z80.Reg.A & 0x0F
			diff := byte(0)
			z80.Reg.F &= FN
			if cf != 0 || z80.Reg.A > 0x99 {
				diff = 0x60
				z80.Reg.F |= FC
			}
			if hf != 0 || lb > 0x09 {
				diff += 0x06
			}
			if nf == 0 {
				z80.Reg.A += diff
			} else {
				z80.Reg.A -= diff
			}
			z80.Reg.F |= parity[z80.Reg.A] | z80.Reg.A&(FS|FY|FX)
			if z80.Reg.A == 0 {
				z80.Reg.F |= FZ
			}
			if nf == 0 && lb > 0x09 || nf != 0 && hf != 0 && lb < 0x06 {
				z80.Reg.F |= FH
			}
		case ex_af_af:
			z80.Reg.A, z80.Reg.A_ = z80.Reg.A_, z80.Reg.A
			z80.Reg.F, z80.Reg.F_ = z80.Reg.F_, z80.Reg.F
		case exx:
			z80.Reg.B, z80.Reg.B_, z80.Reg.C, z80.Reg.C_ = z80.Reg.B_, z80.Reg.B, z80.Reg.C_, z80.Reg.C
			z80.Reg.D, z80.Reg.D_, z80.Reg.E, z80.Reg.E_ = z80.Reg.D_, z80.Reg.D, z80.Reg.E_, z80.Reg.E
			z80.Reg.H, z80.Reg.H_, z80.Reg.L, z80.Reg.L_ = z80.Reg.H_, z80.Reg.H, z80.Reg.L_, z80.Reg.L
		case ex_de_hl:
			z80.Reg.D, z80.Reg.E, z80.Reg.H, z80.Reg.L = z80.Reg.H, z80.Reg.L, z80.Reg.D, z80.Reg.E
		case ex_sp_hl:
			h, l := z80.Reg.r(rH), z80.Reg.r(rL)
			y, x := z80.read(z80.Reg.SP), z80.read(z80.Reg.SP+1)
			z80.delay(1, z80.Reg.SP+1)
			z80.write(z80.Reg.SP, *l)
			z80.write(z80.Reg.SP+1, *h)
			z80.delay(2, z80.Reg.SP+1)
			*h, *l = x, y
		case add_a_n, add_a_a, add_a_b, add_a_c, add_a_d, add_a_e, add_a_h, add_a_l, add_a_hl:
			a := z80.Reg.A
			var n byte
			switch opcode {
			case add_a_n:
				n = z80.nextByte()
			case add_a_hl:
				hl := z80.getHL()
				if z80.Reg.prefix != noPrefix {
					z80.delay(5, z80.Reg.PC-1)
				}
				n = z80.read(hl)
			default:
				n = *z80.Reg.r(opcode & 0b00000111)
			}
			z80.Reg.A += n
			z80.Reg.F = (FS | FY | FX) & z80.Reg.A
			if z80.Reg.A == 0 {
				z80.Reg.F |= FZ
			}
			z80.Reg.F |= (a ^ n ^ z80.Reg.A) & FH
			if (a^n)&0x80 == 0 && (a^z80.Reg.A)&0x80 != 0 {
				z80.Reg.F |= FP
			}
			if z80.Reg.A < a {
				z80.Reg.F |= FC
			}
		case adc_a_n, adc_a_a, adc_a_b, adc_a_c, adc_a_d, adc_a_e, adc_a_h, adc_a_l, adc_a_hl:
			var n byte
			switch opcode {
			case adc_a_n:
				n = z80.nextByte()
			case adc_a_hl:
				hl := z80.getHL()
				if z80.Reg.prefix != noPrefix {
					z80.delay(5, z80.Reg.PC-1)
				}
				n = z80.read(hl)
			default:
				n = *z80.Reg.r(opcode & 0b00000111)
			}
			cf := z80.Reg.F & FC
			sum_w := uint16(z80.Reg.A) + uint16(n) + uint16(cf)
			sum_b := byte(sum_w)
			z80.Reg.F = (FS | FY | FX) & sum_b
			if sum_b == 0 {
				z80.Reg.F |= FZ
			}
			z80.Reg.F |= (z80.Reg.A ^ n ^ sum_b) & FH
			if (z80.Reg.A^n)&0x80 == 0 && (z80.Reg.A^sum_b)&0x80 != 0 {
				z80.Reg.F |= FP
			}
			if sum_w > 0xff {
				z80.Reg.F |= FC
			}
			z80.Reg.A = sum_b
		case add_hl_bc, add_hl_de, add_hl_hl, add_hl_sp:
			z80.delay(7, z80.Reg.IR())
			hl := z80.Reg.HL()
			var nn uint16
			switch opcode {
			case add_hl_bc:
				nn = z80.Reg.BC()
			case add_hl_de:
				nn = z80.Reg.DE()
			case add_hl_hl:
				nn = hl
			case add_hl_sp:
				nn = z80.Reg.SP
			}
			sum := hl + nn
			z80.Reg.SetHL(sum)
			z80.Reg.F = z80.Reg.F & ^(FH|FN|FC) | byte((hl^nn^sum)>>8)&FH | byte(sum>>8)&(FY|FX)
			if sum < hl {
				z80.Reg.F |= FC
			}
		case sub_a, sub_b, sub_c, sub_d, sub_e, sub_h, sub_l, sub_hl, sub_n:
			a := z80.Reg.A
			var n byte
			switch opcode {
			case sub_n:
				n = z80.nextByte()
			case sub_hl:
				hl := z80.getHL()
				if z80.Reg.prefix != noPrefix {
					z80.delay(5, z80.Reg.PC-1)
				}
				n = z80.read(hl)
			default:
				n = *z80.Reg.r(opcode & 0b00000111)
			}
			z80.Reg.A -= n
			z80.Reg.F = (FS|FY|FX)&z80.Reg.A | FN | (a^n^z80.Reg.A)&FH
			if z80.Reg.A == 0 {
				z80.Reg.F |= FZ
			}
			if (a^n)&0x80 != 0 && (a^z80.Reg.A)&0x80 != 0 {
				z80.Reg.F |= FP
			}
			if z80.Reg.A > a {
				z80.Reg.F |= FC
			}
		case cp_a, cp_b, cp_c, cp_d, cp_e, cp_h, cp_l, cp_hl, cp_n:
			var n byte
			switch opcode {
			case cp_n:
				n = z80.nextByte()
			case cp_hl:
				hl := z80.getHL()
				if z80.Reg.prefix != noPrefix {
					z80.delay(5, z80.Reg.PC-1)
				}
				n = z80.read(hl)
			default:
				n = *z80.Reg.r(opcode & 0b00000111)
			}
			test := z80.Reg.A - n
			z80.Reg.F = FN | FS&test | n&(FY|FX) | byte(z80.Reg.A^n^test)&FH
			if test == 0 {
				z80.Reg.F |= FZ
			}
			if (z80.Reg.A^n)&0x80 != 0 && (z80.Reg.A^test)&0x80 != 0 {
				z80.Reg.F |= FP
			}
			if test > z80.Reg.A {
				z80.Reg.F |= FC
			}
		case sbc_a_a, sbc_a_b, sbc_a_c, sbc_a_d, sbc_a_e, sbc_a_h, sbc_a_l, sbc_a_hl, sbc_a_n:
			var n byte
			switch opcode {
			case sbc_a_n:
				n = z80.nextByte()
			case sbc_a_hl:
				hl := z80.getHL()
				if z80.Reg.prefix != noPrefix {
					z80.delay(5, z80.Reg.PC-1)
				}
				n = z80.read(hl)
			default:
				n = *z80.Reg.r(opcode & 0b00000111)
			}
			cf := z80.Reg.F & FC
			sub_w := uint16(z80.Reg.A) - uint16(n) - uint16(cf)
			sub_b := byte(sub_w)
			z80.Reg.F = (FS|FY|FX)&sub_b | FN | byte(z80.Reg.A^n^sub_b)&FH
			if sub_b == 0 {
				z80.Reg.F |= FZ
			}
			if (z80.Reg.A^n)&0x80 != 0 && (sub_b^z80.Reg.A)&0x80 != 0 {
				z80.Reg.F |= FP
			}
			if sub_w > 0xff {
				z80.Reg.F |= FC
			}
			z80.Reg.A = sub_b
		case and_a, and_b, and_c, and_d, and_e, and_h, and_l, and_hl, and_n:
			var n byte
			switch opcode {
			case and_n:
				n = z80.nextByte()
			case and_hl:
				hl := z80.getHL()
				if z80.Reg.prefix != noPrefix {
					z80.delay(5, z80.Reg.PC-1)
				}
				n = z80.read(hl)
			default:
				n = *z80.Reg.r(opcode & 0b00000111)
			}
			z80.Reg.A &= n
			z80.Reg.F = (FS|FY|FX)&z80.Reg.A | FH | parity[z80.Reg.A]
			if z80.Reg.A == 0 {
				z80.Reg.F |= FZ
			}
		case or_a, or_b, or_c, or_d, or_e, or_h, or_l, or_hl, or_n:
			var n byte
			switch opcode {
			case or_n:
				n = z80.nextByte()
			case or_hl:
				hl := z80.getHL()
				if z80.Reg.prefix != noPrefix {
					z80.delay(5, z80.Reg.PC-1)
				}
				n = z80.read(hl)
			default:
				n = *z80.Reg.r(opcode & 0b00000111)
			}
			z80.Reg.A |= n
			z80.Reg.F = (FS|FY|FX)&z80.Reg.A | parity[z80.Reg.A]
			if z80.Reg.A == 0 {
				z80.Reg.F |= FZ
			}
		case xor_a, xor_b, xor_c, xor_d, xor_e, xor_h, xor_l, xor_hl, xor_n:
			var n byte
			switch opcode {
			case xor_n:
				n = z80.nextByte()
			case xor_hl:
				hl := z80.getHL()
				if z80.Reg.prefix != noPrefix {
					z80.delay(5, z80.Reg.PC-1)
				}
				n = z80.read(hl)
			default:
				n = *z80.Reg.r(opcode & 0b00000111)
			}
			z80.Reg.A ^= n
			z80.Reg.F = (FS|FY|FX)&z80.Reg.A | parity[z80.Reg.A]
			if z80.Reg.A == 0 {
				z80.Reg.F |= FZ
			}
		case ld_a_n, ld_b_n, ld_c_n, ld_d_n, ld_e_n, ld_h_n, ld_l_n:
			r := z80.Reg.r(opcode & 0b00111000 >> 3)
			*r = z80.nextByte()
		case
			ld_a_a, ld_a_b, ld_a_c, ld_a_d, ld_a_e, ld_a_h, ld_a_l,
			ld_b_a, ld_b_b, ld_b_c, ld_b_d, ld_b_e, ld_b_h, ld_b_l,
			ld_c_a, ld_c_b, ld_c_c, ld_c_d, ld_c_e, ld_c_h, ld_c_l,
			ld_d_a, ld_d_b, ld_d_c, ld_d_d, ld_d_e, ld_d_h, ld_d_l,
			ld_e_a, ld_e_b, ld_e_c, ld_e_d, ld_e_e, ld_e_h, ld_e_l,
			ld_h_a, ld_h_b, ld_h_c, ld_h_d, ld_h_e, ld_h_h, ld_h_l,
			ld_l_a, ld_l_b, ld_l_c, ld_l_d, ld_l_e, ld_l_h, ld_l_l:
			rs := z80.Reg.r(opcode & 0b00000111)
			rd := z80.Reg.r(opcode & 0b00111000 >> 3)
			*rd = *rs
		case ld_bc_nn:
			z80.Reg.C, z80.Reg.B = z80.nextByte(), z80.nextByte()
		case ld_de_nn:
			z80.Reg.E, z80.Reg.D = z80.nextByte(), z80.nextByte()
		case ld_hl_nn:
			h, l := z80.Reg.r(rH), z80.Reg.r(rL)
			*l, *h = z80.nextByte(), z80.nextByte()
		case ld_sp_nn:
			z80.Reg.SP = z80.nextWord()
		case ld_sp_hl:
			z80.delay(2, z80.Reg.IR())
			z80.Reg.SP = z80.Reg.HL()
		case ld_hl_mm:
			addr := z80.nextWord()
			h, l := z80.Reg.r(rH), z80.Reg.r(rL)
			*l = z80.read(addr)
			*h = z80.read(addr + 1)
		case ld_mm_hl:
			addr := z80.nextWord()
			h, l := z80.Reg.r(rH), z80.Reg.r(rL)
			z80.write(addr, *l)
			z80.write(addr+1, *h)
		case ld_mhl_n:
			hl := z80.getHL()
			n := z80.nextByte()
			if z80.Reg.prefix != noPrefix {
				z80.delay(2, z80.Reg.PC-1)
			}
			z80.write(hl, n)
		case ld_mm_a:
			z80.write(z80.nextWord(), z80.Reg.A)
		case ld_a_mm:
			z80.Reg.A = z80.read(z80.nextWord())
		case ld_bc_a:
			z80.write(z80.Reg.BC(), z80.Reg.A)
		case ld_de_a:
			z80.write(z80.Reg.DE(), z80.Reg.A)
		case ld_a_bc:
			z80.Reg.A = z80.read(z80.Reg.BC())
		case ld_a_de:
			z80.Reg.A = z80.read(z80.Reg.DE())
		case ld_a_hl, ld_b_hl, ld_c_hl, ld_d_hl, ld_e_hl, ld_h_hl, ld_l_hl:
			hl := z80.getHL()
			if z80.Reg.prefix != noPrefix {
				z80.delay(5, z80.Reg.PC-1)
			}
			*z80.Reg.raw[opcode&0b00111000>>3] = z80.read(hl)
		case ld_hl_a, ld_hl_b, ld_hl_c, ld_hl_d, ld_hl_e, ld_hl_h, ld_hl_l:
			hl := z80.getHL()
			if z80.Reg.prefix != noPrefix {
				z80.delay(5, z80.Reg.PC-1)
			}
			z80.write(hl, *z80.Reg.raw[opcode&0b00000111])
		case inc_a, inc_b, inc_c, inc_d, inc_e, inc_h, inc_l:
			r := z80.Reg.r(opcode & 0b00111000 >> 3)
			z80.Reg.F &= FC
			if *r == 0x7F {
				z80.Reg.F |= FP
			}
			if *r&0x0F == 0x0F {
				z80.Reg.F |= FH
			}
			*r += 1
			z80.Reg.F |= *r & (FS | FY | FX)
			if *r == 0 {
				z80.Reg.F |= FZ
			}
		case inc_bc:
			z80.delay(2, z80.Reg.IR())
			z80.Reg.SetBC(z80.Reg.BC() + 1)
		case inc_de:
			z80.delay(2, z80.Reg.IR())
			z80.Reg.SetDE(z80.Reg.DE() + 1)
		case inc_hl:
			z80.delay(2, z80.Reg.IR())
			z80.Reg.SetHL(z80.Reg.HL() + 1)
		case inc_sp:
			z80.delay(2, z80.Reg.IR())
			z80.Reg.SP += 1
		case inc_mhl:
			hl := z80.getHL()
			if z80.Reg.prefix != noPrefix {
				z80.delay(5, z80.Reg.PC-1)
			}
			b := z80.read(hl)
			z80.delay(1, hl)
			z80.Reg.F &= FC
			if b == 0x7F {
				z80.Reg.F |= FP
			}
			if b&0x0F == 0x0F {
				z80.Reg.F |= FH
			}
			b += 1
			if b == 0x00 {
				z80.Reg.F |= FZ
			}
			z80.Reg.F |= b & (FS | FY | FX)
			z80.write(hl, b)
		case dec_a, dec_b, dec_c, dec_d, dec_e, dec_h, dec_l:
			r := z80.Reg.r(opcode & 0b00111000 >> 3)
			z80.Reg.F = z80.Reg.F&FC | FN
			if *r == 0x80 {
				z80.Reg.F |= FP
			}
			if *r&0x0F == 0 {
				z80.Reg.F |= FH
			}
			*r -= 1
			z80.Reg.F |= *r & (FS | FY | FX)
			if *r == 0 {
				z80.Reg.F |= FZ
			}
		case dec_bc:
			z80.delay(2, z80.Reg.IR())
			z80.Reg.SetBC(z80.Reg.BC() - 1)
		case dec_de:
			z80.delay(2, z80.Reg.IR())
			z80.Reg.SetDE(z80.Reg.DE() - 1)
		case dec_hl:
			z80.delay(2, z80.Reg.IR())
			z80.Reg.SetHL(z80.Reg.HL() - 1)
		case dec_sp:
			z80.delay(2, z80.Reg.IR())
			z80.Reg.SP -= 1
		case dec_mhl:
			hl := z80.getHL()
			if z80.Reg.prefix != noPrefix {
				z80.delay(5, z80.Reg.PC-1)
			}
			b := z80.read(hl)
			z80.delay(1, hl)
			z80.Reg.F = z80.Reg.F&FC | FN
			if b == 0x80 {
				z80.Reg.F |= FP
			}
			if b&0x0F == 0 {
				z80.Reg.F |= FH
			}
			b -= 1
			if b == 0x00 {
				z80.Reg.F |= FZ
			}
			z80.Reg.F |= b & (FS | FY | FX)
			z80.write(hl, b)
		case jr_o:
			o := z80.nextByte()
			z80.delay(5, z80.Reg.PC-1)
			if o&0x80 == 0 {
				z80.Reg.PC += uint16(o)
			} else {
				z80.Reg.PC -= uint16(^o + 1)
			}
		case jr_z_o:
			o := z80.nextByte()
			if z80.Reg.F&FZ == FZ {
				z80.delay(5, z80.Reg.PC-1)
				if o&0x80 == 0 {
					z80.Reg.PC += uint16(o)
				} else {
					z80.Reg.PC -= uint16(^o + 1)
				}
			}
		case jr_nz_o:
			o := z80.nextByte()
			if z80.Reg.F&FZ == 0 {
				z80.delay(5, z80.Reg.PC-1)
				if o&0x80 == 0 {
					z80.Reg.PC += uint16(o)
				} else {
					z80.Reg.PC -= uint16(^o + 1)
				}
			}
		case jr_c:
			o := z80.nextByte()
			if z80.Reg.F&FC == FC {
				z80.delay(5, z80.Reg.PC-1)
				if o&0x80 == 0 {
					z80.Reg.PC += uint16(o)
				} else {
					z80.Reg.PC -= uint16(^o + 1)
				}
			}
		case jr_nc_o:
			o := z80.nextByte()
			if z80.Reg.F&FC == 0 {
				z80.delay(5, z80.Reg.PC-1)
				if o&0x80 == 0 {
					z80.Reg.PC += uint16(o)
				} else {
					z80.Reg.PC -= uint16(^o + 1)
				}
			}
		case djnz:
			z80.delay(1, z80.Reg.IR())
			o := z80.nextByte()
			z80.Reg.B -= 1
			if z80.Reg.B != 0 {
				z80.delay(5, z80.Reg.PC-1)
				if o&0x80 == 0 {
					z80.Reg.PC += uint16(o)
				} else {
					z80.Reg.PC -= uint16(^o + 1)
				}
			}
		case jp_nn:
			z80.Reg.PC = z80.nextWord()
		case jp_c_nn, jp_m_nn, jp_nc_nn, jp_nz_nn, jp_p_nn, jp_pe_nn, jp_po_nn, jp_z_nn:
			pc := z80.nextWord()
			if z80.shouldJump(opcode) {
				z80.Reg.PC = pc
			}
		case jp_hl:
			z80.Reg.PC = z80.Reg.HL()
		case call_nn:
			pc := z80.nextWord()
			z80.delay(1, z80.Reg.PC-1)
			z80.Reg.SP -= 1
			z80.write(z80.Reg.SP, byte(z80.Reg.PC>>8))
			z80.Reg.SP -= 1
			z80.write(z80.Reg.SP, byte(z80.Reg.PC))
			z80.Reg.PC = pc
		case call_c_nn, call_m_nn, call_nc_nn, call_nz_nn, call_p_nn, call_pe_nn, call_po_nn, call_z_nn:
			pc := z80.nextWord()
			if z80.shouldJump(opcode) {
				z80.delay(1, z80.Reg.PC-1)
				z80.Reg.SP -= 1
				z80.write(z80.Reg.SP, byte(z80.Reg.PC>>8))
				z80.Reg.SP -= 1
				z80.write(z80.Reg.SP, byte(z80.Reg.PC))
				z80.Reg.PC = pc
			}
		case ret:
			z80.Reg.PC = uint16(z80.read(z80.Reg.SP+1))<<8 | uint16(z80.read(z80.Reg.SP))
			z80.Reg.SP += 2
		case ret_c, ret_m, ret_nc, ret_nz, ret_p, ret_pe, ret_po, ret_z:
			z80.delay(1, z80.Reg.IR())
			if z80.shouldJump(opcode) {
				z80.Reg.PC = uint16(z80.read(z80.Reg.SP+1))<<8 | uint16(z80.read(z80.Reg.SP))
				z80.Reg.SP += 2
			}
		case rst_00h, rst_08h, rst_10h, rst_18h, rst_20h, rst_28h, rst_30h, rst_38h:
			z80.delay(1, z80.Reg.IR())
			z80.Reg.SP -= 1
			z80.write(z80.Reg.SP, byte(z80.Reg.PC>>8))
			z80.Reg.SP -= 1
			z80.write(z80.Reg.SP, byte(z80.Reg.PC))
			z80.Reg.PC = uint16(8 * ((opcode & 0b00111000) >> 3))
		case push_af:
			z80.delay(1, z80.Reg.IR())
			z80.Reg.SP -= 1
			z80.write(z80.Reg.SP, z80.Reg.A)
			z80.Reg.SP -= 1
			z80.write(z80.Reg.SP, z80.Reg.F)
		case push_bc:
			z80.delay(1, z80.Reg.IR())
			z80.Reg.SP -= 1
			z80.write(z80.Reg.SP, z80.Reg.B)
			z80.Reg.SP -= 1
			z80.write(z80.Reg.SP, z80.Reg.C)
		case push_de:
			z80.delay(1, z80.Reg.IR())
			z80.Reg.SP -= 1
			z80.write(z80.Reg.SP, z80.Reg.D)
			z80.Reg.SP -= 1
			z80.write(z80.Reg.SP, z80.Reg.E)
		case push_hl:
			z80.delay(1, z80.Reg.IR())
			z80.Reg.SP -= 1
			z80.write(z80.Reg.SP, *z80.Reg.r(rH))
			z80.Reg.SP -= 1
			z80.write(z80.Reg.SP, *z80.Reg.r(rL))
		case pop_af:
			z80.Reg.A, z80.Reg.F = z80.read(z80.Reg.SP+1), z80.read(z80.Reg.SP)
			z80.Reg.SP += 2
		case pop_bc:
			z80.Reg.B, z80.Reg.C = z80.read(z80.Reg.SP+1), z80.read(z80.Reg.SP)
			z80.Reg.SP += 2
		case pop_de:
			z80.Reg.D, z80.Reg.E = z80.read(z80.Reg.SP+1), z80.read(z80.Reg.SP)
			z80.Reg.SP += 2
		case pop_hl:
			*z80.Reg.r(rH), *z80.Reg.r(rL) = z80.read(z80.Reg.SP+1), z80.read(z80.Reg.SP)
			z80.Reg.SP += 2
		case in_a_n:
			z80.Reg.A = z80.readBus(z80.Reg.A, z80.nextByte())
		case out_n_a:
			z80.writeBus(z80.Reg.A, z80.nextByte(), z80.Reg.A)
		case prefix_cb:
			z80.Reg.IncR()
			z80.prefixCB()
		case prefix_ed:
			z80.Reg.IncR()
			z80.prefixED(z80.fetch())
		case useIX:
			z80.Reg.IncR()
			z80.Reg.prefix = useIX
			continue
		case useIY:
			z80.Reg.IncR()
			z80.Reg.prefix = useIY
			continue
		}
		z80.Reg.prefix = noPrefix
		//log.Println(fmt.Sprintf("OP: %X T: %v", opcode, z80.TC.Current))
	}
}

func (z80 *Z80) shouldJump(opcode byte) bool {
	switch opcode & 0b00111000 {
	case 0b00000000: // Non-Zero (NZ)
		return z80.Reg.F&FZ == 0
	case 0b00001000: // Zero (Z)
		return z80.Reg.F&FZ != 0
	case 0b00010000: // Non Carry (NC)
		return z80.Reg.F&FC == 0
	case 0b00011000: // Carry (C)
		return z80.Reg.F&FC != 0
	case 0b00100000: // Parity Odd (PO)
		return z80.Reg.F&FP == 0
	case 0b00101000: // Parity Even (PE)
		return z80.Reg.F&FP != 0
	case 0b00110000: // Sign Positive (P)
		return z80.Reg.F&FS == 0
	case 0b00111000: // Sign Negative (M)
		return z80.Reg.F&FS != 0
	}

	panic(fmt.Sprintf("Invalid opcode %v", opcode))
}

// Returns value of HL / (IX + d) / (IY + d) register. The current prefix
// determines whether to use IX or IY register instead of HL.
func (z80 *Z80) getHL() uint16 {
	if z80.Reg.prefix != noPrefix {
		return z80.getHLOffset(z80.nextByte())
	}
	return z80.Reg.HL()
}

// For IX or IY add offset to register value, otherwise return HL.
func (z80 *Z80) getHLOffset(offset byte) uint16 {
	hl := z80.Reg.HL()
	if offset == 0 {
		return hl
	} else if offset&0x80 == 0 {
		return hl + uint16(offset)
	} else {
		return hl - uint16(^offset+1)
	}
}

// Handles memory read contention which may require repeated memory access
func (z80 *Z80) delay(count int, addr uint16) {
	for i := 0; i < count; i++ {
		z80.mem.Read(addr)
		z80.TC.Add(1)
	}
}

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
	reg              *registers    // registers
	halt, iff1, iff2 bool          // states of halt, iff1 and iff2
	im               byte          // interrupt mode (im0, im1 or in2)
	TC               *TCounter     // T states counter
}

func NewZ80(mem memory.Memory) *Z80 {
	z80 := &Z80{}
	z80.mem = mem
	z80.Reset()
	return z80
}

// Fetches the opcode and increments PC afterwards. The cost is 4T.
func (z80 *Z80) fetch() byte {
	z80.TC.Add(4)
	b := z80.mem.Read(z80.reg.PC)
	z80.reg.PC += 1
	return b
}

// Reads 8 bit value from memory location specified by current PC value
// and increments PC afterwards. The cost is 3T.
func (z80 *Z80) readByte() byte {
	z80.TC.Add(3)
	b := z80.mem.Read(z80.reg.PC)
	z80.reg.PC += 1
	return b
}

// Reads 8 bit value from memory address. Does not affect PC. The cost is 3T.
func (z80 *Z80) read(addr uint16) byte {
	z80.TC.Add(3)
	b := z80.mem.Read(addr)
	return b
}

// Reads 8 bit value from the bus (IN port).
func (z80 *Z80) readBus(hi, lo byte) byte {
	z80.TC.Add(4)
	if z80.IOBus != nil {
		return z80.IOBus.Read(hi, lo)
	}

	return 0xFF
}

func (z80 *Z80) contention(addr uint16, t int) {
	z80.TC.Add(t)
}

// Writes 8 bit value to memory address. The cost is 3T.
func (z80 *Z80) write(addr uint16, value byte) {
	z80.TC.Add(3)
	z80.mem.Write(addr, value)
}

// Writes 8 bit value to the bus (OUT port).
func (z80 *Z80) writeBus(hi, lo, data byte) {
	z80.TC.Add(4)
	if z80.IOBus != nil {
		z80.IOBus.Write(hi, lo, data)
	}
}

// TODO: pass address
// Reads 16 bit value from memory. The cost is 2*3T.
func (z80 *Z80) readWord() uint16 {
	z80.TC.Add(2 * 3)
	w := uint16(z80.mem.Read(z80.reg.PC)) | uint16(z80.mem.Read(z80.reg.PC+1))<<8
	z80.reg.PC += 2
	return w
}

func (z80 *Z80) pushPC() {
	z80.reg.SP -= 1
	z80.mem.Write(z80.reg.SP, byte(z80.reg.PC>>8))
	z80.reg.SP -= 1
	z80.mem.Write(z80.reg.SP, byte(z80.reg.PC))
}

func (z80 *Z80) incR() {
	z80.reg.R = z80.reg.R&0x80 | (z80.reg.R+1)&0x7F
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
		z80.reg.PC = 0x38 // RST 38h
		z80.TC.Add(13)
	case 2:
		z80.pushPC()
		addr := uint16(z80.reg.I)<<8 + uint16(data)
		z80.reg.PC = uint16(z80.mem.Read(addr+1))<<8 | uint16(z80.mem.Read(addr))
		z80.TC.Add(19)
	}
	z80.incR()
}

// Emulates non-maskable interrupt (NMI)
func (z80 *Z80) NMI() {
	z80.halt = false
	z80.iff2, z80.iff1 = z80.iff1, false
	z80.pushPC()
	z80.reg.PC = 0x66
	z80.TC.Add(11)
	z80.incR()
}

func (z80 *Z80) Reset() {
	z80.reg = newRegisters()
	z80.reg.PC, z80.reg.SP = 0, 0xFFFF
	z80.halt = false
	z80.iff1, z80.iff2 = false, false
	z80.reg.A, z80.reg.F = 0xFF, 0xFF
	z80.reg.I, z80.reg.R = 0x00, 0x00
	z80.halt = false
	z80.TC = &TCounter{}
}

// Executes the instructions until maximum number of T states is reached.
// (tLimit equal to 0 specifies unlimited number of T states to execute)
func (z80 *Z80) Run(limit int) {
	// Update limit with remaining from previous run
	z80.TC.limit(limit)
	for {
		if z80.TC.done() {
			break
		}

		var opcode byte
		if z80.halt {
			z80.TC.halt()
			break
		} else {
			opcode = z80.fetch()
		}

		// debugger.Debug(opcode, z80.reg.prefix, z80.reg.PC, z80.mem)
		z80.incR()

		switch opcode {
		case nop:
		case halt:
			z80.reg.prefix = noPrefix
			z80.halt = true
		case di:
			z80.iff1, z80.iff2 = false, false
		case ei:
			z80.iff1, z80.iff2 = true, true
		case rlca:
			a7 := z80.reg.A >> 7
			z80.reg.A = z80.reg.A<<1 | a7
			z80.reg.F = z80.reg.F&(fS|fZ|fP) | z80.reg.A&(fY|fX) | a7
		case rrca:
			a0 := z80.reg.A & 0x01
			z80.reg.A = z80.reg.A>>1 | a0<<7
			z80.reg.F = z80.reg.F&(fS|fZ|fP) | z80.reg.A&(fY|fX) | a0
		case rla:
			a7 := z80.reg.A >> 7
			z80.reg.A = z80.reg.A<<1 | z80.reg.F&fC
			z80.reg.F = z80.reg.F&(fS|fZ|fP) | z80.reg.A&(fY|fX) | a7
		case rra:
			a0 := z80.reg.A & 0x01
			z80.reg.A = z80.reg.A>>1 | z80.reg.F&fC<<7
			z80.reg.F = z80.reg.F&(fS|fZ|fP) | z80.reg.A&(fY|fX) | a0
		case cpl:
			z80.reg.A = ^z80.reg.A
			z80.reg.F = z80.reg.F&(fS|fZ|fP|fC) | fH | fN | z80.reg.A&(fY|fX)
		case scf:
			z80.reg.F = z80.reg.F&(fS|fZ|fP) | fC | z80.reg.A&(fY|fX)
		case ccf:
			z80.reg.F = (z80.reg.F&(fS|fZ|fP|fC) | z80.reg.F&fC<<4 | z80.reg.A&(fY|fX)) ^ fC
		case daa:
			cf := z80.reg.F & fC
			hf := z80.reg.F & fH
			nf := z80.reg.F & fN
			lb := z80.reg.A & 0x0F
			diff := byte(0)
			z80.reg.F &= fN
			if cf != 0 || z80.reg.A > 0x99 {
				diff = 0x60
				z80.reg.F |= fC
			}
			if hf != 0 || lb > 0x09 {
				diff += 0x06
			}
			if nf == 0 {
				z80.reg.A += diff
			} else {
				z80.reg.A -= diff
			}
			z80.reg.F |= parity[z80.reg.A] | z80.reg.A&(fS|fY|fX)
			if z80.reg.A == 0 {
				z80.reg.F |= fZ
			}
			if nf == 0 && lb > 0x09 || nf != 0 && hf != 0 && lb < 0x06 {
				z80.reg.F |= fH
			}
		case ex_af_af:
			z80.reg.A, z80.reg.A_ = z80.reg.A_, z80.reg.A
			z80.reg.F, z80.reg.F_ = z80.reg.F_, z80.reg.F
		case exx:
			z80.reg.B, z80.reg.B_, z80.reg.C, z80.reg.C_ = z80.reg.B_, z80.reg.B, z80.reg.C_, z80.reg.C
			z80.reg.D, z80.reg.D_, z80.reg.E, z80.reg.E_ = z80.reg.D_, z80.reg.D, z80.reg.E_, z80.reg.E
			z80.reg.H, z80.reg.H_, z80.reg.L, z80.reg.L_ = z80.reg.H_, z80.reg.H, z80.reg.L_, z80.reg.L
		case ex_de_hl:
			z80.reg.D, z80.reg.E, z80.reg.H, z80.reg.L = z80.reg.H, z80.reg.L, z80.reg.D, z80.reg.E
		case ex_sp_hl:
			h, l := z80.reg.r(rH), z80.reg.r(rL)
			x, y := z80.read(z80.reg.SP+1), z80.read(z80.reg.SP)
			z80.contention(z80.reg.SP+1, 1)
			z80.write(z80.reg.SP, *l)
			z80.write(z80.reg.SP+1, *h)
			z80.contention(z80.reg.SP, 2)
			*h, *l = x, y
		case add_a_n, add_a_a, add_a_b, add_a_c, add_a_d, add_a_e, add_a_h, add_a_l, add_a_hl:
			a := z80.reg.A
			var n byte
			switch opcode {
			case add_a_n:
				n = z80.readByte()
			case add_a_hl:
				hl := z80.getHL()
				if z80.reg.prefix != noPrefix {
					z80.contention(z80.reg.PC-1, 5)
				}
				n = z80.read(hl)
			default:
				n = *z80.reg.r(opcode & 0b00000111)
			}
			z80.reg.A += n
			z80.reg.F = (fS | fY | fX) & z80.reg.A
			if z80.reg.A == 0 {
				z80.reg.F |= fZ
			}
			z80.reg.F |= (a ^ n ^ z80.reg.A) & fH
			if (a^n)&0x80 == 0 && (a^z80.reg.A)&0x80 != 0 {
				z80.reg.F |= fP
			}
			if z80.reg.A < a {
				z80.reg.F |= fC
			}
		case adc_a_n, adc_a_a, adc_a_b, adc_a_c, adc_a_d, adc_a_e, adc_a_h, adc_a_l, adc_a_hl:
			var n byte
			switch opcode {
			case adc_a_n:
				n = z80.readByte()
			case adc_a_hl:
				hl := z80.getHL()
				if z80.reg.prefix != noPrefix {
					z80.contention(z80.reg.PC-1, 5)
				}
				n = z80.read(hl)
			default:
				n = *z80.reg.r(opcode & 0b00000111)
			}
			cf := z80.reg.F & fC
			sum_w := uint16(z80.reg.A) + uint16(n) + uint16(cf)
			sum_b := byte(sum_w)
			z80.reg.F = (fS | fY | fX) & sum_b
			if sum_b == 0 {
				z80.reg.F |= fZ
			}
			z80.reg.F |= (z80.reg.A ^ n ^ sum_b) & fH
			if (z80.reg.A^n)&0x80 == 0 && (z80.reg.A^sum_b)&0x80 != 0 {
				z80.reg.F |= fP
			}
			if sum_w > 0xff {
				z80.reg.F |= fC
			}
			z80.reg.A = sum_b
		case add_hl_bc, add_hl_de, add_hl_hl, add_hl_sp:
			z80.contention(z80.reg.IR(), 7)
			hl := z80.reg.HL()
			var nn uint16
			switch opcode {
			case add_hl_bc:
				nn = z80.reg.BC()
			case add_hl_de:
				nn = z80.reg.DE()
			case add_hl_hl:
				nn = hl
			case add_hl_sp:
				nn = z80.reg.SP
			}
			sum := hl + nn
			z80.reg.setHL(sum)
			z80.reg.F = z80.reg.F & ^(fH|fN|fC) | byte((hl^nn^sum)>>8)&fH | byte(sum>>8)&(fY|fX)
			if sum < hl {
				z80.reg.F |= fC
			}
		case sub_a, sub_b, sub_c, sub_d, sub_e, sub_h, sub_l, sub_hl, sub_n:
			a := z80.reg.A
			var n byte
			switch opcode {
			case sub_n:
				n = z80.readByte()
			case sub_hl:
				hl := z80.getHL()
				if z80.reg.prefix != noPrefix {
					z80.contention(z80.reg.PC-1, 1)
				}
				n = z80.read(hl)
			default:
				n = *z80.reg.r(opcode & 0b00000111)
			}
			z80.reg.A -= n
			z80.reg.F = (fS|fY|fX)&z80.reg.A | fN | (a^n^z80.reg.A)&fH
			if z80.reg.A == 0 {
				z80.reg.F |= fZ
			}
			if (a^n)&0x80 != 0 && (a^z80.reg.A)&0x80 != 0 {
				z80.reg.F |= fP
			}
			if z80.reg.A > a {
				z80.reg.F |= fC
			}
		case cp_a, cp_b, cp_c, cp_d, cp_e, cp_h, cp_l, cp_hl, cp_n:
			var n byte
			switch opcode {
			case cp_n:
				n = z80.readByte()
			case cp_hl:
				hl := z80.getHL()
				if z80.reg.prefix != noPrefix {
					z80.contention(z80.reg.PC-1, 5)
				}
				n = z80.read(hl)
			default:
				n = *z80.reg.r(opcode & 0b00000111)
			}
			test := z80.reg.A - n
			z80.reg.F = fN | fS&test | n&(fY|fX) | byte(z80.reg.A^n^test)&fH
			if test == 0 {
				z80.reg.F |= fZ
			}
			if (z80.reg.A^n)&0x80 != 0 && (z80.reg.A^test)&0x80 != 0 {
				z80.reg.F |= fP
			}
			if test > z80.reg.A {
				z80.reg.F |= fC
			}
		case sbc_a_a, sbc_a_b, sbc_a_c, sbc_a_d, sbc_a_e, sbc_a_h, sbc_a_l, sbc_a_hl, sbc_a_n:
			var n byte
			switch opcode {
			case sbc_a_n:
				n = z80.readByte()
			case sbc_a_hl:
				hl := z80.getHL()
				if z80.reg.prefix != noPrefix {
					z80.contention(z80.reg.PC-1, 5)
				}
				n = z80.read(hl)
			default:
				n = *z80.reg.r(opcode & 0b00000111)
			}
			cf := z80.reg.F & fC
			sub_w := uint16(z80.reg.A) - uint16(n) - uint16(cf)
			sub_b := byte(sub_w)
			z80.reg.F = (fS|fY|fX)&sub_b | fN | byte(z80.reg.A^n^sub_b)&fH
			if sub_b == 0 {
				z80.reg.F |= fZ
			}
			if (z80.reg.A^n)&0x80 != 0 && (sub_b^z80.reg.A)&0x80 != 0 {
				z80.reg.F |= fP
			}
			if sub_w > 0xff {
				z80.reg.F |= fC
			}
			z80.reg.A = sub_b
		case and_a, and_b, and_c, and_d, and_e, and_h, and_l, and_hl, and_n:
			var n byte
			switch opcode {
			case and_n:
				n = z80.readByte()
			case and_hl:
				hl := z80.getHL()
				if z80.reg.prefix != noPrefix {
					z80.contention(z80.reg.PC-1, 5)
				}
				n = z80.read(hl)
			default:
				n = *z80.reg.r(opcode & 0b00000111)
			}
			z80.reg.A &= n
			z80.reg.F = (fS|fY|fX)&z80.reg.A | fH | parity[z80.reg.A]
			if z80.reg.A == 0 {
				z80.reg.F |= fZ
			}
		case or_a, or_b, or_c, or_d, or_e, or_h, or_l, or_hl, or_n:
			var n byte
			switch opcode {
			case or_n:
				n = z80.readByte()
			case or_hl:
				hl := z80.getHL()
				if z80.reg.prefix != noPrefix {
					z80.contention(z80.reg.PC-1, 5)
				}
				n = z80.read(hl)
			default:
				n = *z80.reg.r(opcode & 0b00000111)
			}
			z80.reg.A |= n
			z80.reg.F = (fS|fY|fX)&z80.reg.A | parity[z80.reg.A]
			if z80.reg.A == 0 {
				z80.reg.F |= fZ
			}
		case xor_a, xor_b, xor_c, xor_d, xor_e, xor_h, xor_l, xor_hl, xor_n:
			var n byte
			switch opcode {
			case xor_n:
				n = z80.readByte()
			case xor_hl:
				hl := z80.getHL()
				if z80.reg.prefix != noPrefix {
					z80.contention(z80.reg.PC-1, 5)
				}
				n = z80.read(hl)
			default:
				n = *z80.reg.r(opcode & 0b00000111)
			}
			z80.reg.A ^= n
			z80.reg.F = (fS|fY|fX)&z80.reg.A | parity[z80.reg.A]
			if z80.reg.A == 0 {
				z80.reg.F |= fZ
			}
		case ld_a_n, ld_b_n, ld_c_n, ld_d_n, ld_e_n, ld_h_n, ld_l_n:
			r := z80.reg.r(opcode & 0b00111000 >> 3)
			*r = z80.readByte()
		case
			ld_a_a, ld_a_b, ld_a_c, ld_a_d, ld_a_e, ld_a_h, ld_a_l,
			ld_b_a, ld_b_b, ld_b_c, ld_b_d, ld_b_e, ld_b_h, ld_b_l,
			ld_c_a, ld_c_b, ld_c_c, ld_c_d, ld_c_e, ld_c_h, ld_c_l,
			ld_d_a, ld_d_b, ld_d_c, ld_d_d, ld_d_e, ld_d_h, ld_d_l,
			ld_e_a, ld_e_b, ld_e_c, ld_e_d, ld_e_e, ld_e_h, ld_e_l,
			ld_h_a, ld_h_b, ld_h_c, ld_h_d, ld_h_e, ld_h_h, ld_h_l,
			ld_l_a, ld_l_b, ld_l_c, ld_l_d, ld_l_e, ld_l_h, ld_l_l:
			rs := z80.reg.r(opcode & 0b00000111)
			rd := z80.reg.r(opcode & 0b00111000 >> 3)
			*rd = *rs
		case ld_bc_nn:
			z80.reg.C, z80.reg.B = z80.readByte(), z80.readByte()
		case ld_de_nn:
			z80.reg.E, z80.reg.D = z80.readByte(), z80.readByte()
		case ld_hl_nn:
			h, l := z80.reg.r(rH), z80.reg.r(rL)
			*l, *h = z80.readByte(), z80.readByte()
		case ld_sp_nn:
			z80.reg.SP = z80.readWord()
		case ld_sp_hl:
			z80.contention(z80.reg.IR(), 2)
			z80.reg.SP = z80.reg.HL()
		case ld_hl_mm:
			addr := z80.readWord()
			h, l := z80.reg.r(rH), z80.reg.r(rL)
			*l = z80.read(addr)
			*h = z80.read(addr + 1)
		case ld_mm_hl:
			addr := z80.readWord()
			h, l := z80.reg.r(rH), z80.reg.r(rL)
			z80.write(addr, *l)
			z80.write(addr+1, *h)
		case ld_mhl_n:
			hl := z80.getHL()
			n := z80.readByte()
			if z80.reg.prefix != noPrefix {
				z80.contention(z80.reg.PC-1, 2)
			}
			z80.write(hl, n)
		case ld_mm_a:
			z80.write(z80.readWord(), z80.reg.A)
		case ld_a_mm:
			z80.reg.A = z80.read(z80.readWord())
		case ld_bc_a:
			z80.write(z80.reg.BC(), z80.reg.A)
		case ld_de_a:
			z80.write(z80.reg.DE(), z80.reg.A)
		case ld_a_bc:
			z80.reg.A = z80.read(z80.reg.BC())
		case ld_a_de:
			z80.reg.A = z80.read(z80.reg.DE())
		case ld_a_hl, ld_b_hl, ld_c_hl, ld_d_hl, ld_e_hl, ld_h_hl, ld_l_hl:
			hl := z80.getHL()
			if z80.reg.prefix != noPrefix {
				z80.contention(z80.reg.PC-1, 5)
			}
			*z80.reg.raw[opcode&0b00111000>>3] = z80.read(hl)
		case ld_hl_a, ld_hl_b, ld_hl_c, ld_hl_d, ld_hl_e, ld_hl_h, ld_hl_l:
			hl := z80.getHL()
			if z80.reg.prefix != noPrefix {
				z80.contention(z80.reg.PC-1, 5)
			}
			z80.write(hl, *z80.reg.raw[opcode&0b00000111])
		case inc_a, inc_b, inc_c, inc_d, inc_e, inc_h, inc_l:
			r := z80.reg.r(opcode & 0b00111000 >> 3)
			z80.reg.F &= fC
			if *r == 0x7F {
				z80.reg.F |= fP
			}
			if *r&0x0F == 0x0F {
				z80.reg.F |= fH
			}
			*r += 1
			z80.reg.F |= *r & (fS | fY | fX)
			if *r == 0 {
				z80.reg.F |= fZ
			}
		case inc_bc:
			z80.contention(z80.reg.IR(), 2)
			z80.reg.setBC(z80.reg.BC() + 1)
		case inc_de:
			z80.contention(z80.reg.IR(), 2)
			z80.reg.setDE(z80.reg.DE() + 1)
		case inc_hl:
			z80.contention(z80.reg.IR(), 2)
			z80.reg.setHL(z80.reg.HL() + 1)
		case inc_sp:
			z80.contention(z80.reg.IR(), 2)
			z80.reg.SP += 1
		case inc_mhl:
			addr := z80.getHL()
			if z80.reg.prefix != noPrefix {
				z80.contention(z80.reg.PC-1, 5)
			}
			b := z80.read(addr)
			z80.contention(addr, 1)
			z80.reg.F &= fC
			if b == 0x7F {
				z80.reg.F |= fP
			}
			if b&0x0F == 0x0F {
				z80.reg.F |= fH
			}
			b += 1
			if b == 0x00 {
				z80.reg.F |= fZ
			}
			z80.reg.F |= b & (fS | fY | fX)
			z80.write(addr, b)
		case dec_a, dec_b, dec_c, dec_d, dec_e, dec_h, dec_l:
			r := z80.reg.r(opcode & 0b00111000 >> 3)
			z80.reg.F = z80.reg.F&fC | fN
			if *r == 0x80 {
				z80.reg.F |= fP
			}
			if *r&0x0F == 0 {
				z80.reg.F |= fH
			}
			*r -= 1
			z80.reg.F |= *r & (fS | fY | fX)
			if *r == 0 {
				z80.reg.F |= fZ
			}
		case dec_bc:
			z80.contention(z80.reg.IR(), 2)
			z80.reg.setBC(z80.reg.BC() - 1)
		case dec_de:
			z80.contention(z80.reg.IR(), 2)
			z80.reg.setDE(z80.reg.DE() - 1)
		case dec_hl:
			z80.contention(z80.reg.IR(), 2)
			z80.reg.setHL(z80.reg.HL() - 1)
		case dec_sp:
			z80.contention(z80.reg.IR(), 2)
			z80.reg.SP -= 1
		case dec_mhl:
			addr := z80.getHL()
			if z80.reg.prefix != noPrefix {
				z80.contention(z80.reg.PC-1, 5)
			}
			b := z80.read(addr)
			z80.contention(addr, 1)
			z80.reg.F = z80.reg.F&fC | fN
			if b == 0x80 {
				z80.reg.F |= fP
			}
			if b&0x0F == 0 {
				z80.reg.F |= fH
			}
			b -= 1
			if b == 0x00 {
				z80.reg.F |= fZ
			}
			z80.reg.F |= b & (fS | fY | fX)
			z80.write(addr, b)
		case jr_o:
			o := z80.readByte()
			z80.contention(z80.reg.PC-1, 5)
			if o&0x80 == 0 {
				z80.reg.PC += uint16(o)
			} else {
				z80.reg.PC -= uint16(^o + 1)
			}
		case jr_z_o:
			o := z80.readByte()
			if z80.reg.F&fZ == fZ {
				z80.contention(z80.reg.PC-1, 5)
				if o&0x80 == 0 {
					z80.reg.PC += uint16(o)
				} else {
					z80.reg.PC -= uint16(^o + 1)
				}
			}
		case jr_nz_o:
			o := z80.readByte()
			if z80.reg.F&fZ == 0 {
				z80.contention(z80.reg.PC-1, 5)
				if o&0x80 == 0 {
					z80.reg.PC += uint16(o)
				} else {
					z80.reg.PC -= uint16(^o + 1)
				}
			}
		case jr_c:
			o := z80.readByte()
			if z80.reg.F&fC == fC {
				z80.contention(z80.reg.PC-1, 5)
				if o&0x80 == 0 {
					z80.reg.PC += uint16(o)
				} else {
					z80.reg.PC -= uint16(^o + 1)
				}
			}
		case jr_nc_o:
			o := z80.readByte()
			if z80.reg.F&fC == 0 {
				z80.contention(z80.reg.PC-1, 5)
				if o&0x80 == 0 {
					z80.reg.PC += uint16(o)
				} else {
					z80.reg.PC -= uint16(^o + 1)
				}
			}
		case djnz:
			z80.contention(z80.reg.IR(), 1)
			o := z80.readByte()
			z80.reg.B -= 1
			if z80.reg.B != 0 {
				z80.contention(z80.reg.PC, 5)
				if o&0x80 == 0 {
					z80.reg.PC += uint16(o)
				} else {
					z80.reg.PC -= uint16(^o + 1)
				}
			}
		case jp_nn:
			z80.reg.PC = z80.readWord()
		case jp_c_nn, jp_m_nn, jp_nc_nn, jp_nz_nn, jp_p_nn, jp_pe_nn, jp_po_nn, jp_z_nn:
			pc := z80.readWord()
			if z80.shouldJump(opcode) {
				z80.reg.PC = pc
			}
		case jp_hl:
			z80.reg.PC = z80.reg.HL()
		case call_nn:
			pc := z80.readWord()
			z80.contention(z80.reg.PC, 1)
			z80.reg.SP -= 1
			z80.write(z80.reg.SP, byte(z80.reg.PC>>8))
			z80.reg.SP -= 1
			z80.write(z80.reg.SP, byte(z80.reg.PC))
			z80.reg.PC = pc
		case call_c_nn, call_m_nn, call_nc_nn, call_nz_nn, call_p_nn, call_pe_nn, call_po_nn, call_z_nn:
			pc := z80.readWord()
			if z80.shouldJump(opcode) {
				z80.contention(z80.reg.PC, 1)
				z80.reg.SP -= 1
				z80.write(z80.reg.SP, byte(z80.reg.PC>>8))
				z80.reg.SP -= 1
				z80.write(z80.reg.SP, byte(z80.reg.PC))
				z80.reg.PC = pc
			}
		case ret:
			z80.reg.PC = uint16(z80.read(z80.reg.SP+1))<<8 | uint16(z80.read(z80.reg.SP))
			z80.reg.SP += 2
		case ret_c, ret_m, ret_nc, ret_nz, ret_p, ret_pe, ret_po, ret_z:
			z80.contention(z80.reg.IR(), 1)
			if z80.shouldJump(opcode) {
				z80.reg.PC = uint16(z80.read(z80.reg.SP+1))<<8 | uint16(z80.read(z80.reg.SP))
				z80.reg.SP += 2
			}
		case rst_00h, rst_08h, rst_10h, rst_18h, rst_20h, rst_28h, rst_30h, rst_38h:
			z80.contention(z80.reg.IR(), 1)
			z80.reg.SP -= 1
			z80.write(z80.reg.SP, byte(z80.reg.PC>>8))
			z80.reg.SP -= 1
			z80.write(z80.reg.SP, byte(z80.reg.PC))
			z80.reg.PC = uint16(8 * ((opcode & 0b00111000) >> 3))
		case push_af:
			z80.contention(z80.reg.IR(), 1)
			z80.reg.SP -= 1
			z80.write(z80.reg.SP, z80.reg.A)
			z80.reg.SP -= 1
			z80.write(z80.reg.SP, z80.reg.F)
		case push_bc:
			z80.contention(z80.reg.IR(), 1)
			z80.reg.SP -= 1
			z80.write(z80.reg.SP, z80.reg.B)
			z80.reg.SP -= 1
			z80.write(z80.reg.SP, z80.reg.C)
		case push_de:
			z80.contention(z80.reg.IR(), 1)
			z80.reg.SP -= 1
			z80.write(z80.reg.SP, z80.reg.D)
			z80.reg.SP -= 1
			z80.write(z80.reg.SP, z80.reg.E)
		case push_hl:
			z80.contention(z80.reg.IR(), 1)
			z80.reg.SP -= 1
			z80.write(z80.reg.SP, *z80.reg.r(rH))
			z80.reg.SP -= 1
			z80.write(z80.reg.SP, *z80.reg.r(rL))
		case pop_af:
			z80.reg.A, z80.reg.F = z80.read(z80.reg.SP+1), z80.read(z80.reg.SP)
			z80.reg.SP += 2
		case pop_bc:
			z80.reg.B, z80.reg.C = z80.read(z80.reg.SP+1), z80.read(z80.reg.SP)
			z80.reg.SP += 2
		case pop_de:
			z80.reg.D, z80.reg.E = z80.read(z80.reg.SP+1), z80.read(z80.reg.SP)
			z80.reg.SP += 2
		case pop_hl:
			*z80.reg.r(rH), *z80.reg.r(rL) = z80.read(z80.reg.SP+1), z80.read(z80.reg.SP)
			z80.reg.SP += 2
		case in_a_n:
			z80.reg.A = z80.readBus(z80.reg.A, z80.readByte())
		case out_n_a:
			z80.writeBus(z80.reg.A, z80.readByte(), z80.reg.A)
		case prefix_cb:
			z80.incR()
			z80.prefixCB()
		case prefix_ed:
			z80.incR()
			z80.prefixED(z80.fetch())
		case useIX:
			z80.incR()
			z80.reg.prefix = useIX
			continue
		case useIY:
			z80.incR()
			z80.reg.prefix = useIY
			continue
		}
		z80.reg.prefix = noPrefix
	}
}

func (z80 *Z80) shouldJump(opcode byte) bool {
	switch opcode & 0b00111000 {
	case 0b00000000: // Non-Zero (NZ)
		return z80.reg.F&fZ == 0
	case 0b00001000: // Zero (Z)
		return z80.reg.F&fZ != 0
	case 0b00010000: // Non Carry (NC)
		return z80.reg.F&fC == 0
	case 0b00011000: // Carry (C)
		return z80.reg.F&fC != 0
	case 0b00100000: // Parity Odd (PO)
		return z80.reg.F&fP == 0
	case 0b00101000: // Parity Even (PE)
		return z80.reg.F&fP != 0
	case 0b00110000: // Sign Positive (P)
		return z80.reg.F&fS == 0
	case 0b00111000: // Sign Negative (M)
		return z80.reg.F&fS != 0
	}

	panic(fmt.Sprintf("Invalid opcode %v", opcode))
}

// Returns value of HL / (IX + d) / (IY + d) register. The current prefix
// determines whether to use IX or IY register instead of HL.
func (z80 *Z80) getHL() uint16 {
	if z80.reg.prefix != noPrefix {
		return z80.getHLOffset(z80.readByte())
	}
	return z80.reg.HL()
}

// For IX or IY add offset to register value, otherwise return HL.
func (z80 *Z80) getHLOffset(offset byte) uint16 {
	hl := z80.reg.HL()
	if offset == 0 {
		return hl
	} else if offset&0x80 == 0 {
		return hl + uint16(offset)
	} else {
		return hl - uint16(^offset+1)
	}
}

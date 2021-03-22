package z80

import "fmt"

const (
	rB = 0b000
	rC = 0b001
	rD = 0b010
	rE = 0b011
	rH = 0b100
	rL = 0b101
	rA = 0b111

	rBC = 0b00
	rDE = 0b01
	rHL = 0b10
	rSP = 0b11
)

// The Flag registers, F and F', supply information to the user about the status of the Z80
// CPU at any particular time.
// 		 7   6   5   4   3   2   1   0
// 		 S   Z   Y   H   X   P   N   C
// S = sign, Z = zero, H = half carry, P = parity/overflow, N = add/substract, C = carry
//
// X and Y flags are undocumented.
const (
	fNONE byte = 0x00
	fC    byte = 0x01
	fN    byte = 0x02
	fP    byte = 0x04
	fX    byte = 0x08
	fH    byte = 0x10
	fY    byte = 0x20
	fZ    byte = 0x40
	fS    byte = 0x80
	fALL  byte = 0xFF
)

type registers struct {
	A, B, C, D, E, H, L, F         byte      // Standard registers
	A_, B_, C_, D_, E_, H_, L_, F_ byte      // Shadow registers
	IXH, IXL, IYH, IYL             byte      // Index registers, supports unofficial
	SP                             uint16    // Stack Pointer
	PC                             uint16    // Program Counter
	I, R                           byte      // Interrupt Vector / Memory Refresh
	raw                            []*byte   // raw index of 8-bit registers
	prefixed                       [][]*byte // index of registered for IX, IY or no prefix
	prefix                         byte      // indicates IX or IY prefix, or no prefix
}

// Initialises new instance of the registers struct.
func newRegisters() *registers {
	r := &registers{}
	r.prefixed = [][]*byte{
		noPrefix: {
			rA: &r.A, rB: &r.B, rC: &r.C, rD: &r.D,
			rE: &r.E, rH: &r.H, rL: &r.L,
		},
		useIX: {
			rA: &r.A, rB: &r.B, rC: &r.C, rD: &r.D,
			rE: &r.E, rH: &r.IXH, rL: &r.IXL,
		},
		useIY: {
			rA: &r.A, rB: &r.B, rC: &r.C, rD: &r.D,
			rE: &r.E, rH: &r.IYH, rL: &r.IYL,
		},
	}
	r.raw = r.prefixed[noPrefix]

	return r
}

// Gets the specified 8-bit register, respecting operation may be prefixed.
func (r *registers) r(reg byte) *byte {
	return r.prefixed[r.prefix][reg]
}

// Sets the value of the specified 8-bit register, respecting operation may be prefixed.
func (r *registers) setR(reg, value byte) {
	*r.prefixed[r.prefix][reg] = value
}

// Gets the value of the specified 16-bit register, respecting operation may be prefixed.
func (r *registers) rr(reg byte) uint16 {
	switch reg {
	case rBC:
		return r.BC()
	case rDE:
		return r.DE()
	case rHL:
		return r.HL()
	case rSP:
		return r.SP
	}

	panic(fmt.Sprintf("Invalid 16 bit register %v", reg))
}

// Sets the value of the specified 16-bit registers, respecting operation may be prefixed.
func (r *registers) setRR(reg byte, value uint16) {
	switch reg {
	case rBC:
		r.setBC(value)
	case rDE:
		r.setDE(value)
	case rHL:
		r.setHL(value)
	case rSP:
		r.SP = value
	}
}

// Gets the BC register value.
func (r *registers) BC() uint16 {
	return uint16(r.B)<<8 | uint16(r.C)
}

// Sets the BC register value.
func (r *registers) setBC(nn uint16) {
	r.B, r.C = byte(nn>>8), byte(nn)
}

// Gets the DE register value.
func (r *registers) DE() uint16 {
	return uint16(r.D)<<8 | uint16(r.E)
}

// Sets the DE register value.
func (r *registers) setDE(nn uint16) {
	r.D, r.E = byte(nn>>8), byte(nn)
}

// Gets the HL register value, respecting operation may be prefixed.
func (r *registers) HL() uint16 {
	switch r.prefix {
	case useIX:
		return uint16(r.IXH)<<8 | uint16(r.IXL)
	case useIY:
		return uint16(r.IYH)<<8 | uint16(r.IYL)
	default:
		return uint16(r.H)<<8 | uint16(r.L)
	}
}

// Sets the HL register value, respecting operation may be prefixed.
func (r *registers) setHL(value uint16) {
	h, l := byte(value>>8), byte(value)
	switch r.prefix {
	case useIX:
		r.IXH, r.IXL = h, l
	case useIY:
		r.IYH, r.IYL = h, l
	default:
		r.H, r.L = h, l
	}
}

// Gets the virtual IR register value.
func (r *registers) IR() uint16 {
	return uint16(r.I)<<8 | uint16(r.R)
}

// Increments R register.
func (r *registers) IncR() {
	r.R = r.R&0x80 | (r.R+1)&0x7F
}

package z80

import "fmt"

const (
	r_B = 0b000
	r_C = 0b001
	r_D = 0b010
	r_E = 0b011
	r_H = 0b100
	r_L = 0b101
	r_A = 0b111

	r_BC = 0b00
	r_DE = 0b01
	r_HL = 0b10
	r_SP = 0b11
)

// The Flag registers, F and F', supply information to the user about the status of the Z80
// CPU at any particular time.
// 		 7   6   5   4   3   2   1   0
// 		 S   Z   Y   H   X   P   N   C
// S = sign, Z = zero, H = half carry, P = parity/overflow, N = add/substract, C = carry
//
// X and Y flags are undocumented.
const (
	f_NONE byte = 0x00
	f_C    byte = 0x01
	f_N    byte = 0x02
	f_P    byte = 0x04
	f_X    byte = 0x08
	f_H    byte = 0x10
	f_Y    byte = 0x20
	f_Z    byte = 0x40
	f_S    byte = 0x80
	f_ALL  byte = 0xFF
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
			r_A: &r.A, r_B: &r.B, r_C: &r.C, r_D: &r.D,
			r_E: &r.E, r_H: &r.H, r_L: &r.L,
		},
		useIX: {
			r_A: &r.A, r_B: &r.B, r_C: &r.C, r_D: &r.D,
			r_E: &r.E, r_H: &r.IXH, r_L: &r.IXL,
		},
		useIY: {
			r_A: &r.A, r_B: &r.B, r_C: &r.C, r_D: &r.D,
			r_E: &r.E, r_H: &r.IYH, r_L: &r.IYL,
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
	case r_BC:
		return r.BC()
	case r_DE:
		return r.DE()
	case r_HL:
		return r.HL()
	case r_SP:
		return r.SP
	}

	panic(fmt.Sprintf("Invalid 16 bit register %v", reg))
}

// Sets the value of the specified 16-bit registers, respecting operation may be prefixed.
func (r *registers) setRR(reg byte, value uint16) {
	switch reg {
	case r_BC:
		r.setBC(value)
	case r_DE:
		r.setDE(value)
	case r_HL:
		r.setHL(value)
	case r_SP:
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

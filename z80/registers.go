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
// X and Y are undocumented and currently not used in the implementation.
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
	f_ALL  byte = f_S | f_Z | f_H | f_P | f_N | f_C
)

type registers struct {
	// Standard registers
	A, B, C, D, E, H, L, F byte
	// Shadow registers
	A_, B_, C_, D_, E_, H_, L_, F_ byte
	// Unofficial registers
	IXH, IXL, IYH, IYL byte
	// Other & special registers
	SP   uint16
	PC   uint16
	I, R byte
	// Helper register index
	get      []*byte
	prefixed [][]*byte
	prefix   byte
}

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
	r.get = r.prefixed[noPrefix]

	return r
}

func (r *registers) getReg(reg byte) *byte {
	return r.prefixed[r.prefix][reg]
}

func (r *registers) setReg(reg, value byte) {
	*r.prefixed[r.prefix][reg] = value
}

// Gets the value of one of the specified 16-bit registers
func (r *registers) getReg16(reg byte) uint16 {
	switch reg {
	case r_BC:
		return r.getBC()
	case r_DE:
		return r.getDE()
	case r_HL:
		return r.getHL()
	case r_SP:
		return r.SP
	}

	panic(fmt.Sprintf("Invalid 16 bit register %v", reg))
}

// Sets the value of one of the specified 16-bit registers
func (r *registers) setReg16(reg byte, value uint16) {
	switch reg {
	case r_BC:
		r.setBC(value)
	case r_DE:
		r.setDE(value)
	case r_HL:
		r.setHLw(value)
	case r_SP:
		r.SP = value
	}
}

func (r *registers) getBC() uint16 {
	return uint16(r.B)<<8 | uint16(r.C)
}

func (r *registers) setBC(nn uint16) {
	r.B, r.C = byte(nn>>8), byte(nn)
}

func (r *registers) getDE() uint16 {
	return uint16(r.D)<<8 | uint16(r.E)
}

func (r *registers) setDE(nn uint16) {
	r.D, r.E = byte(nn>>8), byte(nn)
}

func (r *registers) getHL() uint16 {
	switch r.prefix {
	case useIX:
		return uint16(r.IXH)<<8 | uint16(r.IXL)
	case useIY:
		return uint16(r.IYH)<<8 | uint16(r.IYL)
	default:
		return uint16(r.H)<<8 | uint16(r.L)
	}
}

// Sets HL register using 16-bit value
func (r *registers) setHLw(value uint16) {
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

// Sets HL register using two 8-bit values
func (r *registers) setHLb(h, l byte) {
	switch r.prefix {
	case useIX:
		r.IXH, r.IXL = h, l
	case useIY:
		r.IYH, r.IYL = h, l
	default:
		r.H, r.L = h, l
	}
}

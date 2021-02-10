package z80

import "fmt"

const (
	r_B  = 0b000
	r_C  = 0b001
	r_D  = 0b010
	r_E  = 0b011
	r_H  = 0b100
	r_L  = 0b101
	r_HL = 0b110
	r_A  = 0b111
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
	A byte
	B byte
	C byte
	D byte
	E byte
	H byte
	L byte
	F byte
	// Shadow registers
	A_ byte
	B_ byte
	C_ byte
	D_ byte
	E_ byte
	H_ byte
	L_ byte
	F_ byte
	// Other & special registers
	IX [2]byte
	IY [2]byte
	SP word
	I  byte
	R  byte
	// Helper register index
	regs8 map[byte]*byte
}

func newRegisters() *registers {
	r := &registers{}
	r.regs8 = map[byte]*byte{
		r_A: &r.A,
		r_B: &r.B,
		r_C: &r.C,
		r_D: &r.D,
		r_E: &r.E,
		r_H: &r.H,
		r_L: &r.L,
	}

	return r
}

func (r *registers) getR(reg byte) *byte {
	return r.regs8[reg]
}

func (r *registers) getReg(reg, prefix byte) *byte {
	switch reg {
	case r_A:
		return &r.A
	case r_B:
		return &r.B
	case r_C:
		return &r.C
	case r_D:
		return &r.D
	case r_E:
		return &r.E
	case r_H:
		switch prefix {
		case use_ix:
			return &r.IX[0]
		case use_iy:
			return &r.IY[0]
		default:
			return &r.H
		}
	case r_L:
		switch prefix {
		case use_ix:
			return &r.IX[1]
		case use_iy:
			return &r.IY[1]
		default:
			return &r.L
		}
	}

	panic(fmt.Sprintf("getReg: Invalid register %v", reg))
}

func (r *registers) setReg(reg, prefix, value byte) {
	switch reg {
	case r_A:
		r.A = value
	case r_B:
		r.B = value
	case r_C:
		r.C = value
	case r_D:
		r.D = value
	case r_E:
		r.E = value
	case r_H:
		switch prefix {
		case use_ix:
			r.IX[0] = value
		case use_iy:
			r.IY[0] = value
		default:
			r.H = value
		}
	case r_L:
		switch prefix {
		case use_ix:
			r.IX[1] = value
		case use_iy:
			r.IY[1] = value
		default:
			r.L = value
		}
	}
}

func (r *registers) getBC() word {
	return word(r.B)<<8 | word(r.C)
}

func (r *registers) setBC(nn word) {
	r.B, r.C = byte(nn>>8), byte(nn)
}

func (r *registers) getDE() word {
	return word(r.D)<<8 | word(r.E)
}

func (r *registers) setDE(nn word) {
	r.D, r.E = byte(nn>>8), byte(nn)
}

func (r *registers) getHL(prefix byte) word {
	switch prefix {
	case use_ix:
		return word(r.IX[0])<<8 | word(r.IX[1])
	case use_iy:
		return word(r.IY[0])<<8 | word(r.IY[1])
	default:
		return word(r.H)<<8 | word(r.L)
	}
}

func (r *registers) setHLw(value word, prefix byte) {
	h, l := byte(value>>8), byte(value)
	switch prefix {
	case use_ix:
		r.IX[0], r.IX[1] = h, l
	case use_iy:
		r.IY[0], r.IY[1] = h, l
	default:
		r.H, r.L = h, l
	}
}

func (r *registers) setHLb(h, l, prefix byte) {
	switch prefix {
	case use_ix:
		r.IX[0], r.IX[1] = h, l
	case use_iy:
		r.IY[0], r.IY[1] = h, l
	default:
		r.H, r.L = h, l
	}
}

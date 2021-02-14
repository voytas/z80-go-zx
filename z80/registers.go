package z80

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
	A, B, C, D, E, H, L, F byte
	// Shadow registers
	A_, B_, C_, D_, E_, H_, L_, F_ byte
	// Other & special registers
	IX, IY [2]byte
	SP     word
	I, R   byte
	// Helper register index
	get      []*byte
	prefixed map[byte][]*byte
	prefix   byte
}

func newRegisters() *registers {
	r := &registers{}
	r.prefixed = map[byte][]*byte{
		useHL: {
			r_A: &r.A, r_B: &r.B, r_C: &r.C, r_D: &r.D,
			r_E: &r.E, r_H: &r.H, r_L: &r.L,
		},
		useIX: {
			r_A: &r.A, r_B: &r.B, r_C: &r.C, r_D: &r.D,
			r_E: &r.E, r_H: &r.IX[0], r_L: &r.IX[1],
		},
		useIY: {
			r_A: &r.A, r_B: &r.B, r_C: &r.C, r_D: &r.D,
			r_E: &r.E, r_H: &r.IY[0], r_L: &r.IY[1],
		},
	}
	r.get = r.prefixed[useHL]

	return r
}

func (r *registers) getReg(reg byte) *byte {
	return r.prefixed[r.prefix][reg]
}

func (r *registers) setReg(reg, value byte) {
	*r.prefixed[r.prefix][reg] = value
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

func (r *registers) getHL() word {
	switch r.prefix {
	case useIX:
		return word(r.IX[0])<<8 | word(r.IX[1])
	case useIY:
		return word(r.IY[0])<<8 | word(r.IY[1])
	default:
		return word(r.H)<<8 | word(r.L)
	}
}

func (r *registers) setHLw(value word) {
	h, l := byte(value>>8), byte(value)
	switch r.prefix {
	case useIX:
		r.IX[0], r.IX[1] = h, l
	case useIY:
		r.IY[0], r.IY[1] = h, l
	default:
		r.H, r.L = h, l
	}
}

func (r *registers) setHLb(h, l byte) {
	switch r.prefix {
	case useIX:
		r.IX[0], r.IX[1] = h, l
	case useIY:
		r.IY[0], r.IY[1] = h, l
	default:
		r.H, r.L = h, l
	}
}

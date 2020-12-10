package z80

const (
	r_BC = 0b00
	r_DE = 0b01
	r_HL = 0b10
	r_SP = 0b11
)

// The Flag registers, F and F', supply information to the user about the status of the Z80
// CPU at any particular time.
// 		 7   6   5   4   3   2   1   0
// 		 S   Z   X5  H   X3  PV  N   C
// S = sign, Z = zero, H = half carry, PV = parity/overflow, N = add/substract, C = carry
//
// X5 and X3 are undocumented and currently not used in the implementation.
const (
	f_NONE byte = 0x00
	f_C    byte = 0x01
	f_N    byte = 0x02
	f_PV   byte = 0x04
	// f_X3   byte = 0x08 // Undocumented, not implemented yet
	f_H byte = 0x10
	// f_X5   byte = 0x20 // Undocumented, not implemented yet
	f_Z   byte = 0x40
	f_S   byte = 0x80
	f_ALL      = f_S | f_Z | f_H | f_PV | f_N | f_C
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
	IX word
	IY word
	SP word
	I  byte
	R  byte
	// Helper register index
	regs8 map[byte]*byte
}

func newRegisters() *registers {
	r := &registers{}
	r.regs8 = map[byte]*byte{
		0b000: &r.B,
		0b001: &r.C,
		0b010: &r.D,
		0b011: &r.E,
		0b100: &r.H,
		0b101: &r.L,
		0b111: &r.A,
	}

	return r
}

func (r *registers) setRRn(reg byte, nn word) {
	switch reg {
	case 0b00:
		r.B, r.C = byte(nn>>8), byte(nn)
	case 0b01:
		r.D, r.E = byte(nn>>8), byte(nn)
	case 0b10:
		r.H, r.L = byte(nn>>8), byte(nn)
	case 0b11:
		r.SP = nn
	}
}

func (r *registers) setRRnn(reg byte, lo byte, hi byte) {
	switch reg {
	case 0b00:
		r.B, r.C = hi, lo
	case 0b01:
		r.D, r.E = hi, lo
	case 0b10:
		r.H, r.L = hi, lo
	case 0b11:
		r.SP = word(hi)<<8 | word(lo)
	}
}

func (r *registers) getR(reg byte) *byte {
	return r.regs8[reg]
}

func (r *registers) getRR(reg byte) word {
	switch reg {
	case r_BC:
		return word(r.B)<<8 | word(r.C)
	case r_DE:
		return word(r.D)<<8 | word(r.E)
	case r_HL:
		return word(r.H)<<8 | word(r.L)
	case r_SP:
		return r.SP
	}

	panic("invalid 16 bit register code")
}

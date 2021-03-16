package z80

// CPU state that can be loaded or saved
type CPUState struct {
	AF, BC, DE, HL     uint16 // standard registers
	AF_, BC_, DE_, HL_ uint16 // shadow registers
	IX, IY             uint16
	PC, SP             uint16
	I, R, IM           byte
	IFF1, IFF2         bool
}

func (z80 *Z80) State(state *CPUState) {
	z80.reg.A = byte(state.AF >> 8)
	z80.reg.F = byte(state.AF)
	z80.reg.B = byte(state.BC >> 8)
	z80.reg.C = byte(state.BC)
	z80.reg.D = byte(state.DE >> 8)
	z80.reg.E = byte(state.DE)
	z80.reg.H = byte(state.HL >> 8)
	z80.reg.L = byte(state.HL)
	z80.reg.A_ = byte(state.AF_ >> 8)
	z80.reg.F_ = byte(state.AF_)
	z80.reg.B_ = byte(state.BC_ >> 8)
	z80.reg.C_ = byte(state.BC_)
	z80.reg.D_ = byte(state.DE_ >> 8)
	z80.reg.E_ = byte(state.DE_)
	z80.reg.H_ = byte(state.HL_ >> 8)
	z80.reg.L_ = byte(state.HL_)
	z80.reg.IXH = byte(state.IX >> 8)
	z80.reg.IXL = byte(state.IX)
	z80.reg.IYH = byte(state.IY >> 8)
	z80.reg.IYL = byte(state.IY)
	z80.reg.PC = state.PC
	z80.reg.SP = state.SP
	z80.reg.I = state.I
	z80.reg.R = state.R
	z80.im = state.IM
	z80.iff1 = state.IFF1
	z80.iff2 = state.IFF2
}

func (z80 *Z80) GetState() CPUState {
	return CPUState{
		AF:   uint16(z80.reg.A)<<8 | uint16(z80.reg.F),
		BC:   uint16(z80.reg.B)<<8 | uint16(z80.reg.C),
		DE:   uint16(z80.reg.D)<<8 | uint16(z80.reg.E),
		HL:   uint16(z80.reg.H)<<8 | uint16(z80.reg.L),
		AF_:  uint16(z80.reg.A_)<<8 | uint16(z80.reg.F_),
		BC_:  uint16(z80.reg.B_)<<8 | uint16(z80.reg.C_),
		DE_:  uint16(z80.reg.D_)<<8 | uint16(z80.reg.E_),
		HL_:  uint16(z80.reg.H_)<<8 | uint16(z80.reg.L_),
		IX:   uint16(z80.reg.IXH)<<8 | uint16(z80.reg.IXL),
		IY:   uint16(z80.reg.IYH)<<8 | uint16(z80.reg.IYL),
		PC:   z80.reg.PC,
		SP:   z80.reg.SP,
		I:    z80.reg.I,
		R:    z80.reg.R,
		IM:   z80.im,
		IFF1: z80.iff1,
		IFF2: z80.iff2,
	}
}

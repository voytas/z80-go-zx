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
	z80.Reg.A = byte(state.AF >> 8)
	z80.Reg.F = byte(state.AF)
	z80.Reg.B = byte(state.BC >> 8)
	z80.Reg.C = byte(state.BC)
	z80.Reg.D = byte(state.DE >> 8)
	z80.Reg.E = byte(state.DE)
	z80.Reg.H = byte(state.HL >> 8)
	z80.Reg.L = byte(state.HL)
	z80.Reg.A_ = byte(state.AF_ >> 8)
	z80.Reg.F_ = byte(state.AF_)
	z80.Reg.B_ = byte(state.BC_ >> 8)
	z80.Reg.C_ = byte(state.BC_)
	z80.Reg.D_ = byte(state.DE_ >> 8)
	z80.Reg.E_ = byte(state.DE_)
	z80.Reg.H_ = byte(state.HL_ >> 8)
	z80.Reg.L_ = byte(state.HL_)
	z80.Reg.IXH = byte(state.IX >> 8)
	z80.Reg.IXL = byte(state.IX)
	z80.Reg.IYH = byte(state.IY >> 8)
	z80.Reg.IYL = byte(state.IY)
	z80.Reg.PC = state.PC
	z80.Reg.SP = state.SP
	z80.Reg.I = state.I
	z80.Reg.R = state.R
	z80.im = state.IM
	z80.iff1 = state.IFF1
	z80.iff2 = state.IFF2
}

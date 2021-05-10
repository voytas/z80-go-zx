package z80

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voytas/z80-go-zx/z80/memory"
)

func Test_State(t *testing.T) {
	mem := &memory.BasicMemory{}
	z80 := NewZ80(mem)
	state := &CPUState{
		AF:   0x1234,
		BC:   0x2345,
		DE:   0x3456,
		HL:   0x4567,
		AF_:  0x5678,
		BC_:  0x6789,
		DE_:  0x789A,
		HL_:  0x89AB,
		IX:   0x9ABC,
		IY:   0xABCD,
		PC:   0xBCDE,
		SP:   0xCDEF,
		I:    0x1A,
		R:    0x2B,
		IM:   2,
		IFF1: true,
		IFF2: true,
	}

	z80.State(state)

	assert.Equal(t, byte(state.AF>>8), z80.Reg.A)
	assert.Equal(t, byte(state.AF), z80.Reg.F)
	assert.Equal(t, byte(state.BC>>8), z80.Reg.B)
	assert.Equal(t, byte(state.BC), z80.Reg.C)
	assert.Equal(t, byte(state.DE>>8), z80.Reg.D)
	assert.Equal(t, byte(state.DE), z80.Reg.E)
	assert.Equal(t, byte(state.HL>>8), z80.Reg.H)
	assert.Equal(t, byte(state.HL), z80.Reg.L)
	assert.Equal(t, byte(state.AF_>>8), z80.Reg.A_)
	assert.Equal(t, byte(state.AF_), z80.Reg.F_)
	assert.Equal(t, byte(state.BC_>>8), z80.Reg.B_)
	assert.Equal(t, byte(state.BC_), z80.Reg.C_)
	assert.Equal(t, byte(state.DE_>>8), z80.Reg.D_)
	assert.Equal(t, byte(state.DE_), z80.Reg.E_)
	assert.Equal(t, byte(state.HL_>>8), z80.Reg.H_)
	assert.Equal(t, byte(state.HL_), z80.Reg.L_)
	assert.Equal(t, byte(state.IX>>8), z80.Reg.IXH)
	assert.Equal(t, byte(state.IX), z80.Reg.IXL)
	assert.Equal(t, byte(state.IY>>8), z80.Reg.IYH)
	assert.Equal(t, byte(state.IY), z80.Reg.IYL)
	assert.Equal(t, state.PC, z80.Reg.PC)
	assert.Equal(t, state.SP, z80.Reg.SP)
	assert.Equal(t, state.I, z80.Reg.I)
	assert.Equal(t, state.R, z80.Reg.R)
	assert.Equal(t, byte(2), z80.im)
	assert.Equal(t, true, z80.iff1)
	assert.Equal(t, true, z80.iff2)
}

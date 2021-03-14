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

	assert.Equal(t, byte(state.AF>>8), z80.reg.A)
	assert.Equal(t, byte(state.AF), z80.reg.F)
	assert.Equal(t, byte(state.BC>>8), z80.reg.B)
	assert.Equal(t, byte(state.BC), z80.reg.C)
	assert.Equal(t, byte(state.DE>>8), z80.reg.D)
	assert.Equal(t, byte(state.DE), z80.reg.E)
	assert.Equal(t, byte(state.HL>>8), z80.reg.H)
	assert.Equal(t, byte(state.HL), z80.reg.L)
	assert.Equal(t, byte(state.AF_>>8), z80.reg.A_)
	assert.Equal(t, byte(state.AF_), z80.reg.F_)
	assert.Equal(t, byte(state.BC_>>8), z80.reg.B_)
	assert.Equal(t, byte(state.BC_), z80.reg.C_)
	assert.Equal(t, byte(state.DE_>>8), z80.reg.D_)
	assert.Equal(t, byte(state.DE_), z80.reg.E_)
	assert.Equal(t, byte(state.HL_>>8), z80.reg.H_)
	assert.Equal(t, byte(state.HL_), z80.reg.L_)
	assert.Equal(t, byte(state.IX>>8), z80.reg.IXH)
	assert.Equal(t, byte(state.IX), z80.reg.IXL)
	assert.Equal(t, byte(state.IY>>8), z80.reg.IYH)
	assert.Equal(t, byte(state.IY), z80.reg.IYL)
	assert.Equal(t, state.PC, z80.reg.PC)
	assert.Equal(t, state.SP, z80.reg.SP)
	assert.Equal(t, state.I, z80.reg.I)
	assert.Equal(t, state.R, z80.reg.R)
	assert.Equal(t, byte(2), z80.im)
	assert.Equal(t, true, z80.iff1)
	assert.Equal(t, true, z80.iff2)
}

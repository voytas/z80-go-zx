package z80

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EX_AF_AF(t *testing.T) {
	mem := &Memory{
		Cells: []byte{EX_AF_AF, HALT},
	}
	cpu := NewCPU(mem)

	var a, a_, f, f_ byte = 0xcc, 0x55, 0x35, 0x97

	cpu.r.A, cpu.r.A_ = a, a_
	cpu.r.F = f
	cpu.r.F_ = f_
	cpu.Run()

	assert.Equal(t, a_, cpu.r.A)
	assert.Equal(t, a, cpu.r.A_)
	assert.Equal(t, f_, cpu.r.F)
	assert.Equal(t, f, cpu.r.F_)

	cpu.PC = 0
	cpu.Run()
	assert.Equal(t, a, cpu.r.A)
	assert.Equal(t, a_, cpu.r.A_)
	assert.Equal(t, f, cpu.r.F)
	assert.Equal(t, f_, cpu.r.F_)
}

func Test_ADD_HL_DE(t *testing.T) {
	// mem := &Memory{
	// 	Cells: []byte{LD_HL_nn, 0, INC_A, HALT},
	// }
	// cpu := NewCPU(mem)
	// cpu.r.F = f_ALL
	// cpu.Run()
}

func Test_INC_R(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 0, INC_A, HALT},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_ALL
	cpu.Run()

	assert.Equal(t, f_C, cpu.r.F)

	cpu.Reset()
	cpu.r.F = f_ALL & ^f_Z
	mem.Cells[1] = 0xFF
	cpu.Run()

	assert.Equal(t, f_Z|f_H|f_C, cpu.r.F)

	cpu.Reset()
	cpu.r.F = f_N
	mem.Cells[1] = 0x7F
	cpu.Run()

	assert.Equal(t, f_S|f_H|f_PV, cpu.r.F)

	cpu.Reset()
	mem.Cells[1] = 0x92
	cpu.Run()

	assert.Equal(t, f_S, cpu.r.F)

	cpu.Reset()
	mem.Cells[1] = 0x10
	cpu.Run()

	assert.Equal(t, f_NONE, cpu.r.F)
}

func Test_DEC_RR(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_BC_nn, 0x34, 0x12, DEC_BC, HALT},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_ALL
	cpu.Run()

	assert.Equal(t, word(0x1233), cpu.r.getRR(r_BC))

	mem.Cells[1], mem.Cells[2] = 0x00, 0x00
	cpu.PC = 0
	cpu.Run()
}

func Test_LD_RR_nn(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_BC_nn, 0x34, 0x12, HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, word(0x1234), cpu.r.getRR(r_BC))

	mem = &Memory{
		Cells: []byte{LD_SP_nn, 0x34, 0x12, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, word(0x1234), cpu.r.getRR(r_SP))
}

func Test_LD_BC_A(t *testing.T) {
	var n byte = 0x76
	mem := &Memory{
		Cells: []byte{LD_A_n, n, LD_BC_nn, 0x07, 0x00, LD_BC_A, HALT, 0x00},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, n, cpu.mem.Cells[0x07])
}

func Test_LD_A_BC(t *testing.T) {
	var n byte = 0x76
	mem := &Memory{
		Cells: []byte{LD_BC_nn, 0x05, 0x00, LD_A_BC, HALT, n},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, n, cpu.r.A)
}

func Test_LD_R_n(t *testing.T) {
	var a, b, c, d, e, h, l byte = 1, 2, 3, 4, 5, 6, 7
	mem := &Memory{
		Cells: []byte{LD_A_n, a, LD_B_n, b, LD_C_n, c, LD_D_n, d, LD_E_n, e, LD_H_n, h, LD_L_n, l, HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, a, cpu.r.A)
	assert.Equal(t, b, cpu.r.B)
	assert.Equal(t, c, cpu.r.C)
	assert.Equal(t, d, cpu.r.D)
	assert.Equal(t, e, cpu.r.E)
	assert.Equal(t, h, cpu.r.H)
	assert.Equal(t, l, cpu.r.L)
}

func Test_RLCA(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 0x55, RLCA, HALT},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)

	cpu.Reset()
	mem.Cells[1] = 0xAA
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.r.A)
	assert.Equal(t, f_C, cpu.r.F)

	cpu.Reset()
	mem.Cells[1] = 0x00
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)

	cpu.Reset()
	mem.Cells[1] = 0xFF
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.r.A)
	assert.Equal(t, f_C, cpu.r.F)
}

func Test_RRCA(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 0x55, RRCA, HALT},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.r.A)
	assert.Equal(t, f_C, cpu.r.F)

	cpu.Reset()
	mem.Cells[1] = 0xAA
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)

	cpu.Reset()
	mem.Cells[1] = 0x00
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)

	cpu.Reset()
	mem.Cells[1] = 0xFF
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.r.A)
	assert.Equal(t, f_C, cpu.r.F)
}

func Test_DJNZ(t *testing.T) {
	var b byte = 0x20
	var o int8 = -3
	mem := &Memory{
		Cells: []byte{LD_B_n, b, LD_A_n, 0, INC_A, DJNZ, byte(o), HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, b, cpu.r.A)
	assert.Equal(t, byte(0), cpu.r.B)

	b = 0xFF
	o = 1
	mem = &Memory{
		Cells: []byte{LD_B_n, b, LD_A_n, 1, DJNZ, byte(o), HALT, INC_A, JR, 0xFA},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, b, cpu.r.A)
	assert.Equal(t, byte(0), cpu.r.B)
}

func Test_JR(t *testing.T) {
	mem := &Memory{
		Cells: []byte{JR, 3, LD_C_n, 0x11, HALT, LD_D_n, 0x22, HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.C)
	assert.Equal(t, byte(0x22), cpu.r.D)

	mem = &Memory{
		Cells: []byte{JR, 6, HALT, LD_C_n, 0x11, LD_B_n, 0x33, HALT, JR, 0xF9, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x33), cpu.r.B)
	assert.Equal(t, byte(0x11), cpu.r.C)
}

package z80

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NOP(t *testing.T) {
	mem := &Memory{
		Cells: []byte{NOP, NOP, NOP, HALT},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_ALL
	cpu.Run()

	assert.Equal(t, f_ALL, cpu.r.F)
}

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

func Test_ADD_A_x(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 0, ADD_A_n, 0, HALT},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_ALL & ^f_Z
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
	assert.Equal(t, f_Z, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0xFF, LD_H_n, 0x00, LD_L_n, 0x08, ADD_A_HL, HALT, 0x01},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_ALL & ^f_Z
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
	assert.Equal(t, f_Z|f_H|f_C, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x70, LD_L_n, 0x70, ADD_A_L, HALT},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_NONE
	cpu.Run()

	assert.Equal(t, byte(0xE0), cpu.r.A)
	assert.Equal(t, f_S|f_PV, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0xF0, ADD_A_n, 0xB0, HALT},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_NONE
	cpu.Run()

	assert.Equal(t, byte(0xA0), cpu.r.A)
	assert.Equal(t, f_S|f_C, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x8f, ADD_A_n, 0x81, HALT},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_NONE
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.r.A)
	assert.Equal(t, f_H|f_PV|f_C, cpu.r.F)
}

func Test_ADC_A_x(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 0, LD_B_n, 0xFF, ADC_A_B, HALT},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
	assert.Equal(t, f_Z|f_H|f_C, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x0F, LD_B_n, 0x00, ADC_A_B, HALT},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.r.A)
	assert.Equal(t, f_H, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x0F, LD_L_n, 0x06, ADC_A_HL, HALT, 0x70},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x80), cpu.r.A)
	assert.Equal(t, f_S|f_H|f_PV, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x0F, LD_B_n, 0x69, ADC_A_B, HALT},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x79), cpu.r.A)
	assert.Equal(t, f_H, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x0E, LD_B_n, 0x01, ADC_A_B, HALT},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.r.A)
	assert.Equal(t, f_H, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x0E, LD_B_n, 0x01, ADC_A_B, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x0F), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)
}

func Test_ADD_HL_RR(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_HL_nn, 0xFF, 0xFF, LD_BC_nn, 0x01, 0, ADD_HL_BC, HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.H)
	assert.Equal(t, byte(0), cpu.r.L)
	assert.Equal(t, f_H|f_C, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_HL_nn, 0x41, 0x42, LD_DE_nn, 0x11, 0x11, ADD_HL_DE, HALT},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_ALL
	cpu.Run()

	assert.Equal(t, byte(0x53), cpu.r.H)
	assert.Equal(t, byte(0x52), cpu.r.L)
	assert.Equal(t, f_S|f_Z|f_PV, cpu.r.F)
}

func Test_SUB_x(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 0, LD_B_n, 0x01, SUB_B, HALT},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_Z
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.r.A)
	assert.Equal(t, f_S|f_H|f_N|f_C, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x20, SUB_A, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.r.A)
	assert.Equal(t, f_Z|f_N, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x90, LD_H_n, 0x20, SUB_H, HALT},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_Z
	cpu.Run()

	assert.Equal(t, byte(0x70), cpu.r.A)
	assert.Equal(t, f_PV|f_N, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x7F, LD_L_n, 0x06, SUB_HL, HALT, 0x80},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.r.A)
	assert.Equal(t, f_S|f_PV|f_N|f_C, cpu.r.F)
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

func Test_INC_RR(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_BC_nn, 0x34, 0x12, INC_BC, HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, word(0x1235), cpu.r.getRR(r_BC))
}

func Test_INC_mHL(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_HL_nn, 0x05, 0x00, INC_mHL, HALT, 0xFF},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_S | f_PV | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x00), mem.Cells[5])
	assert.Equal(t, f_Z|f_H|f_C, cpu.r.F)

	cpu.Reset()
	mem.Cells[5] = 0x7F
	cpu.Run()

	assert.Equal(t, byte(0x80), mem.Cells[5])
	assert.Equal(t, f_S|f_H|f_PV, cpu.r.F)

	cpu.Reset()
	mem.Cells[5] = 0x20
	cpu.Run()

	assert.Equal(t, byte(0x21), mem.Cells[5])
	assert.Equal(t, f_NONE, cpu.r.F)
}

func Test_DEC_R(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 1, DEC_A, HALT},
	}

	cpu := NewCPU(mem)
	cpu.r.F = f_NONE
	cpu.Run()

	assert.Equal(t, f_Z|f_N, cpu.r.F)

	cpu.Reset()
	cpu.r.F = f_ALL & ^(f_Z | f_H | f_N)
	mem.Cells[1] = 0
	cpu.Run()

	cpu.Reset()
	cpu.r.F = f_Z | f_S
	mem.Cells[1] = 0x80
	cpu.Run()

	assert.Equal(t, f_H|f_PV|f_N, cpu.r.F)

	cpu.Reset()
	cpu.r.F = f_ALL
	mem.Cells[1] = 0xab
	cpu.Run()

	assert.Equal(t, f_S|f_N|f_C, cpu.r.F)
}

func Test_DEC_RR(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_BC_nn, 0x34, 0x12, DEC_BC, HALT},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_ALL
	cpu.Run()

	assert.Equal(t, word(0x1233), cpu.r.getRR(r_BC))
}

func Test_DEC_mHL(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_HL_nn, 0x05, 0x00, DEC_mHL, HALT, 0x00},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0xFF), mem.Cells[5])
	assert.Equal(t, f_S|f_H|f_N|f_C, cpu.r.F)

	cpu.Reset()
	mem.Cells[5] = 0x01
	cpu.Run()

	assert.Equal(t, byte(0x00), mem.Cells[5])
	assert.Equal(t, f_Z|f_N, cpu.r.F)

	cpu.Reset()
	mem.Cells[5] = 0x80
	cpu.Run()

	assert.Equal(t, byte(0x7F), mem.Cells[5])
	assert.Equal(t, f_PV|f_H|f_N, cpu.r.F)
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

func Test_LD_mm_HL(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_HL_nn, 0x3A, 0x48, LD_mm_HL, 0x07, 0x00, HALT, 0, 0},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, cpu.r.H, mem.Cells[8])
	assert.Equal(t, cpu.r.L, mem.Cells[7])
}

func Test_LD_HL_mm(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_HL_mm, 0x04, 0x00, HALT, 0x34, 0x12},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x12), cpu.r.H)
	assert.Equal(t, byte(0x34), cpu.r.L)
}

func Test_LD_mHL_n(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_HL_nn, 0x06, 0x00, LD_mHL_n, 0xAB, HALT, 0x00},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.mem.Cells[6])
}

func Test_LD_mm_A(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 0x9F, LD_mm_A, 0x05, 0x00, 0x00, HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, cpu.r.A, mem.Cells[5])
}

func Test_LD_A_mm(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_mm, 0x04, 0x00, HALT, 0xDE},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xDE), cpu.r.A)
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

func Test_LD_DE_A(t *testing.T) {
	var n byte = 0x76
	mem := &Memory{
		Cells: []byte{LD_A_n, n, LD_DE_nn, 0x07, 0x00, LD_DE_A, HALT, 0x00},
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

func Test_LD_A_DE(t *testing.T) {
	var n byte = 0x76
	mem := &Memory{
		Cells: []byte{LD_DE_nn, 0x05, 0x00, LD_A_DE, HALT, n},
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

func Test_LD_R_R(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 0x56, LD_B_A, LD_C_B, LD_D_C, LD_E_D, LD_H_E, LD_L_H, LD_A_n, 0, LD_A_B, HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x56), cpu.r.A)
	assert.Equal(t, byte(0x56), cpu.r.B)
	assert.Equal(t, byte(0x56), cpu.r.C)
	assert.Equal(t, byte(0x56), cpu.r.D)
	assert.Equal(t, byte(0x56), cpu.r.E)
	assert.Equal(t, byte(0x56), cpu.r.H)
	assert.Equal(t, byte(0x56), cpu.r.L)
}

func Test_LD_R_HL(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_HL_nn, 0x06, 0x00, LD_A_HL, LD_L_HL, HALT, 0xA7},
	}

	cpu := NewCPU(mem)
	cpu.Run()
	assert.Equal(t, byte(0xA7), cpu.r.A)
	assert.Equal(t, byte(0xA7), cpu.r.L)
}

func Test_LD_HL_R(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_D_n, 0x99, LD_HL_nn, 0x07, 0x00, LD_HL_D, HALT, 0x00},
	}

	cpu := NewCPU(mem)
	cpu.Run()
	assert.Equal(t, byte(0x99), cpu.mem.Cells[7])
}

func Test_CPL(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 0x5B, CPL, HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xA4), cpu.r.A)
	assert.Equal(t, f_H|f_N, cpu.r.F)
}

func Test_SCF(t *testing.T) {
	mem := &Memory{
		Cells: []byte{SCF, HALT},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_S | f_Z | f_H | f_PV | f_N
	cpu.Run()

	assert.Equal(t, f_S|f_Z|f_PV|f_C, cpu.r.F)
}

func Test_CCF(t *testing.T) {
	mem := &Memory{
		Cells: []byte{CCF, HALT},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_ALL
	cpu.Run()

	assert.Equal(t, f_S|f_Z|f_PV, cpu.r.F)

	cpu.Reset()
	cpu.r.F = f_Z | f_N | f_C
	cpu.Run()

	assert.Equal(t, f_Z|f_H, cpu.r.F)
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

func Test_RLA(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 0x80, RLA, HALT},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.r.A)
	assert.Equal(t, f_C, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x55, RLA, HALT},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_ALL
	cpu.Run()

	assert.Equal(t, byte(0xab), cpu.r.A)
	assert.Equal(t, f_S|f_Z|f_PV, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x88, RLA, LD_B_A, RLA, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.r.B)
	assert.Equal(t, byte(0x21), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)
}

func Test_RRA(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 0x80, RRA, HALT},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0xC0), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x55, RRA, HALT},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_ALL
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.r.A)
	assert.Equal(t, f_S|f_Z|f_PV|f_C, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x89, RRA, LD_B_A, RRA, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x44), cpu.r.B)
	assert.Equal(t, byte(0xA2), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)
}

func Test_DAA(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 0x9A, DAA, HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
	assert.Equal(t, f_Z|f_H|f_PV|f_C, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x99, DAA, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x99), cpu.r.A)
	assert.Equal(t, f_S|f_PV, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x8F, DAA, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x95), cpu.r.A)
	assert.Equal(t, f_S|f_H|f_PV, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0x8F, DAA, HALT},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_N
	cpu.Run()

	assert.Equal(t, byte(0x89), cpu.r.A)
	assert.Equal(t, f_S|f_N, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0xCA, DAA, HALT},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_N
	cpu.Run()

	assert.Equal(t, byte(0x64), cpu.r.A)
	assert.Equal(t, f_N|f_C, cpu.r.F)

	mem = &Memory{
		Cells: []byte{LD_A_n, 0xC5, DAA, HALT},
	}
	cpu = NewCPU(mem)
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x5F), cpu.r.A)
	assert.Equal(t, f_H|f_PV|f_N|f_C, cpu.r.F)
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

func Test_JR_Z(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 2, DEC_A, JR_Z, 0x02, LD_B_n, 0xab, HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xab), cpu.r.B)

	mem = &Memory{
		Cells: []byte{LD_A_n, 1, DEC_A, JR_Z, 0x02, LD_B_n, 0xab, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	mem = &Memory{
		Cells: []byte{LD_A_n, 1, DEC_A, JR_Z, 0xFD, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.r.A)
}

func Test_JR_NZ(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_A_n, 2, DEC_A, JR_NZ, 0x02, LD_B_n, 0xab, HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.B)

	mem = &Memory{
		Cells: []byte{LD_A_n, 1, DEC_A, JR_NZ, 0x02, LD_B_n, 0xab, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xab), cpu.r.B)

	mem = &Memory{
		Cells: []byte{LD_A_n, 2, DEC_A, JR_NZ, 0xFD, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
}

func Test_JR_C(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_B_n, 0xAB, DEC_A, ADD_A_n, 1, JR_C, 1, LD_B_A, HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.r.B)

	mem = &Memory{
		Cells: []byte{LD_B_n, 0xAB, INC_A, ADD_A_n, 1, JR_C, 1, LD_B_A, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(2), cpu.r.B)

	mem = &Memory{
		Cells: []byte{DEC_A, ADD_A_n, 1, JR_C, 0xFC, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(1), cpu.r.A)
}

func Test_JR_NC(t *testing.T) {
	mem := &Memory{
		Cells: []byte{LD_B_n, 0xAB, INC_A, ADD_A_n, 1, JR_NC, 1, LD_B_A, HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.r.B)

	mem = &Memory{
		Cells: []byte{LD_B_n, 0xAB, DEC_A, ADD_A_n, 1, JR_NC, 1, LD_B_A, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.B)

	mem = &Memory{
		Cells: []byte{ADD_A_n, 1, JR_NC, 0xFC, HALT},
	}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
}

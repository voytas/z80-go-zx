package z80

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NOP(t *testing.T) {
	mem := &BasicMemory{cells: []byte{NOP, NOP, NOP, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_ALL
	cpu.Run()

	assert.Equal(t, f_ALL, cpu.r.F)
}

func Test_EX_AF_AF(t *testing.T) {
	mem := &BasicMemory{cells: []byte{EX_AF_AF, HALT}}
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

func Test_EXX(t *testing.T) {
	mem := &BasicMemory{cells: []byte{EXX, HALT}}
	cpu := NewCPU(mem)
	cpu.r.B, cpu.r.C, cpu.r.B_, cpu.r.C_ = 0x01, 0x02, 0x03, 0x04
	cpu.r.D, cpu.r.E, cpu.r.D_, cpu.r.E_ = 0x05, 0x06, 0x07, 0x08
	cpu.r.H, cpu.r.L, cpu.r.H_, cpu.r.L_ = 0x09, 0x0A, 0x0B, 0x0C
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.r.B_)
	assert.Equal(t, byte(0x02), cpu.r.C_)
	assert.Equal(t, byte(0x03), cpu.r.B)
	assert.Equal(t, byte(0x04), cpu.r.C)
	assert.Equal(t, byte(0x05), cpu.r.D_)
	assert.Equal(t, byte(0x06), cpu.r.E_)
	assert.Equal(t, byte(0x07), cpu.r.D)
	assert.Equal(t, byte(0x08), cpu.r.E)
	assert.Equal(t, byte(0x09), cpu.r.H_)
	assert.Equal(t, byte(0x0A), cpu.r.L_)
	assert.Equal(t, byte(0x0B), cpu.r.H)
	assert.Equal(t, byte(0x0C), cpu.r.L)

}

func Test_ADD_A_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0, ADD_A_n, 0, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_ALL & ^f_Z
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
	assert.Equal(t, f_Z, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0xFF, LD_H_n, 0x00, LD_L_n, 0x08, ADD_A_HL, HALT, 0x01}}
	cpu = NewCPU(mem)
	cpu.r.F = f_ALL & ^f_Z
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
	assert.Equal(t, f_Z|f_H|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x70, LD_L_n, 0x70, ADD_A_L, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_NONE
	cpu.Run()

	assert.Equal(t, byte(0xE0), cpu.r.A)
	assert.Equal(t, f_S|f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0xF0, ADD_A_n, 0xB0, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_NONE
	cpu.Run()

	assert.Equal(t, byte(0xA0), cpu.r.A)
	assert.Equal(t, f_S|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x8f, ADD_A_n, 0x81, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_NONE
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.r.A)
	assert.Equal(t, f_H|f_P|f_C, cpu.r.F)
}

func Test_ADC_A_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0x20, ADC_A_n, 0x20, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x41), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0, LD_B_n, 0xFF, ADC_A_B, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
	assert.Equal(t, f_Z|f_H|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x0F, LD_B_n, 0x00, ADC_A_B, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.r.A)
	assert.Equal(t, f_H, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x0F, LD_L_n, 0x06, ADC_A_HL, HALT, 0x70}}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x80), cpu.r.A)
	assert.Equal(t, f_S|f_H|f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x0F, LD_B_n, 0x69, ADC_A_B, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x79), cpu.r.A)
	assert.Equal(t, f_H, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x0E, LD_B_n, 0x01, ADC_A_B, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.r.A)
	assert.Equal(t, f_H, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x0E, LD_B_n, 0x01, ADC_A_B, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x0F), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)
}

func Test_ADD_HL_RR(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_HL_nn, 0xFF, 0xFF, LD_BC_nn, 0x01, 0, ADD_HL_BC, HALT}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.H)
	assert.Equal(t, byte(0), cpu.r.L)
	assert.Equal(t, f_H|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_HL_nn, 0x41, 0x42, LD_DE_nn, 0x11, 0x11, ADD_HL_DE, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_ALL
	cpu.Run()

	assert.Equal(t, byte(0x53), cpu.r.H)
	assert.Equal(t, byte(0x52), cpu.r.L)
	assert.Equal(t, f_S|f_Z|f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_HL_nn, 0x41, 0x42, ADD_HL_HL, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x84), cpu.r.H)
	assert.Equal(t, byte(0x82), cpu.r.L)
	assert.Equal(t, f_NONE, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_HL_nn, 0xFE, 0xFF, LD_SP_nn, 0x02, 0, ADD_HL_SP, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.H)
	assert.Equal(t, byte(0), cpu.r.L)
	assert.Equal(t, f_H|f_C, cpu.r.F)
}

func Test_SUB_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0, SUB_n, 0x01, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_Z
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.r.A)
	assert.Equal(t, f_S|f_H|f_N|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x20, SUB_A, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.r.A)
	assert.Equal(t, f_Z|f_N, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x90, LD_H_n, 0x20, SUB_H, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_Z
	cpu.Run()

	assert.Equal(t, byte(0x70), cpu.r.A)
	assert.Equal(t, f_P|f_N, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x7F, LD_L_n, 0x06, SUB_HL, HALT, 0x80}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.r.A)
	assert.Equal(t, f_S|f_P|f_N|f_C, cpu.r.F)
}

func Test_CP_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0, LD_B_n, 0x01, CP_B, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_Z
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
	assert.Equal(t, f_S|f_H|f_N|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x20, CP_A, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x20), cpu.r.A)
	assert.Equal(t, f_Z|f_N, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x90, LD_H_n, 0x20, CP_H, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_Z
	cpu.Run()

	assert.Equal(t, byte(0x90), cpu.r.A)
	assert.Equal(t, f_P|f_N, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x7F, LD_L_n, 0x06, CP_HL, HALT, 0x80}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x7F), cpu.r.A)
	assert.Equal(t, f_S|f_P|f_N|f_C, cpu.r.F)
}

func Test_SBC_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0x01, LD_B_n, 0x01, SBC_A_B, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.r.A)
	assert.Equal(t, f_S|f_H|f_N|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x7F, LD_L_n, 0x80, SBC_A_L, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0xFE), cpu.r.A)
	assert.Equal(t, f_S|f_P|f_N|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x02, LD_D_n, 0x01, SBC_A_D, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.r.A)
	assert.Equal(t, f_Z|f_N, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x81, LD_L_n, 0x06, SBC_A_HL, HALT, 0x01}}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x7F), cpu.r.A)
	assert.Equal(t, f_H|f_P|f_N, cpu.r.F)
}

func Test_AND_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0x0F, LD_B_n, 0xF0, AND_B, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
	assert.Equal(t, f_Z|f_H|f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x8F, LD_D_n, 0xF3, AND_D, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x83), cpu.r.A)
	assert.Equal(t, f_S|f_H, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0xFF, LD_L_n, 0x06, AND_HL, HALT, 0x81}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x81), cpu.r.A)
	assert.Equal(t, f_S|f_H|f_P, cpu.r.F)
}

func Test_OR_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0x00, LD_B_n, 0x00, OR_B, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_S | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
	assert.Equal(t, f_Z|f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x8A, LD_L_n, 0x06, OR_HL, HALT, 0x85}}
	cpu = NewCPU(mem)
	cpu.r.F = f_S | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x8F), cpu.r.A)
	assert.Equal(t, f_S, cpu.r.F)
}

func Test_XOR_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0x1F, LD_B_n, 0x1F, XOR_B, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
	assert.Equal(t, f_Z|f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x1F, LD_L_n, 0x06, XOR_HL, HALT, 0x8F}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x90), cpu.r.A)
	assert.Equal(t, f_S|f_P, cpu.r.F)
}

func Test_INC_R(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0, INC_A, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_ALL
	cpu.Run()

	assert.Equal(t, f_C, cpu.r.F)

	cpu.Reset()
	cpu.r.F = f_ALL & ^f_Z
	mem.cells[1] = 0xFF
	cpu.Run()

	assert.Equal(t, f_Z|f_H|f_C, cpu.r.F)

	cpu.Reset()
	cpu.r.F = f_N
	mem.cells[1] = 0x7F
	cpu.Run()

	assert.Equal(t, f_S|f_H|f_P, cpu.r.F)

	cpu.Reset()
	mem.cells[1] = 0x92
	cpu.Run()

	assert.Equal(t, f_S, cpu.r.F)

	cpu.Reset()
	mem.cells[1] = 0x10
	cpu.Run()

	assert.Equal(t, f_NONE, cpu.r.F)
}

func Test_INC_RR(t *testing.T) {
	mem := &BasicMemory{
		cells: []byte{
			LD_BC_nn, 0x34, 0x12, INC_BC, LD_DE_nn, 0x35, 0x13, INC_DE,
			LD_HL_nn, 0x36, 0x14, INC_HL, LD_SP_nn, 0x37, 0x15, INC_SP,
			HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, word(0x1235), cpu.r.getBC())
	assert.Equal(t, word(0x1336), cpu.r.getDE())
	assert.Equal(t, word(0x1437), cpu.r.getHL())
	assert.Equal(t, word(0x1538), cpu.r.SP)
}

func Test_INC_mHL(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_HL_nn, 0x05, 0x00, INC_mHL, HALT, 0xFF}}
	cpu := NewCPU(mem)
	cpu.r.F = f_S | f_P | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x00), mem.cells[5])
	assert.Equal(t, f_Z|f_H|f_C, cpu.r.F)

	cpu.Reset()
	mem.cells[5] = 0x7F
	cpu.Run()

	assert.Equal(t, byte(0x80), mem.cells[5])
	assert.Equal(t, f_S|f_H|f_P, cpu.r.F)

	cpu.Reset()
	mem.cells[5] = 0x20
	cpu.Run()

	assert.Equal(t, byte(0x21), mem.cells[5])
	assert.Equal(t, f_NONE, cpu.r.F)
}

func Test_DEC_R(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 1, DEC_A, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_NONE
	cpu.Run()

	assert.Equal(t, f_Z|f_N, cpu.r.F)

	cpu.Reset()
	cpu.r.F = f_ALL & ^(f_Z | f_H | f_N)
	mem.cells[1] = 0
	cpu.Run()

	cpu.Reset()
	cpu.r.F = f_Z | f_S
	mem.cells[1] = 0x80
	cpu.Run()

	assert.Equal(t, f_H|f_P|f_N, cpu.r.F)

	cpu.Reset()
	cpu.r.F = f_ALL
	mem.cells[1] = 0xab
	cpu.Run()

	assert.Equal(t, f_S|f_N|f_C, cpu.r.F)
}

func Test_DEC_RR(t *testing.T) {
	mem := &BasicMemory{
		cells: []byte{
			LD_BC_nn, 0x34, 0x12, DEC_BC, LD_DE_nn, 0x35, 0x13, DEC_DE,
			LD_HL_nn, 0x36, 0x14, DEC_HL, LD_SP_nn, 0x37, 0x15, DEC_SP,
			HALT},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, word(0x1233), cpu.r.getBC())
	assert.Equal(t, word(0x1334), cpu.r.getDE())
	assert.Equal(t, word(0x1435), cpu.r.getHL())
	assert.Equal(t, word(0x1536), cpu.r.SP)
}

func Test_DEC_mHL(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_HL_nn, 0x05, 0x00, DEC_mHL, HALT, 0x00}}
	cpu := NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0xFF), mem.cells[5])
	assert.Equal(t, f_S|f_H|f_N|f_C, cpu.r.F)

	cpu.Reset()
	mem.cells[5] = 0x01
	cpu.Run()

	assert.Equal(t, byte(0x00), mem.cells[5])
	assert.Equal(t, f_Z|f_N, cpu.r.F)

	cpu.Reset()
	mem.cells[5] = 0x80
	cpu.Run()

	assert.Equal(t, byte(0x7F), mem.cells[5])
	assert.Equal(t, f_P|f_H|f_N, cpu.r.F)
}

func Test_LD_RR_nn(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_BC_nn, 0x34, 0x12, HALT}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, word(0x1234), cpu.r.getBC())

	mem = &BasicMemory{cells: []byte{LD_SP_nn, 0x34, 0x12, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, word(0x1234), cpu.r.SP)
}

func Test_LD_mm_HL(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_HL_nn, 0x3A, 0x48, LD_mm_HL, 0x07, 0x00, HALT, 0, 0}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, cpu.r.H, mem.cells[8])
	assert.Equal(t, cpu.r.L, mem.cells[7])
}

func Test_LD_HL_mm(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_HL_mm, 0x04, 0x00, HALT, 0x34, 0x12}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x12), cpu.r.H)
	assert.Equal(t, byte(0x34), cpu.r.L)
}

func Test_LD_mHL_n(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_HL_nn, 0x06, 0x00, LD_mHL_n, 0xAB, HALT, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.mem.read(6))
}

func Test_LD_mm_A(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0x9F, LD_mm_A, 0x06, 0x00, HALT, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, cpu.r.A, mem.cells[6])
}

func Test_LD_A_mm(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_mm, 0x04, 0x00, HALT, 0xDE}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xDE), cpu.r.A)
}

func Test_LD_BC_A(t *testing.T) {
	var n byte = 0x76
	mem := &BasicMemory{cells: []byte{LD_A_n, n, LD_BC_nn, 0x07, 0x00, LD_BC_A, HALT, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, n, cpu.mem.read(7))
}

func Test_LD_DE_A(t *testing.T) {
	var n byte = 0x76
	mem := &BasicMemory{cells: []byte{LD_A_n, n, LD_DE_nn, 0x07, 0x00, LD_DE_A, HALT, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, n, cpu.mem.read(7))
}

func Test_LD_A_BC(t *testing.T) {
	var n byte = 0x76
	mem := &BasicMemory{cells: []byte{LD_BC_nn, 0x05, 0x00, LD_A_BC, HALT, n}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, n, cpu.r.A)
}

func Test_LD_A_DE(t *testing.T) {
	var n byte = 0x76
	mem := &BasicMemory{cells: []byte{LD_DE_nn, 0x05, 0x00, LD_A_DE, HALT, n}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, n, cpu.r.A)
}

func Test_LD_R_n(t *testing.T) {
	var a, b, c, d, e, h, l byte = 1, 2, 3, 4, 5, 6, 7
	mem := &BasicMemory{
		cells: []byte{LD_A_n, a, LD_B_n, b, LD_C_n, c, LD_D_n, d, LD_E_n, e, LD_H_n, h, LD_L_n, l, HALT},
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
	mem := &BasicMemory{
		cells: []byte{LD_A_n, 0x56, LD_B_A, LD_C_B, LD_D_C, LD_E_D, LD_H_E, LD_L_H, LD_A_n, 0, LD_A_B, HALT},
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
	mem := &BasicMemory{cells: []byte{LD_HL_nn, 0x06, 0x00, LD_A_HL, LD_L_HL, HALT, 0xA7}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xA7), cpu.r.A)
	assert.Equal(t, byte(0xA7), cpu.r.L)
}

func Test_LD_HL_R(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_D_n, 0x99, LD_HL_nn, 0x07, 0x00, LD_HL_D, HALT, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x99), cpu.mem.read(7))
}

func Test_CPL(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0x5B, CPL, HALT}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xA4), cpu.r.A)
	assert.Equal(t, f_H|f_N, cpu.r.F)
}

func Test_SCF(t *testing.T) {
	mem := &BasicMemory{cells: []byte{SCF, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_S | f_Z | f_H | f_P | f_N
	cpu.Run()

	assert.Equal(t, f_S|f_Z|f_P|f_C, cpu.r.F)
}

func Test_CCF(t *testing.T) {
	mem := &BasicMemory{cells: []byte{CCF, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_ALL
	cpu.Run()

	assert.Equal(t, f_S|f_Z|f_P, cpu.r.F)

	cpu.Reset()
	cpu.r.F = f_Z | f_N | f_C
	cpu.Run()

	assert.Equal(t, f_Z|f_H, cpu.r.F)
}

func Test_RLCA(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0x55, RLCA, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)

	cpu.Reset()
	mem.cells[1] = 0xAA
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.r.A)
	assert.Equal(t, f_C, cpu.r.F)

	cpu.Reset()
	mem.cells[1] = 0x00
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)

	cpu.Reset()
	mem.cells[1] = 0xFF
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.r.A)
	assert.Equal(t, f_C, cpu.r.F)
}

func Test_RRCA(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0x55, RRCA, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.r.A)
	assert.Equal(t, f_C, cpu.r.F)

	cpu.Reset()
	mem.cells[1] = 0xAA
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)

	cpu.Reset()
	mem.cells[1] = 0x00
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)

	cpu.Reset()
	mem.cells[1] = 0xFF
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.r.A)
	assert.Equal(t, f_C, cpu.r.F)
}

func Test_RLA(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0x80, RLA, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.r.A)
	assert.Equal(t, f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x55, RLA, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_ALL
	cpu.Run()

	assert.Equal(t, byte(0xab), cpu.r.A)
	assert.Equal(t, f_S|f_Z|f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x88, RLA, LD_B_A, RLA, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.r.B)
	assert.Equal(t, byte(0x21), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)
}

func Test_RRA(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0x80, RRA, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0xC0), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x55, RRA, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_ALL
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.r.A)
	assert.Equal(t, f_S|f_Z|f_P|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x89, RRA, LD_B_A, RRA, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x44), cpu.r.B)
	assert.Equal(t, byte(0xA2), cpu.r.A)
	assert.Equal(t, f_NONE, cpu.r.F)
}

func Test_DAA(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 0x9A, DAA, HALT}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
	assert.Equal(t, f_Z|f_H|f_P|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x99, DAA, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x99), cpu.r.A)
	assert.Equal(t, f_S|f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x8F, DAA, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x95), cpu.r.A)
	assert.Equal(t, f_S|f_H|f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x8F, DAA, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_N
	cpu.Run()

	assert.Equal(t, byte(0x89), cpu.r.A)
	assert.Equal(t, f_S|f_N, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0xCA, DAA, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_N
	cpu.Run()

	assert.Equal(t, byte(0x64), cpu.r.A)
	assert.Equal(t, f_N|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0xC5, DAA, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x5F), cpu.r.A)
	assert.Equal(t, f_H|f_P|f_N|f_C, cpu.r.F)
}

func Test_DJNZ(t *testing.T) {
	var b byte = 0x20
	var o int8 = -3
	mem := &BasicMemory{cells: []byte{LD_B_n, b, LD_A_n, 0, INC_A, DJNZ, byte(o), HALT}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, b, cpu.r.A)
	assert.Equal(t, byte(0), cpu.r.B)

	b = 0xFF
	o = 1
	mem = &BasicMemory{cells: []byte{LD_B_n, b, LD_A_n, 1, DJNZ, byte(o), HALT, INC_A, JR_o, 0xFA}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, b, cpu.r.A)
	assert.Equal(t, byte(0), cpu.r.B)
}

func Test_JR_o(t *testing.T) {
	mem := &BasicMemory{cells: []byte{JR_o, 3, LD_C_n, 0x11, HALT, LD_D_n, 0x22, HALT}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.C)
	assert.Equal(t, byte(0x22), cpu.r.D)

	mem = &BasicMemory{cells: []byte{JR_o, 6, HALT, LD_C_n, 0x11, LD_B_n, 0x33, HALT, JR_o, 0xF9, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x33), cpu.r.B)
	assert.Equal(t, byte(0x11), cpu.r.C)
}

func Test_JR_Z_o(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 2, DEC_A, JR_Z_o, 0x02, LD_B_n, 0xab, HALT}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xab), cpu.r.B)

	mem = &BasicMemory{cells: []byte{LD_A_n, 1, DEC_A, JR_Z_o, 0x02, LD_B_n, 0xab, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	mem = &BasicMemory{cells: []byte{LD_A_n, 1, DEC_A, JR_Z_o, 0xFD, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.r.A)
}

func Test_JR_NZ_o(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_A_n, 2, DEC_A, JR_NZ_o, 0x02, LD_B_n, 0xab, HALT}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.B)

	mem = &BasicMemory{cells: []byte{LD_A_n, 1, DEC_A, JR_NZ_o, 0x02, LD_B_n, 0xab, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xab), cpu.r.B)

	mem = &BasicMemory{cells: []byte{LD_A_n, 2, DEC_A, JR_NZ_o, 0xFD, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
}

func Test_JR_C(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_B_n, 0xAB, DEC_A, ADD_A_n, 1, JR_C, 1, LD_B_A, HALT}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.r.B)

	mem = &BasicMemory{cells: []byte{LD_B_n, 0xAB, INC_A, ADD_A_n, 1, JR_C, 1, LD_B_A, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(2), cpu.r.B)

	mem = &BasicMemory{cells: []byte{DEC_A, ADD_A_n, 1, JR_C, 0xFC, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(1), cpu.r.A)
}

func Test_JR_NC_o(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_B_n, 0xAB, INC_A, ADD_A_n, 1, JR_NC_o, 1, LD_B_A, HALT}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.r.B)

	mem = &BasicMemory{cells: []byte{LD_B_n, 0xAB, DEC_A, ADD_A_n, 1, JR_NC_o, 1, LD_B_A, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.B)

	mem = &BasicMemory{cells: []byte{ADD_A_n, 1, JR_NC_o, 0xFC, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.r.A)
}

func Test_JP_nn(t *testing.T) {
	mem := &BasicMemory{cells: []byte{JP_nn, 0x06, 0x00, LD_A_n, 0xAA, HALT, LD_A_n, 0x55, HALT, 0x07, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.r.A)
}

func Test_JP_cc_nn(t *testing.T) {
	var tests = []struct {
		jp       byte
		flag     byte
		expected byte
	}{
		{JP_C_nn, f_C, 0x55}, {JP_NC_nn, f_NONE, 0x55}, {JP_Z_nn, f_Z, 0x55}, {JP_NZ_nn, f_NONE, 0x55},
		{JP_M_nn, f_S, 0x55}, {JP_P_nn, f_NONE, 0x55}, {JP_PE_nn, f_P, 0x55}, {JP_PO_nn, f_NONE, 0x55},
		{JP_C_nn, f_NONE, 0xAA}, {JP_NC_nn, f_C, 0xAA}, {JP_Z_nn, f_NONE, 0xAA}, {JP_NZ_nn, f_Z, 0xAA},
		{JP_M_nn, f_NONE, 0xAA}, {JP_P_nn, f_S, 0xAA}, {JP_PE_nn, f_NONE, 0xAA}, {JP_PO_nn, f_P, 0xAA},
	}

	for _, test := range tests {
		mem := &BasicMemory{
			cells: []byte{test.jp, 0x06, 0x00, LD_A_n, 0xAA, HALT, LD_A_n, 0x55, HALT, 0x07, 0x00},
		}

		cpu := NewCPU(mem)
		cpu.r.F = test.flag
		cpu.Run()

		assert.Equal(t, byte(test.expected), cpu.r.A)
	}
}

func Test_CALL_cc_nn(t *testing.T) {
	var tests = []struct {
		call     byte
		flag     byte
		expected byte
	}{
		{CALL_C_nn, f_C, 0x55}, {CALL_NC_nn, f_NONE, 0x55}, {CALL_Z_nn, f_Z, 0x55}, {CALL_NZ_nn, f_NONE, 0x55},
		{CALL_M_nn, f_S, 0x55}, {CALL_P_nn, f_NONE, 0x55}, {CALL_PE_nn, f_P, 0x55}, {CALL_PO_nn, f_NONE, 0x55},
		{CALL_C_nn, f_NONE, 0xAA}, {CALL_NC_nn, f_C, 0xAA}, {CALL_Z_nn, f_NONE, 0xAA}, {CALL_NZ_nn, f_Z, 0xAA},
		{CALL_M_nn, f_NONE, 0xAA}, {CALL_P_nn, f_S, 0xAA}, {CALL_PE_nn, f_NONE, 0xAA}, {CALL_PO_nn, f_P, 0xAA},
		{CALL_nn, f_NONE, 0x55},
	}

	for _, test := range tests {
		mem := &BasicMemory{
			cells: []byte{LD_SP_nn, 0x10, 0x00, test.call, 0x09, 0x00, LD_A_n, 0xAA, HALT, LD_A_n, 0x55, HALT, 0xFF, 0xFF, 0xFF, 0xFF},
		}

		cpu := NewCPU(mem)
		cpu.r.F = test.flag
		cpu.Run()

		assert.Equal(t, byte(test.expected), cpu.r.A)
		if cpu.r.A == 0x55 {
			assert.Equal(t, word(0x0E), cpu.r.SP)
			assert.Equal(t, byte(0), mem.cells[15])
			assert.Equal(t, byte(0x06), mem.cells[14])
		} else {
			assert.Equal(t, word(0x10), cpu.r.SP)
		}
	}
}

func Test_RET(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_SP_nn, 0x0A, 0x00, RET, LD_A_n, 0xAA, HALT, LD_A_n, 0x55, HALT, 0x07, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.r.A)
}

func Test_RET_cc(t *testing.T) {
	var tests = []struct {
		ret      byte
		flag     byte
		expected byte
	}{
		{RET_C, f_C, 0x55}, {RET_NC, f_NONE, 0x55}, {RET_Z, f_Z, 0x55}, {RET_NZ, f_NONE, 0x55},
		{RET_M, f_S, 0x55}, {RET_P, f_NONE, 0x55}, {RET_PE, f_P, 0x55}, {RET_PO, f_NONE, 0x55},
		{RET_C, f_NONE, 0xAA}, {RET_NC, f_C, 0xAA}, {RET_Z, f_NONE, 0xAA}, {RET_NZ, f_Z, 0xAA},
		{RET_M, f_NONE, 0xAA}, {RET_P, f_S, 0xAA}, {RET_PE, f_NONE, 0xAA}, {RET_PO, f_P, 0xAA},
	}

	for _, test := range tests {
		mem := &BasicMemory{
			cells: []byte{LD_SP_nn, 0x0A, 0x00, test.ret, LD_A_n, 0xAA, HALT, LD_A_n, 0x55, HALT, 0x07, 0x00},
		}
		cpu := NewCPU(mem)
		cpu.r.F = test.flag
		cpu.Run()

		assert.Equal(t, byte(test.expected), cpu.r.A)
	}
}

func Test_RST_xx(t *testing.T) {
	mem := &BasicMemory{
		cells: []byte{
			LD_A_n, 0x01, RET, 0, 0, 0, 0, 0,
			ADD_A_n, 0x02, RET, 0, 0, 0, 0, 0,
			ADD_A_n, 0x04, RET, 0, 0, 0, 0, 0,
			ADD_A_n, 0x08, RET, 0, 0, 0, 0, 0,
			ADD_A_n, 0x10, RET, 0, 0, 0, 0, 0,
			ADD_A_n, 0x20, RET, 0, 0, 0, 0, 0,
			ADD_A_n, 0x40, RET, 0, 0, 0, 0, 0,
			ADD_A_n, 0x80, RET, 0, 0, 0, 0, 0,
			LD_SP_nn, 0x50, 0, RST_00h, RST_08h, RST_10h, RST_18h, RST_20h,
			RST_28h, RST_30h, RST_38h, LD_B_n, 0x55, HALT, 0, 0,
		},
	}
	cpu := NewCPU(mem)
	cpu.PC = 0x40
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.r.A)
	assert.Equal(t, byte(0x55), cpu.r.B)
}

func Test_PUSH_rr(t *testing.T) {
	mem := &BasicMemory{
		cells: []byte{LD_SP_nn, 0x1B, 0x00, LD_A_n, 0x98,
			LD_BC_nn, 0x34, 0x12, LD_DE_nn, 0x35, 0x13, LD_HL_nn, 0x36, 0x14,
			PUSH_AF, PUSH_BC, PUSH_DE, PUSH_HL, HALT,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	cpu := NewCPU(mem)
	cpu.r.F = f_S | f_C
	cpu.Run()

	assert.Equal(t, cpu.mem.read(26), cpu.r.A)
	assert.Equal(t, cpu.mem.read(25), cpu.r.F)
	assert.Equal(t, cpu.mem.read(24), cpu.r.B)
	assert.Equal(t, cpu.mem.read(23), cpu.r.C)
	assert.Equal(t, cpu.mem.read(22), cpu.r.D)
	assert.Equal(t, cpu.mem.read(21), cpu.r.E)
	assert.Equal(t, cpu.mem.read(20), cpu.r.H)
	assert.Equal(t, cpu.mem.read(19), cpu.r.L)
}

func Test_POP_rr(t *testing.T) {
	mem := &BasicMemory{
		cells: []byte{LD_SP_nn, 0x08, 0x00, POP_AF, POP_BC, POP_DE, POP_HL, HALT, 0x43, 0x21, 0x44, 0x22, 0x45, 0x23, 0x46, 0x24},
	}

	cpu := NewCPU(mem)
	cpu.Run()
	assert.Equal(t, byte(0x21), cpu.r.A)
	assert.Equal(t, byte(0x43), cpu.r.F)
	assert.Equal(t, byte(0x22), cpu.r.B)
	assert.Equal(t, byte(0x44), cpu.r.C)
	assert.Equal(t, byte(0x23), cpu.r.D)
	assert.Equal(t, byte(0x45), cpu.r.E)
	assert.Equal(t, byte(0x24), cpu.r.H)
	assert.Equal(t, byte(0x46), cpu.r.L)
	assert.Equal(t, word(0x10), cpu.r.SP)
}

func Test_RLC_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_E_n, 0x55, __CB__, RLC_r | r_E, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.r.E)
	assert.Equal(t, f_S|f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_D_n, 0xAA, __CB__, RLC_r | r_D, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.r.D)
	assert.Equal(t, f_P|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x00, __CB__, RLC_r | r_A, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.r.A)
	assert.Equal(t, f_Z|f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_B_n, 0x80, __CB__, RLC_r | r_B, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.r.B)
	assert.Equal(t, f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_HL_nn, 0x06, 0x00, __CB__, RLC_r | r_HL, HALT, 0x01}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x02), cpu.mem.read(0x06))
	assert.Equal(t, f_NONE, cpu.r.F)
}

func Test_RRC_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_E_n, 0x55, __CB__, RRC_r | r_E, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_Z | f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.r.E)
	assert.Equal(t, f_S|f_P|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_D_n, 0xAA, __CB__, RRC_r | r_D, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_S | f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.r.D)
	assert.Equal(t, f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x00, __CB__, RRC_r | r_A, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.r.A)
	assert.Equal(t, f_Z|f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_B_n, 0x80, __CB__, RRC_r | r_B, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x40), cpu.r.B)
	assert.Equal(t, f_NONE, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_HL_nn, 0x06, 0x00, __CB__, RRC_r | r_HL, HALT, 0x01}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x80), cpu.mem.read(0x06))
	assert.Equal(t, f_S|f_C, cpu.r.F)
}

func Test_RL_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_E_n, 0x55, __CB__, RL_r | r_E, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.r.E)
	assert.Equal(t, f_S, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_D_n, 0xAA, __CB__, RL_r | r_D, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.r.D)
	assert.Equal(t, f_P|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x80, __CB__, RL_r | r_A, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.r.A)
	assert.Equal(t, f_Z|f_P|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_B_n, 0x80, __CB__, RL_r | r_B, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.r.B)
	assert.Equal(t, f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_HL_nn, 0x06, 0x00, __CB__, RL_r | r_HL, HALT, 0x81}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x02), cpu.mem.read(0x06))
	assert.Equal(t, f_C, cpu.r.F)
}

func Test_RR_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_E_n, 0x55, __CB__, RR_r | r_E, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.r.E)
	assert.Equal(t, f_S|f_P|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_D_n, 0xAA, __CB__, RR_r | r_D, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0xD5), cpu.r.D)
	assert.Equal(t, f_S, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_A_n, 0x01, __CB__, RR_r | r_A, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.r.A)
	assert.Equal(t, f_Z|f_P|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_B_n, 0x80, __CB__, RR_r | r_B, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0xC0), cpu.r.B)
	assert.Equal(t, f_S|f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_HL_nn, 0x06, 0x00, __CB__, RR_r | r_HL, HALT, 0x81}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x40), cpu.mem.read(0x06))
	assert.Equal(t, f_C, cpu.r.F)
}

func Test_SLA_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_E_n, 0x55, __CB__, SLA_r | r_E, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.r.E)
	assert.Equal(t, f_S|f_P, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_D_n, 0xAA, __CB__, SLA_r | r_D, HALT}}
	cpu = NewCPU(mem)
	cpu.r.F = f_Z | f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x54), cpu.r.D)
	assert.Equal(t, f_C, cpu.r.F)
}

func Test_SRA_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_E_n, 0x85, __CB__, SRA_r | r_E, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_Z | f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xC2), cpu.r.E)
	assert.Equal(t, f_S|f_C, cpu.r.F)
}

func Test_SLL_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_E_n, 0x95, __CB__, SLL_r | r_E, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_S | f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x2B), cpu.r.E)
	assert.Equal(t, f_P|f_C, cpu.r.F)
}

func Test_SRL_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_H_n, 0x85, __CB__, SRL_r | r_H, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_S | f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x42), cpu.r.H)
	assert.Equal(t, f_P|f_C, cpu.r.F)
}

func Test_BIT_b(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_E_n, 0x40, __CB__, BIT_b | r_E | BIT_6, HALT}}
	cpu := NewCPU(mem)
	cpu.r.F = f_Z | f_N | f_C
	cpu.Run()

	assert.Equal(t, f_H|f_C, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_L_n, 0xFE, __CB__, BIT_b | r_L | BIT_0, HALT}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, f_Z|f_H, cpu.r.F)

	mem = &BasicMemory{cells: []byte{LD_HL_nn, 0x06, 0x00, __CB__, BIT_b | r_HL | BIT_2, HALT, 0xFD}}
	cpu = NewCPU(mem)
	cpu.r.F = f_Z | f_N
	cpu.Run()

	assert.Equal(t, f_H, cpu.r.F)
}

func Test_RES_b(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_D_n, 0xFF, __CB__, RES_b | r_D | BIT_7, HALT}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x7F), cpu.r.D)

	mem = &BasicMemory{cells: []byte{LD_HL_nn, 0x06, 0x00, __CB__, RES_b | r_HL | BIT_2, HALT, 0xFF}}
	cpu = NewCPU(mem)
	cpu.r.F = f_Z | f_N
	cpu.Run()

	assert.Equal(t, byte(0xFB), cpu.mem.read(0x06))
}

func Test_SET_b(t *testing.T) {
	mem := &BasicMemory{cells: []byte{LD_D_n, 0x00, __CB__, SET_b | r_D | BIT_7, HALT}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x80), cpu.r.D)

	mem = &BasicMemory{cells: []byte{LD_HL_nn, 0x06, 0x00, __CB__, SET_b | r_HL | BIT_2, HALT, 0x00}}
	cpu = NewCPU(mem)
	cpu.r.F = f_Z | f_N
	cpu.Run()

	assert.Equal(t, byte(0x04), cpu.mem.read(0x06))
}

func Test_shouldJump(t *testing.T) {
	var tests = []struct {
		flags    byte
		code     byte
		expected bool
	}{
		{f_NONE, 0b00000000, true}, {f_Z, 0b00000000, false},
		{f_NONE, 0b00001000, false}, {f_Z, 0b00001000, true},
		{f_NONE, 0b00010000, true}, {f_C, 0b00010000, false},
		{f_NONE, 0b00011000, false}, {f_C, 0b00011000, true},
		{f_NONE, 0b00100000, true}, {f_P, 0b00100000, false},
		{f_NONE, 0b00101000, false}, {f_P, 0b00101000, true},
		{f_NONE, 0b00110000, true}, {f_S, 0b00110000, false},
		{f_NONE, 0b00111000, false}, {f_S, 0b00111000, true},
	}

	cpu := NewCPU(&BasicMemory{})
	for _, test := range tests {
		cpu.r.F = test.flags
		result := cpu.shouldJump(test.code)

		assert.Equal(t, test.expected, result)
	}
}

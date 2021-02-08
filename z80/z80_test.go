package z80

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NOP(t *testing.T) {
	mem := &BasicMemory{cells: []byte{nop, nop, nop, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_ALL
	cpu.Run()

	assert.Equal(t, f_ALL, cpu.reg.F)
}

func Test_DI(t *testing.T) {
	mem := &BasicMemory{cells: []byte{di, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, false, cpu.iff1)
	assert.Equal(t, false, cpu.iff2)
}

func Test_EI(t *testing.T) {
	mem := &BasicMemory{cells: []byte{di, ei, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, true, cpu.iff1)
	assert.Equal(t, true, cpu.iff2)
}

func Test_EX_AF_AF(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ex_af_af, halt}}
	cpu := NewCPU(mem)
	var a, a_, f, f_ byte = 0xcc, 0x55, 0x35, 0x97
	cpu.reg.A, cpu.reg.A_ = a, a_
	cpu.reg.F = f
	cpu.reg.F_ = f_
	cpu.Run()

	assert.Equal(t, a_, cpu.reg.A)
	assert.Equal(t, a, cpu.reg.A_)
	assert.Equal(t, f_, cpu.reg.F)
	assert.Equal(t, f, cpu.reg.F_)

	cpu.PC = 0
	cpu.Run()
	assert.Equal(t, a, cpu.reg.A)
	assert.Equal(t, a_, cpu.reg.A_)
	assert.Equal(t, f, cpu.reg.F)
	assert.Equal(t, f_, cpu.reg.F_)
}

func Test_EXX(t *testing.T) {
	mem := &BasicMemory{cells: []byte{exx, halt}}
	cpu := NewCPU(mem)
	cpu.reg.B, cpu.reg.C, cpu.reg.B_, cpu.reg.C_ = 0x01, 0x02, 0x03, 0x04
	cpu.reg.D, cpu.reg.E, cpu.reg.D_, cpu.reg.E_ = 0x05, 0x06, 0x07, 0x08
	cpu.reg.H, cpu.reg.L, cpu.reg.H_, cpu.reg.L_ = 0x09, 0x0A, 0x0B, 0x0C
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.reg.B_)
	assert.Equal(t, byte(0x02), cpu.reg.C_)
	assert.Equal(t, byte(0x03), cpu.reg.B)
	assert.Equal(t, byte(0x04), cpu.reg.C)
	assert.Equal(t, byte(0x05), cpu.reg.D_)
	assert.Equal(t, byte(0x06), cpu.reg.E_)
	assert.Equal(t, byte(0x07), cpu.reg.D)
	assert.Equal(t, byte(0x08), cpu.reg.E)
	assert.Equal(t, byte(0x09), cpu.reg.H_)
	assert.Equal(t, byte(0x0A), cpu.reg.L_)
	assert.Equal(t, byte(0x0B), cpu.reg.H)
	assert.Equal(t, byte(0x0C), cpu.reg.L)
}

func Test_EX_DE_HL(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ex_de_hl, halt}}
	cpu := NewCPU(mem)
	cpu.reg.D, cpu.reg.E, cpu.reg.H, cpu.reg.L = 0x01, 0x02, 0x03, 0x04
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.reg.H)
	assert.Equal(t, byte(0x02), cpu.reg.L)
	assert.Equal(t, byte(0x03), cpu.reg.D)
	assert.Equal(t, byte(0x04), cpu.reg.E)
}

func Test_EX_SP_HL(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_hl_nn, 0x12, 0x70, ld_sp_nn, 0x08, 0x00, ex_sp_hl, halt, 0x11, 0x22}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x22), cpu.reg.H)
	assert.Equal(t, byte(0x11), cpu.reg.L)
	assert.Equal(t, byte(0x12), cpu.mem.read(8))
	assert.Equal(t, byte(0x70), cpu.mem.read(9))
}

func Test_EX_SP_IXY(t *testing.T) {
	for _, prefix := range []byte{prefix_ix, prefix_iy} {
		mem := &BasicMemory{cells: []byte{prefix, ld_hl_nn, 0x12, 0x70, ld_sp_nn, 0x0A, 0x00, prefix, ex_sp_hl, halt, 0x11, 0x22}}
		cpu := NewCPU(mem)
		cpu.Run()

		switch prefix {
		case prefix_ix:
			assert.Equal(t, byte(0x22), byte(cpu.reg.IX>>8))
			assert.Equal(t, byte(0x11), byte(cpu.reg.IX))
		case prefix_iy:
			assert.Equal(t, byte(0x22), byte(cpu.reg.IY>>8))
			assert.Equal(t, byte(0x11), byte(cpu.reg.IY))
		}
		assert.Equal(t, byte(0x12), cpu.mem.read(10))
		assert.Equal(t, byte(0x70), cpu.mem.read(11))
	}
}

func Test_ADD_A_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0, add_a_n, 0, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_ALL & ^f_Z
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, f_Z, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0xFF, ld_h_n, 0x00, ld_l_n, 0x08, add_a_hl, halt, 0x01}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_ALL & ^f_Z
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, f_Z|f_H|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x70, ld_l_n, 0x70, add_a_l, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_NONE
	cpu.Run()

	assert.Equal(t, byte(0xE0), cpu.reg.A)
	assert.Equal(t, f_S|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0xF0, add_a_n, 0xB0, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_NONE
	cpu.Run()

	assert.Equal(t, byte(0xA0), cpu.reg.A)
	assert.Equal(t, f_S|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x8f, add_a_n, 0x81, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_NONE
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.reg.A)
	assert.Equal(t, f_H|f_P|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x22, prefix_ix, ld_hl_nn, 0x13, 0x00, prefix_ix, add_a_l, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_NONE
	cpu.Run()

	assert.Equal(t, byte(0x35), cpu.reg.A)
	assert.Equal(t, f_NONE, cpu.reg.F)
}

func Test_ADC_A_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0x20, adc_a_n, 0x20, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x41), cpu.reg.A)
	assert.Equal(t, f_NONE, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0, ld_b_n, 0xFF, adc_a_b, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, f_Z|f_H|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x0F, ld_b_n, 0x00, adc_a_b, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.reg.A)
	assert.Equal(t, f_H, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x0F, ld_l_n, 0x06, adc_a_hl, halt, 0x70}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x80), cpu.reg.A)
	assert.Equal(t, f_S|f_H|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x0F, ld_b_n, 0x69, adc_a_b, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x79), cpu.reg.A)
	assert.Equal(t, f_H, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x0E, ld_b_n, 0x01, adc_a_b, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.reg.A)
	assert.Equal(t, f_H, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x0E, ld_b_n, 0x01, adc_a_b, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x0F), cpu.reg.A)
	assert.Equal(t, f_NONE, cpu.reg.F)
}

func Test_ADD_HL_RR(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_hl_nn, 0xFF, 0xFF, ld_bc_nn, 0x01, 0, add_hl_bc, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.H)
	assert.Equal(t, byte(0), cpu.reg.L)
	assert.Equal(t, f_H|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_hl_nn, 0x41, 0x42, ld_de_nn, 0x11, 0x11, add_hl_de, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_ALL
	cpu.Run()

	assert.Equal(t, byte(0x53), cpu.reg.H)
	assert.Equal(t, byte(0x52), cpu.reg.L)
	assert.Equal(t, f_S|f_Z|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_hl_nn, 0x41, 0x42, add_hl_hl, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x84), cpu.reg.H)
	assert.Equal(t, byte(0x82), cpu.reg.L)
	assert.Equal(t, f_NONE, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_hl_nn, 0xFE, 0xFF, ld_sp_nn, 0x02, 0, add_hl_sp, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.H)
	assert.Equal(t, byte(0), cpu.reg.L)
	assert.Equal(t, f_H|f_C, cpu.reg.F)
}

func Test_SUB_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0, sub_n, 0x01, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_Z
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)
	assert.Equal(t, f_S|f_H|f_N|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x20, sub_a, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, f_Z|f_N, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x90, ld_h_n, 0x20, sub_h, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_Z
	cpu.Run()

	assert.Equal(t, byte(0x70), cpu.reg.A)
	assert.Equal(t, f_P|f_N, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x7F, ld_l_n, 0x06, sub_hl, halt, 0x80}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)
	assert.Equal(t, f_S|f_P|f_N|f_C, cpu.reg.F)
}

func Test_CP_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0, ld_b_n, 0x01, cp_b, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_Z
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, f_S|f_H|f_N|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x20, cp_a, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x20), cpu.reg.A)
	assert.Equal(t, f_Z|f_N, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x90, ld_h_n, 0x20, cp_h, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_Z
	cpu.Run()

	assert.Equal(t, byte(0x90), cpu.reg.A)
	assert.Equal(t, f_P|f_N, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x7F, ld_l_n, 0x06, cp_hl, halt, 0x80}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x7F), cpu.reg.A)
	assert.Equal(t, f_S|f_P|f_N|f_C, cpu.reg.F)
}

func Test_SBC_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0x01, ld_b_n, 0x01, sbc_a_b, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)
	assert.Equal(t, f_S|f_H|f_N|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x7F, ld_l_n, 0x80, sbc_a_l, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0xFE), cpu.reg.A)
	assert.Equal(t, f_S|f_P|f_N|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x02, ld_d_n, 0x01, sbc_a_d, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, f_Z|f_N, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x81, ld_l_n, 0x06, sbc_a_hl, halt, 0x01}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x7F), cpu.reg.A)
	assert.Equal(t, f_H|f_P|f_N, cpu.reg.F)
}

func Test_AND_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0x0F, ld_b_n, 0xF0, and_b, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, f_Z|f_H|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x8F, ld_d_n, 0xF3, and_d, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x83), cpu.reg.A)
	assert.Equal(t, f_S|f_H, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0xFF, ld_l_n, 0x06, and_hl, halt, 0x81}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x81), cpu.reg.A)
	assert.Equal(t, f_S|f_H|f_P, cpu.reg.F)
}

func Test_OR_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0x00, ld_b_n, 0x00, or_b, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_S | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, f_Z|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x8A, ld_l_n, 0x06, or_hl, halt, 0x85}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_S | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x8F), cpu.reg.A)
	assert.Equal(t, f_S, cpu.reg.F)
}

func Test_XOR_x(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0x1F, ld_b_n, 0x1F, xor_b, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, f_Z|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x1F, ld_l_n, 0x06, xor_hl, halt, 0x8F}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x90), cpu.reg.A)
	assert.Equal(t, f_S|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x1F, xor_n, 0x0F, halt, 0x8F}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.reg.A)
	assert.Equal(t, f_NONE, cpu.reg.F)
}

func Test_INC_R(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0, inc_a, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_ALL
	cpu.Run()

	assert.Equal(t, f_C, cpu.reg.F)

	cpu.Reset()
	cpu.reg.F = f_ALL & ^f_Z
	mem.cells[1] = 0xFF
	cpu.Run()

	assert.Equal(t, f_Z|f_H|f_C, cpu.reg.F)

	cpu.Reset()
	cpu.reg.F = f_N
	mem.cells[1] = 0x7F
	cpu.Run()

	assert.Equal(t, f_S|f_H|f_P, cpu.reg.F)

	cpu.Reset()
	mem.cells[1] = 0x92
	cpu.Run()

	assert.Equal(t, f_S, cpu.reg.F)

	cpu.Reset()
	mem.cells[1] = 0x10
	cpu.Run()

	assert.Equal(t, f_NONE, cpu.reg.F)
}

func Test_INC_RR(t *testing.T) {
	mem := &BasicMemory{
		cells: []byte{
			ld_bc_nn, 0x34, 0x12, inc_bc, ld_de_nn, 0x35, 0x13, inc_de,
			ld_hl_nn, 0x36, 0x14, inc_hl, ld_sp_nn, 0x37, 0x15, inc_sp,
			halt},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, word(0x1235), cpu.reg.getBC())
	assert.Equal(t, word(0x1336), cpu.reg.getDE())
	assert.Equal(t, word(0x1437), cpu.reg.getHL())
	assert.Equal(t, word(0x1538), cpu.reg.SP)
}

func Test_INC_mHL(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_hl_nn, 0x05, 0x00, inc_mhl, halt, 0xFF}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_S | f_P | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x00), mem.cells[5])
	assert.Equal(t, f_Z|f_H|f_C, cpu.reg.F)

	cpu.Reset()
	mem.cells[5] = 0x7F
	cpu.Run()

	assert.Equal(t, byte(0x80), mem.cells[5])
	assert.Equal(t, f_S|f_H|f_P, cpu.reg.F)

	cpu.Reset()
	mem.cells[5] = 0x20
	cpu.Run()

	assert.Equal(t, byte(0x21), mem.cells[5])
	assert.Equal(t, f_NONE, cpu.reg.F)
}

func Test_DEC_R(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 1, dec_a, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_NONE
	cpu.Run()

	assert.Equal(t, f_Z|f_N, cpu.reg.F)

	cpu.Reset()
	cpu.reg.F = f_ALL & ^(f_Z | f_H | f_N)
	mem.cells[1] = 0
	cpu.Run()

	cpu.Reset()
	cpu.reg.F = f_Z | f_S
	mem.cells[1] = 0x80
	cpu.Run()

	assert.Equal(t, f_H|f_P|f_N, cpu.reg.F)

	cpu.Reset()
	cpu.reg.F = f_ALL
	mem.cells[1] = 0xab
	cpu.Run()

	assert.Equal(t, f_S|f_N|f_C, cpu.reg.F)
}

func Test_DEC_RR(t *testing.T) {
	mem := &BasicMemory{
		cells: []byte{
			ld_bc_nn, 0x34, 0x12, dec_bc, ld_de_nn, 0x35, 0x13, dec_de,
			ld_hl_nn, 0x36, 0x14, dec_hl, ld_sp_nn, 0x37, 0x15, dec_sp,
			halt},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, word(0x1233), cpu.reg.getBC())
	assert.Equal(t, word(0x1334), cpu.reg.getDE())
	assert.Equal(t, word(0x1435), cpu.reg.getHL())
	assert.Equal(t, word(0x1536), cpu.reg.SP)
}

func Test_DEC_mHL(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_hl_nn, 0x05, 0x00, dec_mhl, halt, 0x00}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0xFF), mem.cells[5])
	assert.Equal(t, f_S|f_H|f_N|f_C, cpu.reg.F)

	cpu.Reset()
	mem.cells[5] = 0x01
	cpu.Run()

	assert.Equal(t, byte(0x00), mem.cells[5])
	assert.Equal(t, f_Z|f_N, cpu.reg.F)

	cpu.Reset()
	mem.cells[5] = 0x80
	cpu.Run()

	assert.Equal(t, byte(0x7F), mem.cells[5])
	assert.Equal(t, f_P|f_H|f_N, cpu.reg.F)
}

func Test_LD_RR_nn(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_bc_nn, 0x34, 0x12, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, word(0x1234), cpu.reg.getBC())

	mem = &BasicMemory{cells: []byte{ld_sp_nn, 0x34, 0x12, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, word(0x1234), cpu.reg.SP)
}

func Test_LD_mm_HL(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_hl_nn, 0x3A, 0x48, ld_mm_hl, 0x07, 0x00, halt, 0, 0}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, cpu.reg.H, mem.cells[8])
	assert.Equal(t, cpu.reg.L, mem.cells[7])
}

func Test_LD_HL_mm(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_hl_mm, 0x04, 0x00, halt, 0x34, 0x12}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x12), cpu.reg.H)
	assert.Equal(t, byte(0x34), cpu.reg.L)
}

func Test_LD_mHL_n(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_hl_nn, 0x06, 0x00, ld_mhl_n, 0xAB, halt, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.mem.read(6))
}

func Test_LD_mIXY_n(t *testing.T) {
	for _, prefix := range []byte{prefix_ix, prefix_iy} {
		mem := &BasicMemory{cells: []byte{prefix, ld_hl_nn, 0x02, 0x00, prefix, ld_mhl_n, 0x07, 0xAB, halt, 0x00}}
		cpu := NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0xAB), cpu.mem.read(9))
	}
}

func Test_LD_mm_A(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0x9F, ld_mm_a, 0x06, 0x00, halt, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, cpu.reg.A, mem.cells[6])
}

func Test_LD_A_mm(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_mm, 0x04, 0x00, halt, 0xDE}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xDE), cpu.reg.A)
}

func Test_LD_BC_A(t *testing.T) {
	var n byte = 0x76
	mem := &BasicMemory{cells: []byte{ld_a_n, n, ld_bc_nn, 0x07, 0x00, ld_bc_a, halt, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, n, cpu.mem.read(7))
}

func Test_LD_DE_A(t *testing.T) {
	var n byte = 0x76
	mem := &BasicMemory{cells: []byte{ld_a_n, n, ld_de_nn, 0x07, 0x00, ld_de_a, halt, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, n, cpu.mem.read(7))
}

func Test_LD_A_BC(t *testing.T) {
	var n byte = 0x76
	mem := &BasicMemory{cells: []byte{ld_bc_nn, 0x05, 0x00, ld_a_bc, halt, n}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, n, cpu.reg.A)
}

func Test_LD_A_DE(t *testing.T) {
	var n byte = 0x76
	mem := &BasicMemory{cells: []byte{ld_de_nn, 0x05, 0x00, ld_a_de, halt, n}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, n, cpu.reg.A)
}

func Test_LD_R_n(t *testing.T) {
	var a, b, c, d, e, h, l byte = 1, 2, 3, 4, 5, 6, 7
	mem := &BasicMemory{
		cells: []byte{ld_a_n, a, ld_b_n, b, ld_c_n, c, ld_d_n, d, ld_e_n, e, ld_h_n, h, ld_l_n, l, halt},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, a, cpu.reg.A)
	assert.Equal(t, b, cpu.reg.B)
	assert.Equal(t, c, cpu.reg.C)
	assert.Equal(t, d, cpu.reg.D)
	assert.Equal(t, e, cpu.reg.E)
	assert.Equal(t, h, cpu.reg.H)
	assert.Equal(t, l, cpu.reg.L)
}

func Test_LD_R_R(t *testing.T) {
	mem := &BasicMemory{
		cells: []byte{ld_a_n, 0x56, ld_b_a, ld_c_b, ld_d_c, ld_e_d, ld_h_e, ld_l_h, ld_a_n, 0, ld_a_b, halt},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x56), cpu.reg.A)
	assert.Equal(t, byte(0x56), cpu.reg.B)
	assert.Equal(t, byte(0x56), cpu.reg.C)
	assert.Equal(t, byte(0x56), cpu.reg.D)
	assert.Equal(t, byte(0x56), cpu.reg.E)
	assert.Equal(t, byte(0x56), cpu.reg.H)
	assert.Equal(t, byte(0x56), cpu.reg.L)
}

func Test_LD_R_HL(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_hl_nn, 0x06, 0x00, ld_a_hl, ld_l_hl, halt, 0xA7}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xA7), cpu.reg.A)
	assert.Equal(t, byte(0xA7), cpu.reg.L)
}

func Test_LD_R_IXY(t *testing.T) {
	for _, prefix := range []byte{prefix_ix, prefix_iy} {
		mem := &BasicMemory{cells: []byte{prefix, ld_hl_nn, 0x01, 0x00, prefix, ld_l_hl, 0x07, halt, 0xA7}}
		cpu := NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0xA7), cpu.reg.L)
	}
}

func Test_LD_HL_R(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_d_n, 0x99, ld_hl_nn, 0x07, 0x00, ld_hl_d, halt, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x99), cpu.mem.read(7))
}

func Test_CPL(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0x5B, cpl, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xA4), cpu.reg.A)
	assert.Equal(t, f_H|f_N, cpu.reg.F)
}

func Test_SCF(t *testing.T) {
	mem := &BasicMemory{cells: []byte{scf, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_S | f_Z | f_H | f_P | f_N
	cpu.Run()

	assert.Equal(t, f_S|f_Z|f_P|f_C, cpu.reg.F)
}

func Test_CCF(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ccf, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_ALL
	cpu.Run()

	assert.Equal(t, f_S|f_Z|f_P, cpu.reg.F)

	cpu.Reset()
	cpu.reg.F = f_Z | f_N | f_C
	cpu.Run()

	assert.Equal(t, f_Z|f_H, cpu.reg.F)
}

func Test_RLCA(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0x55, rlca, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.reg.A)
	assert.Equal(t, f_NONE, cpu.reg.F)

	cpu.Reset()
	mem.cells[1] = 0xAA
	cpu.reg.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.reg.A)
	assert.Equal(t, f_C, cpu.reg.F)

	cpu.Reset()
	mem.cells[1] = 0x00
	cpu.reg.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, f_NONE, cpu.reg.F)

	cpu.Reset()
	mem.cells[1] = 0xFF
	cpu.reg.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)
	assert.Equal(t, f_C, cpu.reg.F)
}

func Test_RRCA(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0x55, rrca, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.reg.A)
	assert.Equal(t, f_C, cpu.reg.F)

	cpu.Reset()
	mem.cells[1] = 0xAA
	cpu.reg.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.reg.A)
	assert.Equal(t, f_NONE, cpu.reg.F)

	cpu.Reset()
	mem.cells[1] = 0x00
	cpu.reg.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, f_NONE, cpu.reg.F)

	cpu.Reset()
	mem.cells[1] = 0xFF
	cpu.reg.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)
	assert.Equal(t, f_C, cpu.reg.F)
}

func Test_RLA(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0x80, rla, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.reg.A)
	assert.Equal(t, f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x55, rla, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_ALL
	cpu.Run()

	assert.Equal(t, byte(0xab), cpu.reg.A)
	assert.Equal(t, f_S|f_Z|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x88, rla, ld_b_a, rla, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.reg.B)
	assert.Equal(t, byte(0x21), cpu.reg.A)
	assert.Equal(t, f_NONE, cpu.reg.F)
}

func Test_RRA(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0x80, rra, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0xC0), cpu.reg.A)
	assert.Equal(t, f_NONE, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x55, rra, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_ALL
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.reg.A)
	assert.Equal(t, f_S|f_Z|f_P|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x89, rra, ld_b_a, rra, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x44), cpu.reg.B)
	assert.Equal(t, byte(0xA2), cpu.reg.A)
	assert.Equal(t, f_NONE, cpu.reg.F)
}

func Test_DAA(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 0x9A, daa, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, f_Z|f_H|f_P|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x99, daa, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x99), cpu.reg.A)
	assert.Equal(t, f_S|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x8F, daa, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x95), cpu.reg.A)
	assert.Equal(t, f_S|f_H|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x8F, daa, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_N
	cpu.Run()

	assert.Equal(t, byte(0x89), cpu.reg.A)
	assert.Equal(t, f_S|f_N, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0xCA, daa, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_N
	cpu.Run()

	assert.Equal(t, byte(0x64), cpu.reg.A)
	assert.Equal(t, f_N|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0xC5, daa, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x5F), cpu.reg.A)
	assert.Equal(t, f_H|f_P|f_N|f_C, cpu.reg.F)
}

func Test_DJNZ(t *testing.T) {
	var b byte = 0x20
	var o int8 = -3
	mem := &BasicMemory{cells: []byte{ld_b_n, b, ld_a_n, 0, inc_a, djnz, byte(o), halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, b, cpu.reg.A)
	assert.Equal(t, byte(0), cpu.reg.B)

	b = 0xFF
	o = 1
	mem = &BasicMemory{cells: []byte{ld_b_n, b, ld_a_n, 1, djnz, byte(o), halt, inc_a, jr_o, 0xFA}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, b, cpu.reg.A)
	assert.Equal(t, byte(0), cpu.reg.B)
}

func Test_JR_o(t *testing.T) {
	mem := &BasicMemory{cells: []byte{jr_o, 3, ld_c_n, 0x11, halt, ld_d_n, 0x22, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.C)
	assert.Equal(t, byte(0x22), cpu.reg.D)

	mem = &BasicMemory{cells: []byte{jr_o, 6, halt, ld_c_n, 0x11, ld_b_n, 0x33, halt, jr_o, 0xF9, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x33), cpu.reg.B)
	assert.Equal(t, byte(0x11), cpu.reg.C)
}

func Test_JR_Z_o(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 2, dec_a, jr_z_o, 0x02, ld_b_n, 0xab, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xab), cpu.reg.B)

	mem = &BasicMemory{cells: []byte{ld_a_n, 1, dec_a, jr_z_o, 0x02, ld_b_n, 0xab, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	mem = &BasicMemory{cells: []byte{ld_a_n, 1, dec_a, jr_z_o, 0xFD, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)
}

func Test_JR_NZ_o(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_a_n, 2, dec_a, jr_nz_o, 0x02, ld_b_n, 0xab, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.B)

	mem = &BasicMemory{cells: []byte{ld_a_n, 1, dec_a, jr_nz_o, 0x02, ld_b_n, 0xab, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xab), cpu.reg.B)

	mem = &BasicMemory{cells: []byte{ld_a_n, 2, dec_a, jr_nz_o, 0xFD, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
}

func Test_JR_C(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_b_n, 0xAB, dec_a, add_a_n, 1, jr_c, 1, ld_b_a, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.reg.B)

	mem = &BasicMemory{cells: []byte{ld_b_n, 0xAB, inc_a, add_a_n, 1, jr_c, 1, ld_b_a, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(2), cpu.reg.B)

	mem = &BasicMemory{cells: []byte{dec_a, add_a_n, 1, jr_c, 0xFC, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(1), cpu.reg.A)
}

func Test_JR_NC_o(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_b_n, 0xAB, inc_a, add_a_n, 1, jr_nc_o, 1, ld_b_a, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.reg.B)

	mem = &BasicMemory{cells: []byte{ld_b_n, 0xAB, dec_a, add_a_n, 1, jr_nc_o, 1, ld_b_a, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.B)

	mem = &BasicMemory{cells: []byte{add_a_n, 1, jr_nc_o, 0xFC, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
}

func Test_JP_nn(t *testing.T) {
	mem := &BasicMemory{cells: []byte{jp_nn, 0x06, 0x00, ld_a_n, 0xAA, halt, ld_a_n, 0x55, halt, 0x07, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.reg.A)
}

func Test_JP_cc_nn(t *testing.T) {
	var tests = []struct {
		jp       byte
		flag     byte
		expected byte
	}{
		{jp_c_nn, f_C, 0x55}, {jp_nc_nn, f_NONE, 0x55}, {jp_z_nn, f_Z, 0x55}, {jp_nz_nn, f_NONE, 0x55},
		{jp_m_nn, f_S, 0x55}, {jp_p_nn, f_NONE, 0x55}, {jp_pe_nn, f_P, 0x55}, {jp_po_nn, f_NONE, 0x55},
		{jp_c_nn, f_NONE, 0xAA}, {jp_nc_nn, f_C, 0xAA}, {jp_z_nn, f_NONE, 0xAA}, {jp_nz_nn, f_Z, 0xAA},
		{jp_m_nn, f_NONE, 0xAA}, {jp_p_nn, f_S, 0xAA}, {jp_pe_nn, f_NONE, 0xAA}, {jp_po_nn, f_P, 0xAA},
	}

	for _, test := range tests {
		mem := &BasicMemory{
			cells: []byte{test.jp, 0x06, 0x00, ld_a_n, 0xAA, halt, ld_a_n, 0x55, halt, 0x07, 0x00},
		}

		cpu := NewCPU(mem)
		cpu.reg.F = test.flag
		cpu.Run()

		assert.Equal(t, byte(test.expected), cpu.reg.A)
	}
}

func Test_CALL_cc_nn(t *testing.T) {
	var tests = []struct {
		call     byte
		flag     byte
		expected byte
	}{
		{call_c_nn, f_C, 0x55}, {call_nc_nn, f_NONE, 0x55}, {call_z_nn, f_Z, 0x55}, {call_nz_nn, f_NONE, 0x55},
		{call_m_nn, f_S, 0x55}, {call_p_nn, f_NONE, 0x55}, {call_pe_nn, f_P, 0x55}, {call_po_nn, f_NONE, 0x55},
		{call_c_nn, f_NONE, 0xAA}, {call_nc_nn, f_C, 0xAA}, {call_z_nn, f_NONE, 0xAA}, {call_nz_nn, f_Z, 0xAA},
		{call_m_nn, f_NONE, 0xAA}, {call_p_nn, f_S, 0xAA}, {call_pe_nn, f_NONE, 0xAA}, {call_po_nn, f_P, 0xAA},
		{call_nn, f_NONE, 0x55},
	}

	for _, test := range tests {
		mem := &BasicMemory{
			cells: []byte{ld_sp_nn, 0x10, 0x00, test.call, 0x09, 0x00, ld_a_n, 0xAA, halt, ld_a_n, 0x55, halt, 0xFF, 0xFF, 0xFF, 0xFF},
		}

		cpu := NewCPU(mem)
		cpu.reg.F = test.flag
		cpu.Run()

		assert.Equal(t, byte(test.expected), cpu.reg.A)
		if cpu.reg.A == 0x55 {
			assert.Equal(t, word(0x0E), cpu.reg.SP)
			assert.Equal(t, byte(0), mem.cells[15])
			assert.Equal(t, byte(0x06), mem.cells[14])
		} else {
			assert.Equal(t, word(0x10), cpu.reg.SP)
		}
	}
}

func Test_RET(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_sp_nn, 0x0A, 0x00, ret, ld_a_n, 0xAA, halt, ld_a_n, 0x55, halt, 0x07, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.reg.A)
}

func Test_RET_cc(t *testing.T) {
	var tests = []struct {
		ret      byte
		flag     byte
		expected byte
	}{
		{ret_c, f_C, 0x55}, {ret_nc, f_NONE, 0x55}, {ret_z, f_Z, 0x55}, {ret_nz, f_NONE, 0x55},
		{ret_m, f_S, 0x55}, {ret_p, f_NONE, 0x55}, {ret_pe, f_P, 0x55}, {ret_po, f_NONE, 0x55},
		{ret_c, f_NONE, 0xAA}, {ret_nc, f_C, 0xAA}, {ret_z, f_NONE, 0xAA}, {ret_nz, f_Z, 0xAA},
		{ret_m, f_NONE, 0xAA}, {ret_p, f_S, 0xAA}, {ret_pe, f_NONE, 0xAA}, {ret_po, f_P, 0xAA},
	}

	for _, test := range tests {
		mem := &BasicMemory{
			cells: []byte{ld_sp_nn, 0x0A, 0x00, test.ret, ld_a_n, 0xAA, halt, ld_a_n, 0x55, halt, 0x07, 0x00},
		}
		cpu := NewCPU(mem)
		cpu.reg.F = test.flag
		cpu.Run()

		assert.Equal(t, byte(test.expected), cpu.reg.A)
	}
}

func Test_RST_xx(t *testing.T) {
	mem := &BasicMemory{
		cells: []byte{
			ld_a_n, 0x01, ret, 0, 0, 0, 0, 0,
			add_a_n, 0x02, ret, 0, 0, 0, 0, 0,
			add_a_n, 0x04, ret, 0, 0, 0, 0, 0,
			add_a_n, 0x08, ret, 0, 0, 0, 0, 0,
			add_a_n, 0x10, ret, 0, 0, 0, 0, 0,
			add_a_n, 0x20, ret, 0, 0, 0, 0, 0,
			add_a_n, 0x40, ret, 0, 0, 0, 0, 0,
			add_a_n, 0x80, ret, 0, 0, 0, 0, 0,
			ld_sp_nn, 0x50, 0, rst_00h, rst_08h, rst_10h, rst_18h, rst_20h,
			rst_28h, rst_30h, rst_38h, ld_b_n, 0x55, halt, 0, 0,
		},
	}
	cpu := NewCPU(mem)
	cpu.PC = 0x40
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)
	assert.Equal(t, byte(0x55), cpu.reg.B)
}

func Test_PUSH_rr(t *testing.T) {
	mem := &BasicMemory{
		cells: []byte{ld_sp_nn, 0x1B, 0x00, ld_a_n, 0x98,
			ld_bc_nn, 0x34, 0x12, ld_de_nn, 0x35, 0x13, ld_hl_nn, 0x36, 0x14,
			push_af, push_bc, push_de, push_hl, halt,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	cpu := NewCPU(mem)
	cpu.reg.F = f_S | f_C
	cpu.Run()

	assert.Equal(t, cpu.mem.read(26), cpu.reg.A)
	assert.Equal(t, cpu.mem.read(25), cpu.reg.F)
	assert.Equal(t, cpu.mem.read(24), cpu.reg.B)
	assert.Equal(t, cpu.mem.read(23), cpu.reg.C)
	assert.Equal(t, cpu.mem.read(22), cpu.reg.D)
	assert.Equal(t, cpu.mem.read(21), cpu.reg.E)
	assert.Equal(t, cpu.mem.read(20), cpu.reg.H)
	assert.Equal(t, cpu.mem.read(19), cpu.reg.L)
}

func Test_POP_rr(t *testing.T) {
	mem := &BasicMemory{
		cells: []byte{ld_sp_nn, 0x08, 0x00, pop_af, pop_bc, pop_de, pop_hl, halt, 0x43, 0x21, 0x44, 0x22, 0x45, 0x23, 0x46, 0x24},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x21), cpu.reg.A)
	assert.Equal(t, byte(0x43), cpu.reg.F)
	assert.Equal(t, byte(0x22), cpu.reg.B)
	assert.Equal(t, byte(0x44), cpu.reg.C)
	assert.Equal(t, byte(0x23), cpu.reg.D)
	assert.Equal(t, byte(0x45), cpu.reg.E)
	assert.Equal(t, byte(0x24), cpu.reg.H)
	assert.Equal(t, byte(0x46), cpu.reg.L)
	assert.Equal(t, word(0x10), cpu.reg.SP)
}

func Test_IN_A_n(t *testing.T) {
	mem := &BasicMemory{
		cells: []byte{ld_a_n, 0x23, in_a_n, 0x01, halt},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x23), cpu.reg.A)

	cpu.Reset()
	cpu.IN = func(a, n byte) byte {
		if a == 0x23 && n == 0x01 {
			return 0xA5
		}
		return 0
	}
	cpu.Run()
	assert.Equal(t, byte(0xA5), cpu.reg.A)
}

func Test_RLC_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_e_n, 0x55, prefix_bit, rlc_r | r_E, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.reg.E)
	assert.Equal(t, f_S|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_d_n, 0xAA, prefix_bit, rlc_r | r_D, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.reg.D)
	assert.Equal(t, f_P|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x00, prefix_bit, rlc_r | r_A, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, f_Z|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_b_n, 0x80, prefix_bit, rlc_r | r_B, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.reg.B)
	assert.Equal(t, f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_bit, rlc_r | r_HL, halt, 0x01}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x02), cpu.mem.read(0x06))
	assert.Equal(t, f_NONE, cpu.reg.F)
}

func Test_RRC_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_e_n, 0x55, prefix_bit, rrc_r | r_E, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_Z | f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.reg.E)
	assert.Equal(t, f_S|f_P|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_d_n, 0xAA, prefix_bit, rrc_r | r_D, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_S | f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.reg.D)
	assert.Equal(t, f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x00, prefix_bit, rrc_r | r_A, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, f_Z|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_b_n, 0x80, prefix_bit, rrc_r | r_B, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x40), cpu.reg.B)
	assert.Equal(t, f_NONE, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_bit, rrc_r | r_HL, halt, 0x01}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x80), cpu.mem.read(0x06))
	assert.Equal(t, f_S|f_C, cpu.reg.F)
}

func Test_RL_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_e_n, 0x55, prefix_bit, rl_r | r_E, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.reg.E)
	assert.Equal(t, f_S, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_d_n, 0xAA, prefix_bit, rl_r | r_D, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.reg.D)
	assert.Equal(t, f_P|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x80, prefix_bit, rl_r | r_A, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, f_Z|f_P|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_b_n, 0x80, prefix_bit, rl_r | r_B, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.reg.B)
	assert.Equal(t, f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_bit, rl_r | r_HL, halt, 0x81}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x02), cpu.mem.read(0x06))
	assert.Equal(t, f_C, cpu.reg.F)
}

func Test_RR_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_e_n, 0x55, prefix_bit, rr_r | r_E, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.reg.E)
	assert.Equal(t, f_S|f_P|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_d_n, 0xAA, prefix_bit, rr_r | r_D, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0xD5), cpu.reg.D)
	assert.Equal(t, f_S, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_a_n, 0x01, prefix_bit, rr_r | r_A, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, f_Z|f_P|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_b_n, 0x80, prefix_bit, rr_r | r_B, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_C
	cpu.Run()

	assert.Equal(t, byte(0xC0), cpu.reg.B)
	assert.Equal(t, f_S|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_bit, rr_r | r_HL, halt, 0x81}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x40), cpu.mem.read(0x06))
	assert.Equal(t, f_C, cpu.reg.F)
}

func Test_SLA_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_e_n, 0x55, prefix_bit, sla_r | r_E, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.reg.E)
	assert.Equal(t, f_S|f_P, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_d_n, 0xAA, prefix_bit, sla_r | r_D, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_Z | f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x54), cpu.reg.D)
	assert.Equal(t, f_C, cpu.reg.F)
}

func Test_SRA_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_e_n, 0x85, prefix_bit, sra_r | r_E, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_Z | f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0xC2), cpu.reg.E)
	assert.Equal(t, f_S|f_C, cpu.reg.F)
}

func Test_SLL_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_e_n, 0x95, prefix_bit, sll_r | r_E, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_S | f_Z | f_H | f_N | f_C
	cpu.Run()

	assert.Equal(t, byte(0x2B), cpu.reg.E)
	assert.Equal(t, f_P|f_C, cpu.reg.F)
}

func Test_SRL_r(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_h_n, 0x85, prefix_bit, srl_r | r_H, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_S | f_H | f_N
	cpu.Run()

	assert.Equal(t, byte(0x42), cpu.reg.H)
	assert.Equal(t, f_P|f_C, cpu.reg.F)
}

func Test_BIT_b(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_e_n, 0x40, prefix_bit, bit_b | r_E | bit_6, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = f_Z | f_N | f_C
	cpu.Run()

	assert.Equal(t, f_H|f_C, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_l_n, 0xFE, prefix_bit, bit_b | r_L | bit_0, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, f_Z|f_H, cpu.reg.F)

	mem = &BasicMemory{cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_bit, bit_b | r_HL | bit_2, halt, 0xFD}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_Z | f_N
	cpu.Run()

	assert.Equal(t, f_H, cpu.reg.F)
}

func Test_RES_b(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_d_n, 0xFF, prefix_bit, res_b | r_D | bit_7, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x7F), cpu.reg.D)

	mem = &BasicMemory{cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_bit, res_b | r_HL | bit_2, halt, 0xFF}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_Z | f_N
	cpu.Run()

	assert.Equal(t, byte(0xFB), cpu.mem.read(0x06))
}

func Test_SET_b(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_d_n, 0x00, prefix_bit, set_b | r_D | bit_7, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x80), cpu.reg.D)

	mem = &BasicMemory{cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_bit, set_b | r_HL | bit_2, halt, 0x00}}
	cpu = NewCPU(mem)
	cpu.reg.F = f_Z | f_N
	cpu.Run()

	assert.Equal(t, byte(0x04), cpu.mem.read(0x06))
}

func Test_LD_IXY_nn(t *testing.T) {
	mem := &BasicMemory{cells: []byte{
		prefix_ix, ld_hl_nn, 0x06, 0x01, prefix_iy, ld_hl_nn, 0x07, 0x02,
		ld_hl_nn, 0x08, 0x03, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, word(0x0106), cpu.reg.IX)
	assert.Equal(t, word(0x0207), cpu.reg.IY)
	assert.Equal(t, byte(0x03), cpu.reg.H)
	assert.Equal(t, byte(0x08), cpu.reg.L)
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
		cpu.reg.F = test.flags
		result := cpu.shouldJump(test.code)

		assert.Equal(t, test.expected, result)
	}
}

func Test_getHL(t *testing.T) {
	mem := &BasicMemory{cells: []byte{ld_hl_nn, 0x01, 0x02, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	result := cpu.getHL(prefix_none)
	assert.Equal(t, word(0x0201), result)

	for _, prefix := range []byte{prefix_ix, prefix_iy} {
		mem = &BasicMemory{cells: []byte{prefix, ld_a_hl, 0x05, halt, 0x00, 0xE5}}
		cpu = NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0xE5), cpu.reg.A)

		mem = &BasicMemory{cells: []byte{prefix, ld_hl_nn, 0x40, 0x00, prefix, ld_a_hl, 0xC9, halt, 0x00, 0xE5}}
		cpu = NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0xE5), cpu.reg.A)
	}
}

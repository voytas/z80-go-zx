package z80

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voytas/z80-go-zx/z80/memory"
)

func Test_NOP(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{nop, nop, nop, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, fALL, cpu.reg.F)
}

func Test_DI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{di, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, false, cpu.iff1)
	assert.Equal(t, false, cpu.iff2)
}

func Test_EI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{di, ei, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, true, cpu.iff1)
	assert.Equal(t, true, cpu.iff2)
}

func Test_IM_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{prefix_ed, im1, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, im1, cpu.im)
}

func Test_EX_AF_AF(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ex_af_af, halt}}
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

	cpu.reg.PC = 0
	cpu.Run()
	assert.Equal(t, a, cpu.reg.A)
	assert.Equal(t, a_, cpu.reg.A_)
	assert.Equal(t, f, cpu.reg.F)
	assert.Equal(t, f_, cpu.reg.F_)
}

func Test_EXX(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{exx, halt}}
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
	mem := &memory.BasicMemory{Cells: []byte{ex_de_hl, halt}}
	cpu := NewCPU(mem)
	cpu.reg.D, cpu.reg.E, cpu.reg.H, cpu.reg.L = 0x01, 0x02, 0x03, 0x04
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.reg.H)
	assert.Equal(t, byte(0x02), cpu.reg.L)
	assert.Equal(t, byte(0x03), cpu.reg.D)
	assert.Equal(t, byte(0x04), cpu.reg.E)
}

func Test_EX_SP_HL(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x12, 0x70, ld_sp_nn, 0x08, 0x00, ex_sp_hl, halt, 0x11, 0x22}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x22), cpu.reg.H)
	assert.Equal(t, byte(0x11), cpu.reg.L)
	assert.Equal(t, byte(0x12), cpu.mem.Read(8))
	assert.Equal(t, byte(0x70), cpu.mem.Read(9))

	for _, prefix := range []byte{useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x12, 0x70, ld_sp_nn, 0x0A, 0x00, prefix, ex_sp_hl, halt, 0x11, 0x22}}
		cpu := NewCPU(mem)
		cpu.Run()

		switch prefix {
		case useIX:
			assert.Equal(t, byte(0x22), byte(cpu.reg.IXH))
			assert.Equal(t, byte(0x11), byte(cpu.reg.IXL))
		case useIY:
			assert.Equal(t, byte(0x22), byte(cpu.reg.IYH))
			assert.Equal(t, byte(0x11), byte(cpu.reg.IYL))
		}
		assert.Equal(t, byte(0x12), cpu.mem.Read(10))
		assert.Equal(t, byte(0x70), cpu.mem.Read(11))
	}
}

func Test_ADD_A_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0, add_a_n, 0, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fALL & ^fZ
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, fZ, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xFF, ld_h_n, 0x00, ld_l_n, 0x08, add_a_hl, halt, 0x01}}
	cpu = NewCPU(mem)
	cpu.reg.F = fALL & ^fZ
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, fZ|fH|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x70, ld_l_n, 0x70, add_a_l, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0xE0), cpu.reg.A)
	assert.Equal(t, fS|fY|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xF0, add_a_n, 0xB0, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0xA0), cpu.reg.A)
	assert.Equal(t, fS|fY|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x8f, add_a_n, 0x81, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.reg.A)
	assert.Equal(t, fH|fP|fC, cpu.reg.F)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x22, prefix, ld_hl_nn, 0x13, 0x00, prefix, add_a_l, halt}}
		cpu = NewCPU(mem)
		cpu.reg.F = fNONE
		cpu.Run()

		assert.Equal(t, byte(0x35), cpu.reg.A)
		assert.Equal(t, fY, cpu.reg.F)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x22, prefix, ld_hl_nn, 0x09, 0x00, prefix, add_a_hl, 0x01, halt, 0x13}}
		cpu = NewCPU(mem)
		cpu.reg.F = fNONE
		cpu.Run()

		assert.Equal(t, byte(0x35), cpu.reg.A)
		assert.Equal(t, fY, cpu.reg.F)
	}
}

func Test_ADC_A_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x20, adc_a_n, 0x20, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, byte(0x41), cpu.reg.A)
	assert.Equal(t, fNONE, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0, ld_b_n, 0xFF, adc_a_b, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fN | fC
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, fZ|fH|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x0F, ld_b_n, 0x00, adc_a_b, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.reg.A)
	assert.Equal(t, fH, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x0F, ld_l_n, 0x06, adc_a_hl, halt, 0x70}}
	cpu = NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, byte(0x80), cpu.reg.A)
	assert.Equal(t, fS|fH|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x0F, ld_b_n, 0x69, adc_a_b, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, byte(0x79), cpu.reg.A)
	assert.Equal(t, fY|fH|fX, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x0E, ld_b_n, 0x01, adc_a_b, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.reg.A)
	assert.Equal(t, fH, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x0E, ld_b_n, 0x01, adc_a_b, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x0F), cpu.reg.A)
	assert.Equal(t, fX, cpu.reg.F)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x22, prefix, ld_hl_nn, 0x13, 0x00, prefix, adc_a_l, halt}}
		cpu = NewCPU(mem)
		cpu.reg.F = fC
		cpu.Run()

		assert.Equal(t, byte(0x36), cpu.reg.A)
		assert.Equal(t, fY, cpu.reg.F)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x22, prefix, ld_hl_nn, 0x09, 0x00, prefix, adc_a_hl, 0x01, halt, 0x13}}
		cpu = NewCPU(mem)
		cpu.reg.F = fC
		cpu.Run()

		assert.Equal(t, byte(0x36), cpu.reg.A)
		assert.Equal(t, fY, cpu.reg.F)
	}
}

func Test_ADD_HL_RR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0xFF, 0xFF, ld_bc_nn, 0x01, 0, add_hl_bc, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.H)
	assert.Equal(t, byte(0), cpu.reg.L)
	assert.Equal(t, fH|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x41, 0x42, ld_de_nn, 0x11, 0x11, add_hl_de, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, byte(0x53), cpu.reg.H)
	assert.Equal(t, byte(0x52), cpu.reg.L)
	assert.Equal(t, fS|fZ|fY|fX|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x41, 0x42, add_hl_hl, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x84), cpu.reg.H)
	assert.Equal(t, byte(0x82), cpu.reg.L)
	assert.Equal(t, fNONE, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0xFE, 0xFF, ld_sp_nn, 0x02, 0, add_hl_sp, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.H)
	assert.Equal(t, byte(0), cpu.reg.L)
	assert.Equal(t, fH|fC, cpu.reg.F)

	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0xFE, 0xFF, ld_sp_nn, 0x03, 0, prefix, add_hl_sp, halt}}
		cpu = NewCPU(mem)
		cpu.reg.F = fNONE
		cpu.Run()

		assert.Equal(t, byte(0x0), *cpu.reg.prefixed[prefix][rH])
		assert.Equal(t, byte(0x01), *cpu.reg.prefixed[prefix][rL])
		assert.Equal(t, fH|fC, cpu.reg.F)
	}
}

func Test_ADC_HL_RR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{scf, ld_hl_nn, 0xFE, 0xFF, ld_bc_nn, 0x01, 0x00, prefix_ed, adc_hl_bc, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.H)
	assert.Equal(t, byte(0), cpu.reg.L)
	assert.Equal(t, fZ|fH|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{scf, ld_hl_nn, 0xC0, 0x63, ld_de_nn, 0xD0, 0x8A, prefix_ed, adc_hl_de, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xEE), cpu.reg.H)
	assert.Equal(t, byte(0x91), cpu.reg.L)
	assert.Equal(t, fS|fY|fX, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{scf, ld_hl_nn, 0x18, 0x7F, ld_de_nn, 0x48, 0x77, prefix_ed, adc_hl_de, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xF6), cpu.reg.H)
	assert.Equal(t, byte(0x61), cpu.reg.L)
	assert.Equal(t, fS|fY|fH|fP, cpu.reg.F)
}

func Test_SUB_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0, sub_n, 0x01, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fZ
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)
	assert.Equal(t, fS|fY|fH|fX|fN|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x20, sub_a, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, fZ|fN, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x90, ld_h_n, 0x20, sub_h, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fZ
	cpu.Run()

	assert.Equal(t, byte(0x70), cpu.reg.A)
	assert.Equal(t, fY|fP|fN, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x7F, ld_l_n, 0x06, sub_hl, halt, 0x80}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)
	assert.Equal(t, fS|fY|fX|fP|fN|fC, cpu.reg.F)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x12, prefix, ld_h_n, 0x02, prefix, sub_h, halt}}
		cpu = NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x10), cpu.reg.A)
		assert.Equal(t, fN, cpu.reg.F)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x12, prefix, sub_hl, 0x06, halt, 0x01}}
		cpu = NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x11), cpu.reg.A)
		assert.Equal(t, fN, cpu.reg.F)
	}
}

func Test_CP_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0, ld_b_n, 0x01, cp_b, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fZ
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, fS|fH|fN|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x20, cp_a, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x20), cpu.reg.A)
	assert.Equal(t, fZ|fY|fN, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x90, cp_n, 0x20, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fZ
	cpu.Run()

	assert.Equal(t, byte(0x90), cpu.reg.A)
	assert.Equal(t, fY|fP|fN, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x7F, ld_l_n, 0x06, cp_hl, halt, 0x80}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x7F), cpu.reg.A)
	assert.Equal(t, fS|fP|fN|fC, cpu.reg.F)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x7F, prefix, ld_h_n, 0x80, prefix, cp_h, halt}}
		cpu = NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x7F), cpu.reg.A)
		assert.Equal(t, fS|fP|fN|fC, cpu.reg.F)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x7F, prefix, cp_hl, 0x06, halt, 0x80}}
		cpu = NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x7F), cpu.reg.A)
		assert.Equal(t, fS|fP|fN|fC, cpu.reg.F)
	}
}

func Test_SBC_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x01, ld_b_n, 0x01, sbc_a_b, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)
	assert.Equal(t, fS|fY|fH|fX|fN|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x7F, ld_l_n, 0x80, sbc_a_l, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, byte(0xFE), cpu.reg.A)
	assert.Equal(t, fS|fY|fX|fP|fN|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x02, sbc_a_n, 0x01, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, fZ|fN, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x81, ld_l_n, 0x06, sbc_a_hl, halt, 0x01}}
	cpu = NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, byte(0x7F), cpu.reg.A)
	assert.Equal(t, fY|fH|fX|fP|fN, cpu.reg.F)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x12, prefix, ld_h_n, 0x02, prefix, sbc_a_h, halt}}
		cpu = NewCPU(mem)
		cpu.reg.F = fC
		cpu.Run()

		assert.Equal(t, byte(0x0F), cpu.reg.A)
		assert.Equal(t, fH|fX|fN, cpu.reg.F)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x12, prefix, sbc_a_hl, 0x06, halt, 0x01}}
		cpu = NewCPU(mem)
		cpu.reg.F = fC
		cpu.Run()

		assert.Equal(t, byte(0x10), cpu.reg.A)
		assert.Equal(t, fN, cpu.reg.F)
	}
}

func Test_SBC_HL_RR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{scf, ld_hl_nn, 0xFE, 0xFF, ld_bc_nn, 0xFD, 0xFF, prefix_ed, sbc_hl_bc, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.H)
	assert.Equal(t, byte(0), cpu.reg.L)
	assert.Equal(t, fZ|fN, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{scf, ld_hl_nn, 0x01, 0x00, ld_bc_nn, 0xFD, 0x7F, prefix_ed, sbc_hl_bc, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x80), cpu.reg.H)
	assert.Equal(t, byte(0x03), cpu.reg.L)
	assert.Equal(t, fS|fH|fN|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{scf, ld_hl_nn, 0x01, 0x70, ld_bc_nn, 0xFD, 0x8F, prefix_ed, sbc_hl_bc, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xE0), cpu.reg.H)
	assert.Equal(t, byte(0x03), cpu.reg.L)
	assert.Equal(t, fS|fY|fH|fP|fN|fC, cpu.reg.F)
}

func Test_NEG(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x55, prefix_ed, neg, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.reg.A)
	assert.Equal(t, fS|fY|fH|fX|fN|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x00, prefix_ed, neg, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, fZ|fN, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x80, prefix_ed, neg, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x80), cpu.reg.A)
	assert.Equal(t, fS|fP|fN|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xAA, prefix_ed, neg, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x56), cpu.reg.A)
	assert.Equal(t, fH|fN|fC, cpu.reg.F)
}

func Test_AND_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x0F, ld_b_n, 0xF0, and_b, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fN | fC
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, fZ|fH|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x8F, and_n, 0xF3, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x83), cpu.reg.A)
	assert.Equal(t, fS|fH, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xFF, ld_l_n, 0x06, and_hl, halt, 0x81}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x81), cpu.reg.A)
	assert.Equal(t, fS|fH|fP, cpu.reg.F)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x01, prefix, ld_l_n, 0x03, prefix, and_l, halt}}
		cpu = NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x01), cpu.reg.A)
		assert.Equal(t, fH, cpu.reg.F)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x88, prefix, and_hl, 0x06, halt, 0x08}}
		cpu = NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x08), cpu.reg.A)
		assert.Equal(t, fH|fX, cpu.reg.F)
	}
}

func Test_OR_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x00, ld_b_n, 0x00, or_b, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fS | fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, fZ|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x8A, ld_l_n, 0x06, or_hl, halt, 0x85}}
	cpu = NewCPU(mem)
	cpu.reg.F = fS | fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0x8F), cpu.reg.A)
	assert.Equal(t, fS|fX, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x11, or_n, 0x20, halt, 0x85}}
	cpu = NewCPU(mem)
	cpu.reg.F = fS | fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0x31), cpu.reg.A)
	assert.Equal(t, fY, cpu.reg.F)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x01, prefix, ld_l_n, 0x12, prefix, or_l, halt}}
		cpu = NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x13), cpu.reg.A)
		assert.Equal(t, fNONE, cpu.reg.F)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x80, prefix, or_hl, 0x06, halt, 0x08}}
		cpu = NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x88), cpu.reg.A)
		assert.Equal(t, fS|fX|fP, cpu.reg.F)
	}
}

func Test_XOR_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x1F, ld_b_n, 0x1F, xor_b, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, fZ|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x1F, ld_l_n, 0x06, xor_hl, halt, 0x8F}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x90), cpu.reg.A)
	assert.Equal(t, fS|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x1F, xor_n, 0x0F, halt, 0x8F}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.reg.A)
	assert.Equal(t, fNONE, cpu.reg.F)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x01, prefix, ld_l_n, 0x03, prefix, xor_l, halt}}
		cpu = NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x02), cpu.reg.A)
		assert.Equal(t, fNONE, cpu.reg.F)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x88, prefix, xor_hl, 0x06, halt, 0x08}}
		cpu = NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x80), cpu.reg.A)
		assert.Equal(t, fS, cpu.reg.F)
	}
}

func Test_INC_R(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0, inc_a, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, fC, cpu.reg.F)
	assert.Equal(t, byte(0x01), cpu.reg.A)

	cpu.Reset()
	cpu.reg.F = fALL & ^fZ
	mem.Cells[1] = 0xFF
	cpu.Run()

	assert.Equal(t, fZ|fH|fC, cpu.reg.F)
	assert.Equal(t, byte(0x00), cpu.reg.A)

	cpu.Reset()
	cpu.reg.F = fN
	mem.Cells[1] = 0x7F
	cpu.Run()

	assert.Equal(t, fS|fH|fP, cpu.reg.F)
	assert.Equal(t, byte(0x80), cpu.reg.A)

	cpu.Reset()
	mem.Cells[1] = 0x92
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, fS, cpu.reg.F)
	assert.Equal(t, byte(0x93), cpu.reg.A)

	cpu.Reset()
	mem.Cells[1] = 0x10
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, fNONE, cpu.reg.F)
	assert.Equal(t, byte(0x11), cpu.reg.A)

	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{
			prefix, ld_h_n, 0x10, prefix, inc_h,
			prefix, ld_l_n, 0x20, prefix, inc_l, halt},
		}
		cpu = NewCPU(mem)
		cpu.reg.F = fNONE
		cpu.Run()

		assert.Equal(t, fY, cpu.reg.F)
		assert.Equal(t, byte(0x11), *cpu.reg.prefixed[prefix][rH])
		assert.Equal(t, byte(0x21), *cpu.reg.prefixed[prefix][rL])
	}
}

func Test_INC_RR(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{
			ld_bc_nn, 0x34, 0x12, inc_bc, ld_de_nn, 0x35, 0x13, inc_de,
			ld_hl_nn, 0x36, 0x14, inc_hl, ld_sp_nn, 0x37, 0x15, inc_sp,
			useIX, ld_hl_nn, 0x38, 0x16, useIX, inc_hl,
			useIY, ld_hl_nn, 0x39, 0x17, useIY, inc_hl,
			halt},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, uint16(0x1235), cpu.reg.BC())
	assert.Equal(t, uint16(0x1336), cpu.reg.DE())
	assert.Equal(t, uint16(0x1437), cpu.reg.HL())
	assert.Equal(t, uint16(0x1538), cpu.reg.SP)
	assert.Equal(t, uint16(0x1639), uint16(cpu.reg.IXH)<<8|uint16(cpu.reg.IXL))
	assert.Equal(t, uint16(0x173A), uint16(cpu.reg.IYH)<<8|uint16(cpu.reg.IYL))
}

func Test_INC_mHL(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x05, 0x00, inc_mhl, halt, 0xFF}}
	cpu := NewCPU(mem)
	cpu.reg.F = fS | fP | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0x00), mem.Cells[5])
	assert.Equal(t, fZ|fH|fC, cpu.reg.F)

	cpu.Reset()
	mem.Cells[5] = 0x7F
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x80), mem.Cells[5])
	assert.Equal(t, fS|fH|fP, cpu.reg.F)

	cpu.Reset()
	mem.Cells[5] = 0x20
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x21), mem.Cells[5])
	assert.Equal(t, fY, cpu.reg.F)

	for _, prefix := range []byte{useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x05, 0x00, prefix, inc_mhl, 0x03, halt, 0x3F}}
		cpu := NewCPU(mem)
		cpu.reg.F = fNONE
		cpu.Run()

		assert.Equal(t, byte(0x40), mem.Cells[8])
		assert.Equal(t, fH, cpu.reg.F)
	}
}

func Test_DEC_R(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 1, dec_a, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, fZ|fN, cpu.reg.F)
	assert.Equal(t, byte(0x00), cpu.reg.A)

	cpu.Reset()
	cpu.reg.F = fALL & ^(fZ | fH | fN)
	mem.Cells[1] = 0
	cpu.Run()

	assert.Equal(t, fS|fH|fY|fX|fN|fC, cpu.reg.F)
	assert.Equal(t, byte(0xFF), cpu.reg.A)

	cpu.Reset()
	cpu.reg.F = fZ | fS
	mem.Cells[1] = 0x80
	cpu.Run()

	assert.Equal(t, fY|fH|fX|fP|fN, cpu.reg.F)
	assert.Equal(t, byte(0x7F), cpu.reg.A)

	cpu.Reset()
	cpu.reg.F = fALL
	mem.Cells[1] = 0xAB
	cpu.Run()

	assert.Equal(t, fS|fY|fX|fN|fC, cpu.reg.F)
	assert.Equal(t, byte(0xAA), cpu.reg.A)

	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{
			prefix, ld_h_n, 0x10, prefix, dec_h,
			prefix, ld_l_n, 0x20, prefix, dec_l, halt},
		}
		cpu = NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x0F), *cpu.reg.prefixed[prefix][rH])
		assert.Equal(t, byte(0x1F), *cpu.reg.prefixed[prefix][rL])
	}
}

func Test_DEC_RR(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{
			ld_bc_nn, 0x34, 0x12, dec_bc, ld_de_nn, 0x35, 0x13, dec_de,
			ld_hl_nn, 0x36, 0x14, dec_hl, ld_sp_nn, 0x37, 0x15, dec_sp,
			useIX, ld_hl_nn, 0x38, 0x16, useIX, dec_hl,
			useIY, ld_hl_nn, 0x39, 0x17, useIY, dec_hl,
			halt},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, uint16(0x1233), cpu.reg.BC())
	assert.Equal(t, uint16(0x1334), cpu.reg.DE())
	assert.Equal(t, uint16(0x1435), cpu.reg.HL())
	assert.Equal(t, uint16(0x1536), cpu.reg.SP)
	assert.Equal(t, uint16(0x1637), uint16(cpu.reg.IXH)<<8|uint16(cpu.reg.IXL))
	assert.Equal(t, uint16(0x1738), uint16(cpu.reg.IYH)<<8|uint16(cpu.reg.IYL))
}

func Test_DEC_mHL(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x05, 0x00, dec_mhl, halt, 0x00}}
	cpu := NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, byte(0xFF), mem.Cells[5])
	assert.Equal(t, fS|fY|fH|fX|fN|fC, cpu.reg.F)

	cpu.Reset()
	mem.Cells[5] = 0x01
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x00), mem.Cells[5])
	assert.Equal(t, fZ|fN, cpu.reg.F)

	cpu.Reset()
	mem.Cells[5] = 0x80
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x7F), mem.Cells[5])
	assert.Equal(t, fP|fY|fH|fX|fN, cpu.reg.F)

	for _, prefix := range []byte{useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x05, 0x00, prefix, dec_mhl, 0x03, halt, 0x3F}}
		cpu := NewCPU(mem)
		cpu.reg.F = fNONE
		cpu.Run()

		assert.Equal(t, byte(0x3E), mem.Cells[8])
		assert.Equal(t, fY|fX|fN, cpu.reg.F)
	}
}

func Test_LD_RR_nn(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_bc_nn, 0x34, 0x12, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, uint16(0x1234), cpu.reg.BC())

	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x34, 0x12, halt}}
		cpu = NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x12), *cpu.reg.prefixed[prefix][rH])
		assert.Equal(t, byte(0x34), *cpu.reg.prefixed[prefix][rL])
	}
}

func Test_LD_mm_HL(t *testing.T) {
	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x3A, 0x48, prefix, ld_mm_hl, 0x09, 0x00, halt, 0, 0}}
		cpu := NewCPU(mem)
		cpu.Run()

		assert.Equal(t, *cpu.reg.prefixed[prefix][rH], mem.Cells[10])
		assert.Equal(t, *cpu.reg.prefixed[prefix][rL], mem.Cells[9])
	}
}

func Test_LD_mm_RR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_sp_nn, 0x3A, 0x48, prefix_ed, ld_mm_sp, 0x08, 0x00, halt, 0, 0}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x3A), mem.Cells[8])
	assert.Equal(t, byte(0x48), mem.Cells[9])
}

func Test_LD_HL_mm(t *testing.T) {
	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_mm, 0x05, 0x00, halt, 0x34, 0x12}}
		cpu := NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x12), *cpu.reg.prefixed[prefix][rH])
		assert.Equal(t, byte(0x34), *cpu.reg.prefixed[prefix][rL])
	}
}

func Test_LD_RR_mm(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{prefix_ed, ld_de_mm, 0x05, 0x00, halt, 0x34, 0x12}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x12), cpu.reg.D)
	assert.Equal(t, byte(0x34), cpu.reg.E)
}

func Test_LD_mHL_n(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, ld_mhl_n, 0xAB, halt, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.mem.Read(6))

	for _, prefix := range []byte{useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x06, 0x00, prefix, ld_mhl_n, 0x03, 0xAB, halt, 0x00}}
		cpu := NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0xAB), cpu.mem.Read(9))
	}
}

func Test_LD_SP_HL(t *testing.T) {
	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x20, 0x30, prefix, ld_sp_hl, halt}}
		cpu := NewCPU(mem)
		cpu.Run()

		assert.Equal(t, uint16(0x3020), cpu.reg.SP)
	}
}

func Test_LD_mIXY_n(t *testing.T) {
	for _, prefix := range []byte{useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x02, 0x00, prefix, ld_mhl_n, 0x07, 0xAB, halt, 0x00}}
		cpu := NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0xAB), cpu.mem.Read(9))
	}
}

func Test_LD_mm_A(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x9F, ld_mm_a, 0x06, 0x00, halt, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, cpu.reg.A, mem.Cells[6])
}

func Test_LD_A_mm(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_mm, 0x04, 0x00, halt, 0xDE}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xDE), cpu.reg.A)
}

func Test_LD_BC_A(t *testing.T) {
	var n byte = 0x76
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, n, ld_bc_nn, 0x07, 0x00, ld_bc_a, halt, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, n, cpu.mem.Read(7))
}

func Test_LD_DE_A(t *testing.T) {
	var n byte = 0x76
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, n, ld_de_nn, 0x07, 0x00, ld_de_a, halt, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, n, cpu.mem.Read(7))
}

func Test_LD_A_BC(t *testing.T) {
	var n byte = 0x76
	mem := &memory.BasicMemory{Cells: []byte{ld_bc_nn, 0x05, 0x00, ld_a_bc, halt, n}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, n, cpu.reg.A)
}

func Test_LD_A_DE(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_de_nn, 0x05, 0x00, ld_a_de, halt, 0x76}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x76), cpu.reg.A)
}

func Test_LD_A_IR(t *testing.T) {
	for _, ld := range []byte{ld_a_r, ld_a_i} {
		mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x05, prefix_ed, ld, halt}}
		cpu := NewCPU(mem)
		cpu.reg.F = fH | fN | fC
		cpu.iff2 = true
		cpu.Run()

		assert.Equal(t, byte(0), cpu.reg.A)
		assert.Equal(t, fZ|fP|fC, cpu.reg.F)

		mem = &memory.BasicMemory{Cells: []byte{prefix_ed, ld, halt}}
		cpu = NewCPU(mem)
		if ld == ld_a_r {
			cpu.reg.R = 0xFF
		} else {
			cpu.reg.I = 0xFF
		}
		cpu.iff2 = false
		cpu.reg.F = fNONE
		cpu.Run()

		assert.Equal(t, byte(0xFF), cpu.reg.A)
		assert.Equal(t, fS, cpu.reg.F)
	}
}

func Test_LD_R_A(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x85, prefix_ed, ld_r_a, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x85), cpu.reg.R)
}

func Test_LD_I_A(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x85, prefix_ed, ld_i_a, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x85), cpu.reg.I)
}

func Test_LD_R_n(t *testing.T) {
	var a, b, c, d, e, h, l, ixh, ixl, iyh, iyl byte = 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	mem := &memory.BasicMemory{
		Cells: []byte{
			ld_a_n, a, ld_b_n, b, ld_c_n, c, ld_d_n, d, ld_e_n, e, ld_h_n, h, ld_l_n, l,
			useIX, ld_h_n, ixh, useIX, ld_l_n, ixl, useIY, ld_h_n, iyh, useIY, ld_l_n, iyl, halt},
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
	assert.Equal(t, ixh, cpu.reg.IXH)
	assert.Equal(t, ixl, cpu.reg.IXL)
	assert.Equal(t, iyh, cpu.reg.IYH)
	assert.Equal(t, iyl, cpu.reg.IYL)
}

func Test_LD_R_R(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{ld_a_n, 0x56, ld_b_a, ld_c_b, ld_d_c, ld_e_d, ld_h_e, ld_l_h, ld_a_n, 0,
			useIX, ld_h_b, useIX, ld_l_h, useIY, ld_l_b, useIY, ld_l_e, useIY, ld_a_l, halt},
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
	assert.Equal(t, byte(0x56), cpu.reg.IXH)
	assert.Equal(t, byte(0x56), cpu.reg.IXL)
}

func Test_LD_R_HL(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, ld_a_hl, ld_l_hl, halt, 0xA7}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xA7), cpu.reg.A)
	assert.Equal(t, byte(0xA7), cpu.reg.L)

	for _, prefix := range []byte{useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x01, 0x00, prefix, ld_l_hl, 0x07, halt, 0xA7}}
		cpu := NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0xA7), cpu.reg.L)
	}
}

func Test_LD_HL_R(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_d_n, 0x99, ld_hl_nn, 0x07, 0x00, ld_hl_d, halt, 0x00}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x99), cpu.mem.Read(7))

	for _, prefix := range []byte{useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{ld_d_n, 0x99, prefix, ld_hl_nn, 0x07, 0x00, prefix, ld_hl_d, 0x03, halt, 0x00}}
		cpu := NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x99), cpu.mem.Read(10))
	}
}

func Test_CPL(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x5B, cpl, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0xA4), cpu.reg.A)
	assert.Equal(t, fH|fY|fN, cpu.reg.F)
}

func Test_SCF(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{scf, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fS | fZ | fH | fP | fN
	cpu.Run()

	assert.Equal(t, fS|fZ|fY|fX|fP|fC, cpu.reg.F)
}

func Test_CCF(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ccf, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, fS|fZ|fY|fH|fX|fP, cpu.reg.F)

	cpu.Reset()
	cpu.reg.F = fZ | fN
	cpu.Run()

	assert.Equal(t, fZ|fY|fX|fC, cpu.reg.F)

	cpu.Reset()
	cpu.reg.F = fZ | fN | fC
	cpu.Run()

	assert.Equal(t, fZ|fY|fH|fX, cpu.reg.F)
}

func Test_RLCA(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x55, rlca, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fH | fN
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.reg.A)
	assert.Equal(t, fY|fX, cpu.reg.F)

	cpu.Reset()
	mem.Cells[1] = 0xAA
	cpu.reg.F = fH | fN
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.reg.A)
	assert.Equal(t, fC, cpu.reg.F)

	cpu.Reset()
	mem.Cells[1] = 0x00
	cpu.reg.F = fH | fN
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, fNONE, cpu.reg.F)

	cpu.Reset()
	mem.Cells[1] = 0xFF
	cpu.reg.F = fH | fN
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)
	assert.Equal(t, fY|fX|fC, cpu.reg.F)
}

func Test_RRCA(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x55, rrca, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fH | fN
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.reg.A)
	assert.Equal(t, fY|fX|fC, cpu.reg.F)

	cpu.Reset()
	mem.Cells[1] = 0xAA
	cpu.reg.F = fH | fN
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.reg.A)
	assert.Equal(t, fNONE, cpu.reg.F)

	cpu.Reset()
	mem.Cells[1] = 0x00
	cpu.reg.F = fH | fN
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, fNONE, cpu.reg.F)

	cpu.Reset()
	mem.Cells[1] = 0xFF
	cpu.reg.F = fH | fN
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)
	assert.Equal(t, fY|fX|fC, cpu.reg.F)
}

func Test_RLA(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x80, rla, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.reg.A)
	assert.Equal(t, fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x55, rla, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, byte(0xab), cpu.reg.A)
	assert.Equal(t, fS|fZ|fY|fX|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x88, rla, ld_b_a, rla, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x10), cpu.reg.B)
	assert.Equal(t, byte(0x21), cpu.reg.A)
	assert.Equal(t, fY, cpu.reg.F)
}

func Test_RRA(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x80, rra, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0xC0), cpu.reg.A)
	assert.Equal(t, fNONE, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x55, rra, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.reg.A)
	assert.Equal(t, fS|fZ|fY|fX|fP|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x89, rra, ld_b_a, rra, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x44), cpu.reg.B)
	assert.Equal(t, byte(0xA2), cpu.reg.A)
	assert.Equal(t, fY, cpu.reg.F)
}

func Test_RLD(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x7A, ld_hl_nn, 0x08, 0x00, prefix_ed, rld, halt, 0x31}}
	cpu := NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, byte(0x73), cpu.reg.A)
	assert.Equal(t, byte(0x1A), cpu.mem.Read(8))
	assert.Equal(t, fY|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x0F, ld_hl_nn, 0x08, 0x00, prefix_ed, rld, halt, 0x0A}}
	cpu = NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, byte(0xAF), cpu.mem.Read(8))
	assert.Equal(t, fZ|fP|fC, cpu.reg.F)
}

func Test_RRD(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x84, ld_hl_nn, 0x08, 0x00, prefix_ed, rrd, halt, 0x20}}
	cpu := NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, byte(0x80), cpu.reg.A)
	assert.Equal(t, byte(0x42), cpu.mem.Read(8))
	assert.Equal(t, fS|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x03, ld_hl_nn, 0x08, 0x00, prefix_ed, rrd, halt, 0x60}}
	cpu = NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, byte(0x36), cpu.mem.Read(8))
	assert.Equal(t, fZ|fP|fC, cpu.reg.F)
}

func Test_DAA(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x9A, daa, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
	assert.Equal(t, fZ|fH|fP|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x99, daa, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x99), cpu.reg.A)
	assert.Equal(t, fS|fX|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x8F, daa, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x95), cpu.reg.A)
	assert.Equal(t, fS|fH|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x8F, daa, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fN
	cpu.Run()

	assert.Equal(t, byte(0x89), cpu.reg.A)
	assert.Equal(t, fS|fX|fN, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xCA, daa, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fN
	cpu.Run()

	assert.Equal(t, byte(0x64), cpu.reg.A)
	assert.Equal(t, fY|fN|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xC5, daa, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fH | fN
	cpu.Run()

	assert.Equal(t, byte(0x5F), cpu.reg.A)
	assert.Equal(t, fH|fX|fP|fN|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xCA, daa, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, byte(0x64), cpu.reg.A)
	assert.Equal(t, fY|fN|fC, cpu.reg.F)
}

func Test_DJNZ(t *testing.T) {
	var b byte = 0x20
	var o int8 = -3
	mem := &memory.BasicMemory{Cells: []byte{ld_b_n, b, ld_a_n, 0, inc_a, djnz, byte(o), halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, b, cpu.reg.A)
	assert.Equal(t, byte(0), cpu.reg.B)

	b = 0xFF
	o = 1
	mem = &memory.BasicMemory{Cells: []byte{ld_b_n, b, ld_a_n, 1, djnz, byte(o), halt, inc_a, jr_o, 0xFA}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, b, cpu.reg.A)
	assert.Equal(t, byte(0), cpu.reg.B)
}

func Test_JR_o(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{jr_o, 3, ld_c_n, 0x11, halt, ld_d_n, 0x22, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.C)
	assert.Equal(t, byte(0x22), cpu.reg.D)

	mem = &memory.BasicMemory{Cells: []byte{jr_o, 6, halt, ld_c_n, 0x11, ld_b_n, 0x33, halt, jr_o, 0xF9, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x33), cpu.reg.B)
	assert.Equal(t, byte(0x11), cpu.reg.C)
}

func Test_JR_Z_o(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 2, dec_a, jr_z_o, 0x02, ld_b_n, 0xab, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xab), cpu.reg.B)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 1, dec_a, jr_z_o, 0x02, ld_b_n, 0xab, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 1, dec_a, jr_z_o, 0xFD, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)
}

func Test_JR_NZ_o(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 2, dec_a, jr_nz_o, 0x02, ld_b_n, 0xab, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.B)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 1, dec_a, jr_nz_o, 0x02, ld_b_n, 0xab, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xab), cpu.reg.B)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 2, dec_a, jr_nz_o, 0xFD, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
}

func Test_JR_C(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_b_n, 0xAB, dec_a, add_a_n, 1, jr_c, 1, ld_b_a, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.B)

	mem = &memory.BasicMemory{Cells: []byte{ld_b_n, 0xAB, scf, jr_c, 1, halt, ld_b_a, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.B)

	mem = &memory.BasicMemory{Cells: []byte{xor_a, dec_a, add_a_n, 1, jr_c, 0xFC, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(1), cpu.reg.A)
}

func Test_JR_NC_o(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_b_n, 0xAB, inc_a, add_a_n, 1, jr_nc_o, 1, ld_b_a, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.reg.B)

	mem = &memory.BasicMemory{Cells: []byte{ld_b_n, 0xAB, scf, ccf, jr_nc_o, 1, halt, ld_b_a, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.B)

	mem = &memory.BasicMemory{Cells: []byte{add_a_n, 1, jr_nc_o, 0xFC, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.A)
}

func Test_JP_nn(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{jp_nn, 0x06, 0x00, ld_a_n, 0xAA, halt, ld_a_n, 0x55, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.reg.A)
}

func Test_JP_HL(t *testing.T) {
	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x09, 0x00, prefix, jp_hl, ld_a_n, 0xAA, halt, ld_a_n, 0x55, halt}}
		cpu := NewCPU(mem)
		cpu.Run()

		assert.Equal(t, byte(0x55), cpu.reg.A)
	}
}

func Test_JP_cc_nn(t *testing.T) {
	var tests = []struct {
		jp       byte
		flag     byte
		expected byte
	}{
		{jp_c_nn, fC, 0x55}, {jp_nc_nn, fNONE, 0x55}, {jp_z_nn, fZ, 0x55}, {jp_nz_nn, fNONE, 0x55},
		{jp_m_nn, fS, 0x55}, {jp_p_nn, fNONE, 0x55}, {jp_pe_nn, fP, 0x55}, {jp_po_nn, fNONE, 0x55},
		{jp_c_nn, fNONE, 0xAA}, {jp_nc_nn, fC, 0xAA}, {jp_z_nn, fNONE, 0xAA}, {jp_nz_nn, fZ, 0xAA},
		{jp_m_nn, fNONE, 0xAA}, {jp_p_nn, fS, 0xAA}, {jp_pe_nn, fNONE, 0xAA}, {jp_po_nn, fP, 0xAA},
	}

	for _, test := range tests {
		mem := &memory.BasicMemory{
			Cells: []byte{test.jp, 0x06, 0x00, ld_a_n, 0xAA, halt, ld_a_n, 0x55, halt, 0x07, 0x00},
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
		{call_c_nn, fC, 0x55}, {call_nc_nn, fNONE, 0x55}, {call_z_nn, fZ, 0x55}, {call_nz_nn, fNONE, 0x55},
		{call_m_nn, fS, 0x55}, {call_p_nn, fNONE, 0x55}, {call_pe_nn, fP, 0x55}, {call_po_nn, fNONE, 0x55},
		{call_c_nn, fNONE, 0xAA}, {call_nc_nn, fC, 0xAA}, {call_z_nn, fNONE, 0xAA}, {call_nz_nn, fZ, 0xAA},
		{call_m_nn, fNONE, 0xAA}, {call_p_nn, fS, 0xAA}, {call_pe_nn, fNONE, 0xAA}, {call_po_nn, fP, 0xAA},
		{call_nn, fNONE, 0x55},
	}

	for _, test := range tests {
		mem := &memory.BasicMemory{
			Cells: []byte{ld_sp_nn, 0x10, 0x00, test.call, 0x09, 0x00, ld_a_n, 0xAA, halt, ld_a_n, 0x55, halt, 0xFF, 0xFF, 0xFF, 0xFF},
		}

		cpu := NewCPU(mem)
		cpu.reg.F = test.flag
		cpu.Run()

		assert.Equal(t, byte(test.expected), cpu.reg.A)
		if cpu.reg.A == 0x55 {
			assert.Equal(t, uint16(0x0E), cpu.reg.SP)
			assert.Equal(t, byte(0), mem.Cells[15])
			assert.Equal(t, byte(0x06), mem.Cells[14])
		} else {
			assert.Equal(t, uint16(0x10), cpu.reg.SP)
		}
	}
}

func Test_RET(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_sp_nn, 0x0A, 0x00, ret, ld_a_n, 0xAA, halt, ld_a_n, 0x55, halt, 0x07, 0x00}}
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
		{ret_c, fC, 0x55}, {ret_nc, fNONE, 0x55}, {ret_z, fZ, 0x55}, {ret_nz, fNONE, 0x55},
		{ret_m, fS, 0x55}, {ret_p, fNONE, 0x55}, {ret_pe, fP, 0x55}, {ret_po, fNONE, 0x55},
		{ret_c, fNONE, 0xAA}, {ret_nc, fC, 0xAA}, {ret_z, fNONE, 0xAA}, {ret_nz, fZ, 0xAA},
		{ret_m, fNONE, 0xAA}, {ret_p, fS, 0xAA}, {ret_pe, fNONE, 0xAA}, {ret_po, fP, 0xAA},
	}

	for _, test := range tests {
		mem := &memory.BasicMemory{
			Cells: []byte{ld_sp_nn, 0x0A, 0x00, test.ret, ld_a_n, 0xAA, halt, ld_a_n, 0x55, halt, 0x07, 0x00},
		}
		cpu := NewCPU(mem)
		cpu.reg.F = test.flag
		cpu.Run()

		assert.Equal(t, byte(test.expected), cpu.reg.A)
	}
}

func Test_RETN_RETI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_sp_nn, 0x0B, 0x00, prefix_ed, retn, ld_a_n, 0xAA, halt, ld_a_n, 0x55, halt, 0x08, 0x00}}
	cpu := NewCPU(mem)
	cpu.iff2 = true
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.reg.A)
	assert.Equal(t, true, cpu.iff1)
	assert.Equal(t, true, cpu.iff2)
}

func Test_RST_xx(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{
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
	cpu.reg.PC = 0x40
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)
	assert.Equal(t, byte(0x55), cpu.reg.B)
}

func Test_PUSH_rr(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{ld_sp_nn, 0x23, 0x00, ld_a_n, 0x98,
			ld_bc_nn, 0x34, 0x12, ld_de_nn, 0x35, 0x13, ld_hl_nn, 0x36, 0x14,
			push_af, push_bc, push_de, push_hl, useIX, push_hl, useIY, push_hl, halt,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	cpu := NewCPU(mem)
	cpu.reg.F = fS | fC
	cpu.Run()

	assert.Equal(t, cpu.mem.Read(34), cpu.reg.A)
	assert.Equal(t, cpu.mem.Read(33), cpu.reg.F)
	assert.Equal(t, cpu.mem.Read(32), cpu.reg.B)
	assert.Equal(t, cpu.mem.Read(31), cpu.reg.C)
	assert.Equal(t, cpu.mem.Read(30), cpu.reg.D)
	assert.Equal(t, cpu.mem.Read(29), cpu.reg.E)
	assert.Equal(t, cpu.mem.Read(28), cpu.reg.H)
	assert.Equal(t, cpu.mem.Read(27), cpu.reg.L)
	assert.Equal(t, cpu.mem.Read(26), cpu.reg.IXH)
	assert.Equal(t, cpu.mem.Read(25), cpu.reg.IXL)
	assert.Equal(t, cpu.mem.Read(24), cpu.reg.IYH)
	assert.Equal(t, cpu.mem.Read(23), cpu.reg.IYL)
}

func Test_POP_rr(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{ld_sp_nn, 0x0C, 0x00, pop_af, pop_bc, pop_de, pop_hl, useIX, pop_hl, useIY, pop_hl, halt,
			0x43, 0x21, 0x44, 0x22, 0x45, 0x23, 0x46, 0x24, 0x47, 0x25, 0x48, 0x26},
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
	assert.Equal(t, byte(0x25), cpu.reg.IXH)
	assert.Equal(t, byte(0x47), cpu.reg.IXL)
	assert.Equal(t, byte(0x26), cpu.reg.IYH)
	assert.Equal(t, byte(0x48), cpu.reg.IYL)
	assert.Equal(t, uint16(0x18), cpu.reg.SP)
}

func Test_IN_A_n(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{ld_a_n, 0x23, in_a_n, 0x01, halt},
	}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.A)

	cpu.Reset()
	cpu.IN = func(hi, lo byte) byte {
		if hi == 0x23 && lo == 0x01 {
			return 0xA5
		}
		return 0
	}
	cpu.Run()

	assert.Equal(t, byte(0xA5), cpu.reg.A)
}

func Test_IN_R_C(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{ld_bc_nn, 0x23, 0x01, ld_d_n, 0x01, prefix_ed, in_d_c, halt},
	}
	cpu := NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, byte(0xFF), cpu.reg.D)
	assert.Equal(t, fS|fP|fC, cpu.reg.F)

	cpu.Reset()
	cpu.reg.F = fNONE
	cpu.IN = func(hi, lo byte) byte {
		if hi == 0x01 && lo == 0x23 {
			return 0
		}
		return 0xA5
	}
	cpu.Run()

	assert.Equal(t, byte(0), cpu.reg.D)
	assert.Equal(t, fZ|fP, cpu.reg.F)
}

func Test_OUT_n_A(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{ld_a_n, 0x23, out_n_a, 0x01, halt},
	}
	cpu := NewCPU(mem)
	cpu.OUT = func(hi, lo, data byte) {
		assert.Equal(t, byte(0x23), hi)
		assert.Equal(t, byte(0x01), lo)
		assert.Equal(t, byte(0x23), data)
	}
	cpu.Run()
}

func Test_OUT_C_R(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{ld_bc_nn, 0x11, 0x22, ld_h_n, 0x33, prefix_ed, out_c_h, halt},
	}
	cpu := NewCPU(mem)
	cpu.OUT = func(hi, lo, data byte) {
		assert.Equal(t, byte(0x22), hi)
		assert.Equal(t, byte(0x11), lo)
		assert.Equal(t, byte(0x33), data)
	}
	cpu.Run()
}

func Test_RLC_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x55, prefix_cb, rlc_r | rE, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fZ | fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.reg.E)
	assert.Equal(t, fS|fY|fX|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_d_n, 0xAA, prefix_cb, rlc_r | rD, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fZ | fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.reg.D)
	assert.Equal(t, fP|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x00, prefix_cb, rlc_r | rA, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, fZ|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_b_n, 0x80, prefix_cb, rlc_r | rB, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.reg.B)
	assert.Equal(t, fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_cb, rlc_r | 0b110, halt, 0x01}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x02), cpu.mem.Read(0x06))
	assert.Equal(t, fNONE, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{useIX, ld_hl_nn, 0x04, 0x00, useIX, prefix_cb, 0x05, rlc_r | 0b110, halt, 0x01}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x02), cpu.mem.Read(0x09))
	assert.Equal(t, fNONE, cpu.reg.F)
}

func Test_RRC_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x55, prefix_cb, rrc_r | rE, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fZ | fH | fN
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.reg.E)
	assert.Equal(t, fS|fY|fX|fP|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_d_n, 0xAA, prefix_cb, rrc_r | rD, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fS | fZ | fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.reg.D)
	assert.Equal(t, fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x00, prefix_cb, rrc_r | rA, halt}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, fZ|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_b_n, 0x80, prefix_cb, rrc_r | rB, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, byte(0x40), cpu.reg.B)
	assert.Equal(t, fNONE, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_cb, rrc_r | 0b110, halt, 0x01}}
	cpu = NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x80), cpu.mem.Read(0x06))
	assert.Equal(t, fS|fC, cpu.reg.F)
}

func Test_RL_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x55, prefix_cb, rl_r | rE, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fZ | fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0xAB), cpu.reg.E)
	assert.Equal(t, fS|fY|fX, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_d_n, 0xAA, prefix_cb, rl_r | rD, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fZ | fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0x55), cpu.reg.D)
	assert.Equal(t, fP|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x80, prefix_cb, rl_r | rA, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, fZ|fP|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_b_n, 0x80, prefix_cb, rl_r | rB, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.reg.B)
	assert.Equal(t, fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_cb, rl_r | 0b110, halt, 0x81}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x02), cpu.mem.Read(0x06))
	assert.Equal(t, fC, cpu.reg.F)
}

func Test_RR_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x55, prefix_cb, rr_r | rE, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fZ | fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.reg.E)
	assert.Equal(t, fS|fY|fX|fP|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_d_n, 0xAA, prefix_cb, rr_r | rD, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fZ | fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0xD5), cpu.reg.D)
	assert.Equal(t, fS, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x01, prefix_cb, rr_r | rA, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.reg.A)
	assert.Equal(t, fZ|fP|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_b_n, 0x80, prefix_cb, rr_r | rB, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, byte(0xC0), cpu.reg.B)
	assert.Equal(t, fS|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_cb, rr_r | 0b110, halt, 0x81}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0x40), cpu.mem.Read(0x06))
	assert.Equal(t, fC, cpu.reg.F)
}

func Test_SLA_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x55, prefix_cb, sla_r | rE, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fZ | fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0xAA), cpu.reg.E)
	assert.Equal(t, fS|fY|fX|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_d_n, 0xAA, prefix_cb, sla_r | rD, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fZ | fH | fN
	cpu.Run()

	assert.Equal(t, byte(0x54), cpu.reg.D)
	assert.Equal(t, fC, cpu.reg.F)
}

func Test_SRA_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x85, prefix_cb, sra_r | rE, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fZ | fH | fN
	cpu.Run()

	assert.Equal(t, byte(0xC2), cpu.reg.E)
	assert.Equal(t, fS|fC, cpu.reg.F)
}

func Test_SLL_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x95, prefix_cb, sll_r | rE, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fS | fZ | fH | fN | fC
	cpu.Run()

	assert.Equal(t, byte(0x2B), cpu.reg.E)
	assert.Equal(t, fY|fX|fP|fC, cpu.reg.F)
}

func Test_SRL_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_h_n, 0x85, prefix_cb, srl_r | rH, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fS | fH | fN
	cpu.Run()

	assert.Equal(t, byte(0x42), cpu.reg.H)
	assert.Equal(t, fP|fC, cpu.reg.F)
}

func Test_BIT_b(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x40, prefix_cb, bit_b | rE | bit_6, halt}}
	cpu := NewCPU(mem)
	cpu.reg.F = fZ | fN | fC
	cpu.Run()

	assert.Equal(t, fH|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_l_n, 0xFE, prefix_cb, bit_b | rL | bit_0, halt}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, fZ|fY|fH|fX|fP, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_cb, bit_b | 0b110 | bit_2, halt, 0xFD}}
	cpu = NewCPU(mem)
	cpu.reg.F = fZ | fN
	cpu.Run()

	assert.Equal(t, fH, cpu.reg.F)
}

func Test_RES_b(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_d_n, 0xFF, prefix_cb, res_b | rD | bit_7, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x7F), cpu.reg.D)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_cb, res_b | 0b110 | bit_2, halt, 0xFF}}
	cpu = NewCPU(mem)
	cpu.reg.F = fZ | fN
	cpu.Run()

	assert.Equal(t, byte(0xFB), cpu.mem.Read(0x06))
}

func Test_SET_b(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_d_n, 0x00, prefix_cb, set_b | rD | bit_7, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x80), cpu.reg.D)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_cb, set_b | 0b110 | bit_2, halt, 0x00}}
	cpu = NewCPU(mem)
	cpu.reg.F = fZ | fN
	cpu.Run()

	assert.Equal(t, byte(0x04), cpu.mem.Read(0x06))

	mem = &memory.BasicMemory{Cells: []byte{useIY, ld_hl_nn, 0x06, 0x00, useIY, prefix_cb, 0x03, set_b | bit_2, halt, 0x00}}
	cpu = NewCPU(mem)
	cpu.reg.F = fZ | fN
	cpu.Run()

	assert.Equal(t, byte(0x04), cpu.mem.Read(0x09))
	assert.Equal(t, byte(0x04), cpu.reg.B)
}

func Test_LD_IXY_nn(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		useIX, ld_hl_nn, 0x06, 0x01, useIY, ld_hl_nn, 0x07, 0x02,
		ld_hl_nn, 0x08, 0x03, halt}}
	cpu := NewCPU(mem)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.reg.IXH)
	assert.Equal(t, byte(0x06), cpu.reg.IXL)
	assert.Equal(t, byte(0x02), cpu.reg.IYH)
	assert.Equal(t, byte(0x07), cpu.reg.IYL)
	assert.Equal(t, byte(0x03), cpu.reg.H)
	assert.Equal(t, byte(0x08), cpu.reg.L)
}

func Test_LDI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0C, 0x00, ld_de_nn, 0x0D, 0x00, ld_bc_nn, 0x01, 0x00,
		prefix_ed, ldi, halt, 0xA5, 0x00}}
	cpu := NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, byte(0xA5), cpu.mem.Read(0x0D))
	assert.Equal(t, uint16(0), cpu.reg.BC())
	assert.Equal(t, uint16(0x0E), cpu.reg.DE())
	assert.Equal(t, uint16(0x0D), cpu.reg.HL())
	assert.Equal(t, fS|fZ|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0C, 0x00, ld_de_nn, 0x0D, 0x00, ld_bc_nn, 0x02, 0x00,
		prefix_ed, ldi, halt, 0xA5, 0x00}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0xA5), cpu.mem.Read(0x0D))
	assert.Equal(t, uint16(1), cpu.reg.BC())
	assert.Equal(t, uint16(0x0E), cpu.reg.DE())
	assert.Equal(t, uint16(0x0D), cpu.reg.HL())
	assert.Equal(t, fP, cpu.reg.F)
}

func Test_LDIR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0C, 0x00, ld_de_nn, 0x0F, 0x00, ld_bc_nn, 0x03, 0x00,
		prefix_ed, ldir, halt, 0x88, 0x36, 0xA5, 0x00, 0x00, 0x00}}
	cpu := NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, byte(0x88), cpu.mem.Read(0x0F))
	assert.Equal(t, byte(0x36), cpu.mem.Read(0x10))
	assert.Equal(t, byte(0xA5), cpu.mem.Read(0x11))
	assert.Equal(t, uint16(0), cpu.reg.BC())
	assert.Equal(t, uint16(0x12), cpu.reg.DE())
	assert.Equal(t, uint16(0x0F), cpu.reg.HL())
	assert.Equal(t, fS|fZ|fC, cpu.reg.F)
}

func Test_CPI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0B, 0x00, ld_bc_nn, 0x03, 0x00, ld_a_n, 0x88,
		prefix_ed, cpi, halt, 0x88}}
	cpu := NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, uint16(0x02), cpu.reg.BC())
	assert.Equal(t, uint16(0x0C), cpu.reg.HL())
	assert.Equal(t, fZ|fP|fN|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0B, 0x00, ld_bc_nn, 0x01, 0x00, ld_a_n, 0x88,
		prefix_ed, cpi, halt, 0x89}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, uint16(0x00), cpu.reg.BC())
	assert.Equal(t, uint16(0x0C), cpu.reg.HL())
	assert.Equal(t, fS|fY|fH|fX|fN, cpu.reg.F)
}

func Test_CPIR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0B, 0x00, ld_bc_nn, 0xFF, 0x00, ld_a_n, 0x88,
		prefix_ed, cpir, halt, 0x02, 0x04, 0x80, 0x88, 0x90}}
	cpu := NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, uint16(0xFB), cpu.reg.BC())
	assert.Equal(t, uint16(0x0F), cpu.reg.HL())
	assert.Equal(t, fZ|fP|fN, cpu.reg.F)
}

func Test_INI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x09, 0x00, ld_bc_nn, 0x34, 0x01,
		prefix_ed, ini, halt, 0x00}}
	cpu := NewCPU(mem)
	cpu.reg.F = fC
	cpu.IN = func(hi, lo byte) byte {
		if hi == 0x01 && lo == 0x34 {
			return 0x5E
		}
		return 0
	}
	cpu.Run()

	assert.Equal(t, uint16(0x34), cpu.reg.BC())
	assert.Equal(t, byte(0x5E), cpu.mem.Read((9)))
	assert.Equal(t, uint16(0x0A), cpu.reg.HL())
	assert.Equal(t, fZ|fN|fC, cpu.reg.F)
}

func Test_INIR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x09, 0x00, ld_bc_nn, 0x34, 0x05,
		prefix_ed, inir, halt, 0x00, 0x00, 0x00, 0x00, 0x00}}
	cpu := NewCPU(mem)
	cpu.reg.F = fC
	cpu.IN = func(hi, lo byte) byte {
		if lo == 0x34 {
			return hi + 0x20
		}
		return 0
	}
	cpu.Run()

	assert.Equal(t, uint16(0x34), cpu.reg.BC())
	assert.Equal(t, byte(0x25), cpu.mem.Read((9)))
	assert.Equal(t, byte(0x24), cpu.mem.Read((10)))
	assert.Equal(t, byte(0x23), cpu.mem.Read((11)))
	assert.Equal(t, byte(0x22), cpu.mem.Read((12)))
	assert.Equal(t, byte(0x21), cpu.mem.Read((13)))
	assert.Equal(t, uint16(0x0E), cpu.reg.HL())
	assert.Equal(t, fZ|fN|fC, cpu.reg.F)
}

func Test_OUTI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x09, 0x00, ld_bc_nn, 0x34, 0x01,
		prefix_ed, outi, halt, 0x87}}
	cpu := NewCPU(mem)
	cpu.reg.F = fC
	cpu.OUT = func(hi, lo, data byte) {
		assert.Equal(t, byte(0), hi)
		assert.Equal(t, byte(0x34), lo)
		assert.Equal(t, byte(0x87), data)
	}
	cpu.Run()

	assert.Equal(t, uint16(0x34), cpu.reg.BC())
	assert.Equal(t, uint16(0x0A), cpu.reg.HL())
	assert.Equal(t, fZ|fN|fC, cpu.reg.F)
}

func Test_OTIR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x09, 0x00, ld_bc_nn, 0x34, 0x04,
		prefix_ed, otir, halt, 0x87, 0x88, 0x89, 0x8A}}
	cpu := NewCPU(mem)
	cpu.reg.F = fC
	cpu.OUT = func(hi, lo, data byte) {
		assert.Equal(t, byte(0x34), lo)
		assert.Equal(t, byte(0x87+0x03-hi), data)
	}
	cpu.Run()

	assert.Equal(t, uint16(0x34), cpu.reg.BC())
	assert.Equal(t, uint16(0x0D), cpu.reg.HL())
	assert.Equal(t, fZ|fN|fC, cpu.reg.F)
}

func Test_LDD(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0C, 0x00, ld_de_nn, 0x0D, 0x00, ld_bc_nn, 0x01, 0x00,
		prefix_ed, ldd, halt, 0xA5, 0x00}}
	cpu := NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, byte(0xA5), cpu.mem.Read(0x0D))
	assert.Equal(t, uint16(0x00), cpu.reg.BC())
	assert.Equal(t, uint16(0x0C), cpu.reg.DE())
	assert.Equal(t, uint16(0x0B), cpu.reg.HL())
	assert.Equal(t, fS|fZ|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0C, 0x00, ld_de_nn, 0x0D, 0x00, ld_bc_nn, 0x02, 0x00,
		prefix_ed, ldd, halt, 0xA5, 0x00}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, byte(0xA5), cpu.mem.Read(0x0D))
	assert.Equal(t, uint16(0x01), cpu.reg.BC())
	assert.Equal(t, uint16(0x0C), cpu.reg.DE())
	assert.Equal(t, uint16(0x0B), cpu.reg.HL())
	assert.Equal(t, fP, cpu.reg.F)
}

func Test_LDDR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0E, 0x00, ld_de_nn, 0x11, 0x00, ld_bc_nn, 0x03, 0x00,
		prefix_ed, lddr, halt, 0x88, 0x36, 0xA5, 0x00, 0x00, 0x00}}
	cpu := NewCPU(mem)
	cpu.reg.F = fALL
	cpu.Run()

	assert.Equal(t, byte(0x88), cpu.mem.Read(0x0F))
	assert.Equal(t, byte(0x36), cpu.mem.Read(0x10))
	assert.Equal(t, byte(0xA5), cpu.mem.Read(0x11))
	assert.Equal(t, uint16(0x00), cpu.reg.BC())
	assert.Equal(t, uint16(0x0E), cpu.reg.DE())
	assert.Equal(t, uint16(0x0B), cpu.reg.HL())
	assert.Equal(t, fS|fZ|fY|fC, cpu.reg.F)
}

func Test_CPD(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0B, 0x00, ld_bc_nn, 0x03, 0x00, ld_a_n, 0x88,
		prefix_ed, cpd, halt, 0x88}}
	cpu := NewCPU(mem)
	cpu.reg.F = fC
	cpu.Run()

	assert.Equal(t, uint16(0x02), cpu.reg.BC())
	assert.Equal(t, uint16(0x0A), cpu.reg.HL())
	assert.Equal(t, fZ|fP|fN|fC, cpu.reg.F)

	mem = &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0B, 0x00, ld_bc_nn, 0x01, 0x00, ld_a_n, 0x88,
		prefix_ed, cpd, halt, 0x89}}
	cpu = NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, uint16(0x00), cpu.reg.BC())
	assert.Equal(t, uint16(0x0A), cpu.reg.HL())
	assert.Equal(t, fS|fY|fH|fX|fN, cpu.reg.F)
}

func Test_CPDR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0F, 0x00, ld_bc_nn, 0xFF, 0x00, ld_a_n, 0x88,
		prefix_ed, cpdr, halt, 0x02, 0x04, 0x80, 0x88, 0x90}}
	cpu := NewCPU(mem)
	cpu.reg.F = fNONE
	cpu.Run()

	assert.Equal(t, uint16(0xFD), cpu.reg.BC())
	assert.Equal(t, uint16(0x0D), cpu.reg.HL())
	assert.Equal(t, fZ|fP|fN, cpu.reg.F)
}

func Test_IND(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x09, 0x00, ld_bc_nn, 0x34, 0x01,
		prefix_ed, ind, halt, 0x00}}
	cpu := NewCPU(mem)
	cpu.reg.F = fC
	cpu.IN = func(hi, lo byte) byte {
		if hi == 0x01 && lo == 0x34 {
			return 0x5E
		}
		return 0
	}
	cpu.Run()

	assert.Equal(t, uint16(0x34), cpu.reg.BC())
	assert.Equal(t, byte(0x5E), cpu.mem.Read((9)))
	assert.Equal(t, uint16(0x08), cpu.reg.HL())
	assert.Equal(t, fZ|fN|fC, cpu.reg.F)
}

func Test_INDR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0D, 0x00, ld_bc_nn, 0x34, 0x05,
		prefix_ed, indr, halt, 0x00, 0x00, 0x00, 0x00, 0x00}}
	cpu := NewCPU(mem)
	cpu.reg.F = fC
	cpu.IN = func(hi, lo byte) byte {
		if lo == 0x34 {
			return hi + 0x20
		}
		return 0
	}
	cpu.Run()

	assert.Equal(t, uint16(0x34), cpu.reg.BC())
	assert.Equal(t, byte(0x21), cpu.mem.Read((9)))
	assert.Equal(t, byte(0x22), cpu.mem.Read((10)))
	assert.Equal(t, byte(0x23), cpu.mem.Read((11)))
	assert.Equal(t, byte(0x24), cpu.mem.Read((12)))
	assert.Equal(t, byte(0x25), cpu.mem.Read((13)))
	assert.Equal(t, uint16(0x08), cpu.reg.HL())
	assert.Equal(t, fZ|fN|fC, cpu.reg.F)
}

func Test_OUTD(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x09, 0x00, ld_bc_nn, 0x34, 0x01,
		prefix_ed, outd, halt, 0x87}}
	cpu := NewCPU(mem)
	cpu.reg.F = fC
	cpu.OUT = func(hi, lo, data byte) {
		assert.Equal(t, byte(0), hi)
		assert.Equal(t, byte(0x34), lo)
		assert.Equal(t, byte(0x87), data)
	}
	cpu.Run()

	assert.Equal(t, uint16(0x34), cpu.reg.BC())
	assert.Equal(t, uint16(0x08), cpu.reg.HL())
	assert.Equal(t, fZ|fN|fC, cpu.reg.F)
}

func Test_OTDR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0C, 0x00, ld_bc_nn, 0x34, 0x04,
		prefix_ed, otdr, halt, 0x87, 0x88, 0x89, 0x8A}}
	cpu := NewCPU(mem)
	cpu.reg.F = fC
	cpu.OUT = func(hi, lo, data byte) {
		assert.Equal(t, byte(0x34), lo)
		assert.Equal(t, byte(0x8A-(0x03-hi)), data)
	}
	cpu.Run()

	assert.Equal(t, uint16(0x34), cpu.reg.BC())
	assert.Equal(t, uint16(0x08), cpu.reg.HL())
	assert.Equal(t, fZ|fN|fC, cpu.reg.F)
}

func Test_shouldJump(t *testing.T) {
	var tests = []struct {
		flags    byte
		code     byte
		expected bool
	}{
		{fNONE, 0b00000000, true}, {fZ, 0b00000000, false},
		{fNONE, 0b00001000, false}, {fZ, 0b00001000, true},
		{fNONE, 0b00010000, true}, {fC, 0b00010000, false},
		{fNONE, 0b00011000, false}, {fC, 0b00011000, true},
		{fNONE, 0b00100000, true}, {fP, 0b00100000, false},
		{fNONE, 0b00101000, false}, {fP, 0b00101000, true},
		{fNONE, 0b00110000, true}, {fS, 0b00110000, false},
		{fNONE, 0b00111000, false}, {fS, 0b00111000, true},
	}

	cpu := NewCPU(&memory.BasicMemory{})
	for _, test := range tests {
		cpu.reg.F = test.flags
		result := cpu.shouldJump(test.code)

		assert.Equal(t, test.expected, result)
	}
}

func Test_getHL(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x34, 0x12, useIX, ld_hl_nn, 0x45, 0x23, halt, 0x00, 0x02, 0xFE}}
	cpu := NewCPU(mem)
	cpu.Run()

	hl := cpu.getHL()
	assert.Equal(t, uint16(0x1234), hl)

	hl = cpu.getHL()
	assert.Equal(t, uint16(0x1234), hl)

	cpu.reg.prefix = useIX
	hl = cpu.getHL()
	assert.Equal(t, uint16(0x2345), hl)

	hl = cpu.getHL()
	assert.Equal(t, uint16(0x2347), hl)

	hl = cpu.getHL()
	assert.Equal(t, uint16(0x2343), hl)
}

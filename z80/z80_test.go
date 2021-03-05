package z80

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voytas/z80-go-zx/z80/memory"
)

type TestIOBus struct {
	read  func(hi, lo byte) byte
	write func(hi, lo, data byte, t int)
}

func (bus *TestIOBus) Read(hi, lo byte) byte          { return bus.read(hi, lo) }
func (bus *TestIOBus) Write(hi, lo, data byte, t int) { bus.write(hi, lo, data, t) }

func Test_NOP(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{nop}}
	z80 := NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(4)

	assert.Equal(t, fALL, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_HALT(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{xor_a, halt, ld_a_n, 0x87}}
	z80 := NewZ80(mem)
	z80.Run(4 + 4 + 7)

	assert.Equal(t, true, z80.halt)
	assert.Equal(t, byte(0x00), z80.reg.A)
}

func Test_DI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{di}}
	z80 := NewZ80(mem)
	z80.Run(4)

	assert.Equal(t, false, z80.iff1)
	assert.Equal(t, false, z80.iff2)
	assert.Equal(t, 0, z80.t)
}

func Test_EI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{di, ei}}
	z80 := NewZ80(mem)
	z80.Run(4 + 4)

	assert.Equal(t, true, z80.iff1)
	assert.Equal(t, true, z80.iff2)
	assert.Equal(t, 0, z80.t)
}

func Test_IM_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{prefix_ed, im1}}
	z80 := NewZ80(mem)
	z80.Run(8)

	assert.Equal(t, im1, z80.im)
	assert.Equal(t, 0, z80.t)
}

func Test_EX_AF_AF(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ex_af_af}}
	z80 := NewZ80(mem)
	var a, a_, f, f_ byte = 0xcc, 0x55, 0x35, 0x97
	z80.reg.A, z80.reg.A_ = a, a_
	z80.reg.F = f
	z80.reg.F_ = f_
	z80.Run(4)

	assert.Equal(t, a_, z80.reg.A)
	assert.Equal(t, a, z80.reg.A_)
	assert.Equal(t, f_, z80.reg.F)
	assert.Equal(t, f, z80.reg.F_)
	assert.Equal(t, 0, z80.t)

	z80.reg.PC = 0
	z80.Run(4)
	assert.Equal(t, a, z80.reg.A)
	assert.Equal(t, a_, z80.reg.A_)
	assert.Equal(t, f, z80.reg.F)
	assert.Equal(t, f_, z80.reg.F_)
	assert.Equal(t, 0, z80.t)
}

func Test_EXX(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{exx}}
	z80 := NewZ80(mem)
	z80.reg.B, z80.reg.C, z80.reg.B_, z80.reg.C_ = 0x01, 0x02, 0x03, 0x04
	z80.reg.D, z80.reg.E, z80.reg.D_, z80.reg.E_ = 0x05, 0x06, 0x07, 0x08
	z80.reg.H, z80.reg.L, z80.reg.H_, z80.reg.L_ = 0x09, 0x0A, 0x0B, 0x0C
	z80.Run(4)
	assert.Equal(t, 0, z80.t)

	assert.Equal(t, byte(0x01), z80.reg.B_)
	assert.Equal(t, byte(0x02), z80.reg.C_)
	assert.Equal(t, byte(0x03), z80.reg.B)
	assert.Equal(t, byte(0x04), z80.reg.C)
	assert.Equal(t, byte(0x05), z80.reg.D_)
	assert.Equal(t, byte(0x06), z80.reg.E_)
	assert.Equal(t, byte(0x07), z80.reg.D)
	assert.Equal(t, byte(0x08), z80.reg.E)
	assert.Equal(t, byte(0x09), z80.reg.H_)
	assert.Equal(t, byte(0x0A), z80.reg.L_)
	assert.Equal(t, byte(0x0B), z80.reg.H)
	assert.Equal(t, byte(0x0C), z80.reg.L)
}

func Test_EX_DE_HL(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ex_de_hl}}
	z80 := NewZ80(mem)
	z80.reg.D, z80.reg.E, z80.reg.H, z80.reg.L = 0x01, 0x02, 0x03, 0x04
	z80.Run(4)

	assert.Equal(t, byte(0x01), z80.reg.H)
	assert.Equal(t, byte(0x02), z80.reg.L)
	assert.Equal(t, byte(0x03), z80.reg.D)
	assert.Equal(t, byte(0x04), z80.reg.E)
	assert.Equal(t, 0, z80.t)
}

func Test_EX_SP_HL(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x12, 0x70, ld_sp_nn, 0x08, 0x00, ex_sp_hl, nop, 0x11, 0x22}}
	z80 := NewZ80(mem)
	z80.Run(10 + 10 + 19)

	assert.Equal(t, byte(0x22), z80.reg.H)
	assert.Equal(t, byte(0x11), z80.reg.L)
	assert.Equal(t, byte(0x12), z80.mem.Read(8))
	assert.Equal(t, byte(0x70), z80.mem.Read(9))
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x12, 0x70, ld_sp_nn, 0x0A, 0x00, prefix, ex_sp_hl, nop, 0x11, 0x22}}
		z80 := NewZ80(mem)
		z80.Run(14 + 10 + 23)

		switch prefix {
		case useIX:
			assert.Equal(t, byte(0x22), byte(z80.reg.IXH))
			assert.Equal(t, byte(0x11), byte(z80.reg.IXL))
		case useIY:
			assert.Equal(t, byte(0x22), byte(z80.reg.IYH))
			assert.Equal(t, byte(0x11), byte(z80.reg.IYL))
		}
		assert.Equal(t, byte(0x12), z80.mem.Read(10))
		assert.Equal(t, byte(0x70), z80.mem.Read(11))
		assert.Equal(t, 0, z80.t)
	}
}

func Test_ADD_A_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0, add_a_n, 0}}
	z80 := NewZ80(mem)
	z80.reg.F = fALL & ^fZ
	z80.Run(7 + 7)

	assert.Equal(t, byte(0), z80.reg.A)
	assert.Equal(t, fZ, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xFF, ld_h_n, 0x00, ld_l_n, 0x08, add_a_hl, nop, 0x01}}
	z80 = NewZ80(mem)
	z80.reg.F = fALL & ^fZ
	z80.Run(7 + 7 + 7 + 7)

	assert.Equal(t, byte(0), z80.reg.A)
	assert.Equal(t, fZ|fH|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x70, ld_l_n, 0x70, add_a_l}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(7 + 7 + 4)

	assert.Equal(t, byte(0xE0), z80.reg.A)
	assert.Equal(t, fS|fY|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xF0, add_a_n, 0xB0}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(7 + 7)

	assert.Equal(t, byte(0xA0), z80.reg.A)
	assert.Equal(t, fS|fY|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x8f, add_a_n, 0x81}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(7 + 7)

	assert.Equal(t, byte(0x10), z80.reg.A)
	assert.Equal(t, fH|fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x22, prefix, ld_hl_nn, 0x13, 0x00, prefix, add_a_l}}
		z80 = NewZ80(mem)
		z80.reg.F = fNONE
		z80.Run(7 + 14 + 8)

		assert.Equal(t, byte(0x35), z80.reg.A)
		assert.Equal(t, fY, z80.reg.F)
		assert.Equal(t, 0, z80.t)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x22, prefix, ld_hl_nn, 0x09, 0x00, prefix, add_a_hl, 0x01, nop, 0x13}}
		z80 = NewZ80(mem)
		z80.reg.F = fNONE
		z80.Run(7 + 14 + 19)

		assert.Equal(t, byte(0x35), z80.reg.A)
		assert.Equal(t, fY, z80.reg.F)
		assert.Equal(t, 0, z80.t)
	}
}

func Test_ADC_A_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x20, adc_a_n, 0x20}}
	z80 := NewZ80(mem)
	z80.reg.F = fC
	z80.Run(7 + 7)

	assert.Equal(t, byte(0x41), z80.reg.A)
	assert.Equal(t, fNONE, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0, ld_b_n, 0xFF, adc_a_b}}
	z80 = NewZ80(mem)
	z80.reg.F = fN | fC
	z80.Run(7 + 7 + 4)

	assert.Equal(t, byte(0), z80.reg.A)
	assert.Equal(t, fZ|fH|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x0F, ld_b_n, 0x00, adc_a_b}}
	z80 = NewZ80(mem)
	z80.reg.F = fC
	z80.Run(7 + 7 + 4)

	assert.Equal(t, byte(0x10), z80.reg.A)
	assert.Equal(t, fH, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x0F, ld_l_n, 0x06, adc_a_hl, nop, 0x70}}
	z80 = NewZ80(mem)
	z80.reg.F = fC
	z80.Run(7 + 7 + 7)

	assert.Equal(t, byte(0x80), z80.reg.A)
	assert.Equal(t, fS|fH|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x0F, ld_b_n, 0x69, adc_a_b}}
	z80 = NewZ80(mem)
	z80.reg.F = fC
	z80.Run(7 + 7 + 4)

	assert.Equal(t, byte(0x79), z80.reg.A)
	assert.Equal(t, fY|fH|fX, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x0E, ld_b_n, 0x01, adc_a_b}}
	z80 = NewZ80(mem)
	z80.reg.F = fC
	z80.Run(7 + 7 + 4)

	assert.Equal(t, byte(0x10), z80.reg.A)
	assert.Equal(t, fH, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x0E, ld_b_n, 0x01, adc_a_b}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(7 + 7 + 4)

	assert.Equal(t, byte(0x0F), z80.reg.A)
	assert.Equal(t, fX, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x22, prefix, ld_hl_nn, 0x13, 0x00, prefix, adc_a_l}}
		z80 = NewZ80(mem)
		z80.reg.F = fC
		z80.Run(7 + 14 + 8)

		assert.Equal(t, byte(0x36), z80.reg.A)
		assert.Equal(t, fY, z80.reg.F)
		assert.Equal(t, 0, z80.t)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x22, prefix, ld_hl_nn, 0x09, 0x00, prefix, adc_a_hl, 0x01, nop, 0x13}}
		z80 = NewZ80(mem)
		z80.reg.F = fC
		z80.Run(7 + 14 + 19)

		assert.Equal(t, byte(0x36), z80.reg.A)
		assert.Equal(t, fY, z80.reg.F)
		assert.Equal(t, 0, z80.t)
	}
}

func Test_ADD_HL_RR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0xFF, 0xFF, ld_bc_nn, 0x01, 0, add_hl_bc}}
	z80 := NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(10 + 10 + 11)

	assert.Equal(t, byte(0), z80.reg.H)
	assert.Equal(t, byte(0), z80.reg.L)
	assert.Equal(t, fH|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x41, 0x42, ld_de_nn, 0x11, 0x11, add_hl_de}}
	z80 = NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(10 + 10 + 11)

	assert.Equal(t, byte(0x53), z80.reg.H)
	assert.Equal(t, byte(0x52), z80.reg.L)
	assert.Equal(t, fS|fZ|fY|fX|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x41, 0x42, add_hl_hl}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(10 + 11)

	assert.Equal(t, byte(0x84), z80.reg.H)
	assert.Equal(t, byte(0x82), z80.reg.L)
	assert.Equal(t, fNONE, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0xFE, 0xFF, ld_sp_nn, 0x02, 0, add_hl_sp}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(10 + 10 + 11)

	assert.Equal(t, byte(0), z80.reg.H)
	assert.Equal(t, byte(0), z80.reg.L)
	assert.Equal(t, fH|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0xFE, 0xFF, ld_sp_nn, 0x03, 0, prefix, add_hl_sp}}
		z80 = NewZ80(mem)
		z80.reg.F = fNONE
		z80.Run(14 + 10 + 15)

		assert.Equal(t, byte(0x0), *z80.reg.prefixed[prefix][rH])
		assert.Equal(t, byte(0x01), *z80.reg.prefixed[prefix][rL])
		assert.Equal(t, fH|fC, z80.reg.F)
		assert.Equal(t, 0, z80.t)
	}
}

func Test_ADC_HL_RR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{scf, ld_hl_nn, 0xFE, 0xFF, ld_bc_nn, 0x01, 0x00, prefix_ed, adc_hl_bc}}
	z80 := NewZ80(mem)
	z80.Run(4 + 10 + 10 + 15)

	assert.Equal(t, byte(0), z80.reg.H)
	assert.Equal(t, byte(0), z80.reg.L)
	assert.Equal(t, fZ|fH|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{scf, ld_hl_nn, 0xC0, 0x63, ld_de_nn, 0xD0, 0x8A, prefix_ed, adc_hl_de}}
	z80 = NewZ80(mem)
	z80.Run(4 + 10 + 10 + 15)

	assert.Equal(t, byte(0xEE), z80.reg.H)
	assert.Equal(t, byte(0x91), z80.reg.L)
	assert.Equal(t, fS|fY|fX, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{scf, ld_hl_nn, 0x18, 0x7F, ld_de_nn, 0x48, 0x77, prefix_ed, adc_hl_de}}
	z80 = NewZ80(mem)
	z80.Run(4 + 10 + 10 + 15)

	assert.Equal(t, byte(0xF6), z80.reg.H)
	assert.Equal(t, byte(0x61), z80.reg.L)
	assert.Equal(t, fS|fY|fH|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_SUB_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0, sub_n, 0x01}}
	z80 := NewZ80(mem)
	z80.reg.F = fZ
	z80.Run(7 + 7)

	assert.Equal(t, byte(0xFF), z80.reg.A)
	assert.Equal(t, fS|fY|fH|fX|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x20, sub_a}}
	z80 = NewZ80(mem)
	z80.Run(7 + 4)

	assert.Equal(t, byte(0x00), z80.reg.A)
	assert.Equal(t, fZ|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x90, ld_h_n, 0x20, sub_h}}
	z80 = NewZ80(mem)
	z80.reg.F = fZ
	z80.Run(7 + 7 + 4)

	assert.Equal(t, byte(0x70), z80.reg.A)
	assert.Equal(t, fY|fP|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x7F, ld_l_n, 0x06, sub_hl, nop, 0x80}}
	z80 = NewZ80(mem)
	z80.Run(7 + 7 + 7)

	assert.Equal(t, byte(0xFF), z80.reg.A)
	assert.Equal(t, fS|fY|fX|fP|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x12, prefix, ld_h_n, 0x02, prefix, sub_h}}
		z80 = NewZ80(mem)
		z80.Run(7 + 11 + 8)

		assert.Equal(t, byte(0x10), z80.reg.A)
		assert.Equal(t, fN, z80.reg.F)
		assert.Equal(t, 0, z80.t)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x12, prefix, sub_hl, 0x06, nop, 0x01}}
		z80 = NewZ80(mem)
		z80.Run(7 + 19)

		assert.Equal(t, byte(0x11), z80.reg.A)
		assert.Equal(t, fN, z80.reg.F)
		assert.Equal(t, 0, z80.t)
	}
}

func Test_CP_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0, ld_b_n, 0x01, cp_b}}
	z80 := NewZ80(mem)
	z80.reg.F = fZ
	z80.Run(7 + 7 + 4)

	assert.Equal(t, byte(0), z80.reg.A)
	assert.Equal(t, fS|fH|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x20, cp_a}}
	z80 = NewZ80(mem)
	z80.Run(7 + 4)

	assert.Equal(t, byte(0x20), z80.reg.A)
	assert.Equal(t, fZ|fY|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x90, cp_n, 0x20}}
	z80 = NewZ80(mem)
	z80.reg.F = fZ
	z80.Run(7 + 7)

	assert.Equal(t, byte(0x90), z80.reg.A)
	assert.Equal(t, fY|fP|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x7F, ld_l_n, 0x06, cp_hl, nop, 0x80}}
	z80 = NewZ80(mem)
	z80.Run(7 + 7 + 7)

	assert.Equal(t, byte(0x7F), z80.reg.A)
	assert.Equal(t, fS|fP|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x7F, prefix, ld_h_n, 0x80, prefix, cp_h}}
		z80 = NewZ80(mem)
		z80.Run(7 + 11 + 8)

		assert.Equal(t, byte(0x7F), z80.reg.A)
		assert.Equal(t, fS|fP|fN|fC, z80.reg.F)
		assert.Equal(t, 0, z80.t)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x7F, prefix, cp_hl, 0x06, nop, 0x80}}
		z80 = NewZ80(mem)
		z80.Run(7 + 19)

		assert.Equal(t, byte(0x7F), z80.reg.A)
		assert.Equal(t, fS|fP|fN|fC, z80.reg.F)
		assert.Equal(t, 0, z80.t)
	}
}

func Test_SBC_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x01, ld_b_n, 0x01, sbc_a_b}}
	z80 := NewZ80(mem)
	z80.reg.F = fC
	z80.Run(7 + 7 + 4)

	assert.Equal(t, byte(0xFF), z80.reg.A)
	assert.Equal(t, fS|fY|fH|fX|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x7F, ld_l_n, 0x80, sbc_a_l}}
	z80 = NewZ80(mem)
	z80.reg.F = fC
	z80.Run(7 + 7 + 4)

	assert.Equal(t, byte(0xFE), z80.reg.A)
	assert.Equal(t, fS|fY|fX|fP|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x02, sbc_a_n, 0x01}}
	z80 = NewZ80(mem)
	z80.reg.F = fC
	z80.Run(7 + 7)

	assert.Equal(t, byte(0x00), z80.reg.A)
	assert.Equal(t, fZ|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x81, ld_l_n, 0x06, sbc_a_hl, nop, 0x01}}
	z80 = NewZ80(mem)
	z80.reg.F = fC
	z80.Run(7 + 7 + 7)

	assert.Equal(t, byte(0x7F), z80.reg.A)
	assert.Equal(t, fY|fH|fX|fP|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x12, prefix, ld_h_n, 0x02, prefix, sbc_a_h}}
		z80 = NewZ80(mem)
		z80.reg.F = fC
		z80.Run(7 + 11 + 8)

		assert.Equal(t, byte(0x0F), z80.reg.A)
		assert.Equal(t, fH|fX|fN, z80.reg.F)
		assert.Equal(t, 0, z80.t)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x12, prefix, sbc_a_hl, 0x06, nop, 0x01}}
		z80 = NewZ80(mem)
		z80.reg.F = fC
		z80.Run(7 + 19)

		assert.Equal(t, byte(0x10), z80.reg.A)
		assert.Equal(t, fN, z80.reg.F)
		assert.Equal(t, 0, z80.t)
	}
}

func Test_SBC_HL_RR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{scf, ld_hl_nn, 0xFE, 0xFF, ld_bc_nn, 0xFD, 0xFF, prefix_ed, sbc_hl_bc}}
	z80 := NewZ80(mem)
	z80.Run(4 + 10 + 10 + 15)

	assert.Equal(t, byte(0), z80.reg.H)
	assert.Equal(t, byte(0), z80.reg.L)
	assert.Equal(t, fZ|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{scf, ld_hl_nn, 0x01, 0x00, ld_bc_nn, 0xFD, 0x7F, prefix_ed, sbc_hl_bc}}
	z80 = NewZ80(mem)
	z80.Run(4 + 10 + 10 + 15)

	assert.Equal(t, byte(0x80), z80.reg.H)
	assert.Equal(t, byte(0x03), z80.reg.L)
	assert.Equal(t, fS|fH|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{scf, ld_hl_nn, 0x01, 0x70, ld_bc_nn, 0xFD, 0x8F, prefix_ed, sbc_hl_bc}}
	z80 = NewZ80(mem)
	z80.Run(4 + 10 + 10 + 15)

	assert.Equal(t, byte(0xE0), z80.reg.H)
	assert.Equal(t, byte(0x03), z80.reg.L)
	assert.Equal(t, fS|fY|fH|fP|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_NEG(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x55, prefix_ed, neg}}
	z80 := NewZ80(mem)
	z80.Run(7 + 8)

	assert.Equal(t, byte(0xAB), z80.reg.A)
	assert.Equal(t, fS|fY|fH|fX|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x00, prefix_ed, neg}}
	z80 = NewZ80(mem)
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x00), z80.reg.A)
	assert.Equal(t, fZ|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x80, prefix_ed, neg}}
	z80 = NewZ80(mem)
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x80), z80.reg.A)
	assert.Equal(t, fS|fP|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xAA, prefix_ed, neg}}
	z80 = NewZ80(mem)
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x56), z80.reg.A)
	assert.Equal(t, fH|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_AND_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x0F, ld_b_n, 0xF0, and_b}}
	z80 := NewZ80(mem)
	z80.reg.F = fN | fC
	z80.Run(7 + 7 + 4)

	assert.Equal(t, byte(0), z80.reg.A)
	assert.Equal(t, fZ|fH|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x8F, and_n, 0xF3}}
	z80 = NewZ80(mem)
	z80.Run(7 + 7)

	assert.Equal(t, byte(0x83), z80.reg.A)
	assert.Equal(t, fS|fH, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xFF, ld_l_n, 0x06, and_hl, nop, 0x81}}
	z80 = NewZ80(mem)
	z80.Run(7 + 7 + 7)

	assert.Equal(t, byte(0x81), z80.reg.A)
	assert.Equal(t, fS|fH|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x01, prefix, ld_l_n, 0x03, prefix, and_l}}
		z80 = NewZ80(mem)
		z80.Run(7 + 11 + 8)

		assert.Equal(t, byte(0x01), z80.reg.A)
		assert.Equal(t, fH, z80.reg.F)
		assert.Equal(t, 0, z80.t)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x88, prefix, and_hl, 0x06, nop, 0x08}}
		z80 = NewZ80(mem)
		z80.Run(7 + 19)

		assert.Equal(t, byte(0x08), z80.reg.A)
		assert.Equal(t, fH|fX, z80.reg.F)
		assert.Equal(t, 0, z80.t)
	}
}

func Test_OR_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x00, ld_b_n, 0x00, or_b}}
	z80 := NewZ80(mem)
	z80.reg.F = fS | fH | fN | fC
	z80.Run(7 + 7 + 4)

	assert.Equal(t, byte(0), z80.reg.A)
	assert.Equal(t, fZ|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x8A, ld_l_n, 0x06, or_hl, nop, 0x85}}
	z80 = NewZ80(mem)
	z80.reg.F = fS | fH | fN | fC
	z80.Run(7 + 7 + 7)

	assert.Equal(t, byte(0x8F), z80.reg.A)
	assert.Equal(t, fS|fX, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x11, or_n, 0x20}}
	z80 = NewZ80(mem)
	z80.reg.F = fS | fH | fN | fC
	z80.Run(7 + 7)

	assert.Equal(t, byte(0x31), z80.reg.A)
	assert.Equal(t, fY, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x01, prefix, ld_l_n, 0x12, prefix, or_l}}
		z80 = NewZ80(mem)
		z80.Run(7 + 11 + 8)

		assert.Equal(t, byte(0x13), z80.reg.A)
		assert.Equal(t, fNONE, z80.reg.F)
		assert.Equal(t, 0, z80.t)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x80, prefix, or_hl, 0x06, nop, 0x08}}
		z80 = NewZ80(mem)
		z80.Run(7 + 19)

		assert.Equal(t, byte(0x88), z80.reg.A)
		assert.Equal(t, fS|fX|fP, z80.reg.F)
		assert.Equal(t, 0, z80.t)
	}
}

func Test_XOR_x(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x1F, ld_b_n, 0x1F, xor_b}}
	z80 := NewZ80(mem)
	z80.reg.F = fH | fN | fC
	z80.Run(7 + 7 + 4)

	assert.Equal(t, byte(0), z80.reg.A)
	assert.Equal(t, fZ|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x1F, ld_l_n, 0x06, xor_hl, nop, 0x8F}}
	z80 = NewZ80(mem)
	z80.Run(7 + 7 + 7)

	assert.Equal(t, byte(0x90), z80.reg.A)
	assert.Equal(t, fS|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x1F, xor_n, 0x0F}}
	z80 = NewZ80(mem)
	z80.Run(7 + 7)

	assert.Equal(t, byte(0x10), z80.reg.A)
	assert.Equal(t, fNONE, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x01, prefix, ld_l_n, 0x03, prefix, xor_l}}
		z80 = NewZ80(mem)
		z80.Run(7 + 11 + 8)

		assert.Equal(t, byte(0x02), z80.reg.A)
		assert.Equal(t, fNONE, z80.reg.F)
		assert.Equal(t, 0, z80.t)

		mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x88, prefix, xor_hl, 0x06, nop, 0x08}}
		z80 = NewZ80(mem)
		z80.Run(7 + 19)

		assert.Equal(t, byte(0x80), z80.reg.A)
		assert.Equal(t, fS, z80.reg.F)
		assert.Equal(t, 0, z80.t)
	}
}

func Test_INC_R(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0, inc_a}}
	z80 := NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(7 + 4)

	assert.Equal(t, fC, z80.reg.F)
	assert.Equal(t, byte(0x01), z80.reg.A)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	z80.reg.F = fALL & ^fZ
	mem.Cells[1] = 0xFF
	z80.Run(7 + 4)

	assert.Equal(t, fZ|fH|fC, z80.reg.F)
	assert.Equal(t, byte(0x00), z80.reg.A)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	z80.reg.F = fN
	mem.Cells[1] = 0x7F
	z80.Run(7 + 4)

	assert.Equal(t, fS|fH|fP, z80.reg.F)
	assert.Equal(t, byte(0x80), z80.reg.A)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	mem.Cells[1] = 0x92
	z80.reg.F = fNONE
	z80.Run(7 + 4)

	assert.Equal(t, fS, z80.reg.F)
	assert.Equal(t, byte(0x93), z80.reg.A)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	mem.Cells[1] = 0x10
	z80.reg.F = fNONE
	z80.Run(7 + 4)

	assert.Equal(t, fNONE, z80.reg.F)
	assert.Equal(t, byte(0x11), z80.reg.A)
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{
			prefix, ld_h_n, 0x10, prefix, inc_h,
			prefix, ld_l_n, 0x20, prefix, inc_l},
		}
		z80 = NewZ80(mem)
		z80.reg.F = fNONE
		z80.Run(11 + 8 + 11 + 8)

		assert.Equal(t, fY, z80.reg.F)
		assert.Equal(t, byte(0x11), *z80.reg.prefixed[prefix][rH])
		assert.Equal(t, byte(0x21), *z80.reg.prefixed[prefix][rL])
		assert.Equal(t, 0, z80.t)
	}
}

func Test_INC_RR(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{
			ld_bc_nn, 0x34, 0x12, inc_bc, ld_de_nn, 0x35, 0x13, inc_de,
			ld_hl_nn, 0x36, 0x14, inc_hl, ld_sp_nn, 0x37, 0x15, inc_sp,
			useIX, ld_hl_nn, 0x38, 0x16, useIX, inc_hl,
			useIY, ld_hl_nn, 0x39, 0x17, useIY, inc_hl}}
	z80 := NewZ80(mem)
	z80.Run(10 + 6 + 10 + 6 + 10 + 6 + 10 + 6 + 14 + 10 + 14 + 10)

	assert.Equal(t, uint16(0x1235), z80.reg.BC())
	assert.Equal(t, uint16(0x1336), z80.reg.DE())
	assert.Equal(t, uint16(0x1437), z80.reg.HL())
	assert.Equal(t, uint16(0x1538), z80.reg.SP)
	assert.Equal(t, uint16(0x1639), uint16(z80.reg.IXH)<<8|uint16(z80.reg.IXL))
	assert.Equal(t, uint16(0x173A), uint16(z80.reg.IYH)<<8|uint16(z80.reg.IYL))
	assert.Equal(t, 0, z80.t)
}

func Test_INC_mHL(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x05, 0x00, inc_mhl, nop, 0xFF}}
	z80 := NewZ80(mem)
	z80.reg.F = fS | fP | fN | fC
	z80.Run(10 + 11)

	assert.Equal(t, byte(0x00), mem.Cells[5])
	assert.Equal(t, fZ|fH|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	mem.Cells[5] = 0x7F
	z80.reg.F = fNONE
	z80.Run(10 + 11)

	assert.Equal(t, byte(0x80), mem.Cells[5])
	assert.Equal(t, fS|fH|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	mem.Cells[5] = 0x20
	z80.reg.F = fNONE
	z80.Run(10 + 11)

	assert.Equal(t, byte(0x21), mem.Cells[5])
	assert.Equal(t, fY, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x05, 0x00, prefix, inc_mhl, 0x03, nop, 0x3F}}
		z80 := NewZ80(mem)
		z80.reg.F = fNONE
		z80.Run(14 + 23)

		assert.Equal(t, byte(0x40), mem.Cells[8])
		assert.Equal(t, fH, z80.reg.F)
		assert.Equal(t, 0, z80.t)
	}
}

func Test_DEC_R(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 1, dec_a}}
	z80 := NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(7 + 4)

	assert.Equal(t, fZ|fN, z80.reg.F)
	assert.Equal(t, byte(0x00), z80.reg.A)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	z80.reg.F = fALL & ^(fZ | fH | fN)
	mem.Cells[1] = 0
	z80.Run(7 + 4)

	assert.Equal(t, fS|fH|fY|fX|fN|fC, z80.reg.F)
	assert.Equal(t, byte(0xFF), z80.reg.A)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	z80.reg.F = fZ | fS
	mem.Cells[1] = 0x80
	z80.Run(7 + 4)

	assert.Equal(t, fY|fH|fX|fP|fN, z80.reg.F)
	assert.Equal(t, byte(0x7F), z80.reg.A)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	z80.reg.F = fALL
	mem.Cells[1] = 0xAB
	z80.Run(7 + 4)

	assert.Equal(t, fS|fY|fX|fN|fC, z80.reg.F)
	assert.Equal(t, byte(0xAA), z80.reg.A)
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{
			prefix, ld_h_n, 0x10, prefix, dec_h,
			prefix, ld_l_n, 0x20, prefix, dec_l},
		}
		z80 = NewZ80(mem)
		z80.Run(11 + 8 + 11 + 8)

		assert.Equal(t, byte(0x0F), *z80.reg.prefixed[prefix][rH])
		assert.Equal(t, byte(0x1F), *z80.reg.prefixed[prefix][rL])
		assert.Equal(t, 0, z80.t)
	}
}

func Test_DEC_RR(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{
			ld_bc_nn, 0x34, 0x12, dec_bc, ld_de_nn, 0x35, 0x13, dec_de,
			ld_hl_nn, 0x36, 0x14, dec_hl, ld_sp_nn, 0x37, 0x15, dec_sp,
			useIX, ld_hl_nn, 0x38, 0x16, useIX, dec_hl,
			useIY, ld_hl_nn, 0x39, 0x17, useIY, dec_hl}}
	z80 := NewZ80(mem)
	z80.Run(10 + 6 + 10 + 6 + 10 + 6 + 10 + 6 + 14 + 10 + 14 + 10)

	assert.Equal(t, uint16(0x1233), z80.reg.BC())
	assert.Equal(t, uint16(0x1334), z80.reg.DE())
	assert.Equal(t, uint16(0x1435), z80.reg.HL())
	assert.Equal(t, uint16(0x1536), z80.reg.SP)
	assert.Equal(t, uint16(0x1637), uint16(z80.reg.IXH)<<8|uint16(z80.reg.IXL))
	assert.Equal(t, uint16(0x1738), uint16(z80.reg.IYH)<<8|uint16(z80.reg.IYL))
	assert.Equal(t, 0, z80.t)
}

func Test_DEC_mHL(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x05, 0x00, dec_mhl, nop, 0x00}}
	z80 := NewZ80(mem)
	z80.reg.F = fC
	z80.Run(10 + 11)

	assert.Equal(t, byte(0xFF), mem.Cells[5])
	assert.Equal(t, fS|fY|fH|fX|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	mem.Cells[5] = 0x01
	z80.reg.F = fNONE
	z80.Run(10 + 11)

	assert.Equal(t, byte(0x00), mem.Cells[5])
	assert.Equal(t, fZ|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	mem.Cells[5] = 0x80
	z80.reg.F = fNONE
	z80.Run(10 + 11)

	assert.Equal(t, byte(0x7F), mem.Cells[5])
	assert.Equal(t, fP|fY|fH|fX|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x05, 0x00, prefix, dec_mhl, 0x03, nop, 0x3F}}
		z80 := NewZ80(mem)
		z80.reg.F = fNONE
		z80.Run(14 + 23)

		assert.Equal(t, byte(0x3E), mem.Cells[8])
		assert.Equal(t, fY|fX|fN, z80.reg.F)
		assert.Equal(t, 0, z80.t)
	}
}

func Test_LD_RR_nn(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_bc_nn, 0x34, 0x12}}
	z80 := NewZ80(mem)
	z80.Run(10)

	assert.Equal(t, uint16(0x1234), z80.reg.BC())
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem = &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x34, 0x12}}
		z80 = NewZ80(mem)
		z80.Run(14)

		assert.Equal(t, byte(0x12), *z80.reg.prefixed[prefix][rH])
		assert.Equal(t, byte(0x34), *z80.reg.prefixed[prefix][rL])
		assert.Equal(t, 0, z80.t)
	}
}

func Test_LD_mm_HL(t *testing.T) {
	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x3A, 0x48, prefix, ld_mm_hl, 0x09, 0x00, nop, 0, 0}}
		z80 := NewZ80(mem)
		z80.Run(14 + 20)

		assert.Equal(t, *z80.reg.prefixed[prefix][rH], mem.Cells[10])
		assert.Equal(t, *z80.reg.prefixed[prefix][rL], mem.Cells[9])
		assert.Equal(t, 0, z80.t)
	}
}

func Test_LD_mm_RR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_sp_nn, 0x3A, 0x48, prefix_ed, ld_mm_sp, 0x08, 0x00, nop, 0, 0}}
	z80 := NewZ80(mem)
	z80.Run(10 + 20)

	assert.Equal(t, byte(0x3A), mem.Cells[8])
	assert.Equal(t, byte(0x48), mem.Cells[9])
	assert.Equal(t, 0, z80.t)
}

func Test_LD_HL_mm(t *testing.T) {
	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_mm, 0x05, 0x00, nop, 0x34, 0x12}}
		z80 := NewZ80(mem)
		z80.Run(20)

		assert.Equal(t, byte(0x12), *z80.reg.prefixed[prefix][rH])
		assert.Equal(t, byte(0x34), *z80.reg.prefixed[prefix][rL])
		assert.Equal(t, 0, z80.t)
	}
}

func Test_LD_RR_mm(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{prefix_ed, ld_de_mm, 0x05, 0x00, nop, 0x34, 0x12}}
	z80 := NewZ80(mem)
	z80.Run(20)

	assert.Equal(t, byte(0x12), z80.reg.D)
	assert.Equal(t, byte(0x34), z80.reg.E)
	assert.Equal(t, 0, z80.t)
}

func Test_LD_mHL_n(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, ld_mhl_n, 0xAB, nop, 0x00}}
	z80 := NewZ80(mem)
	z80.Run(10 + 10)

	assert.Equal(t, byte(0xAB), z80.mem.Read(6))
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x06, 0x00, prefix, ld_mhl_n, 0x03, 0xAB, nop, 0x00}}
		z80 := NewZ80(mem)
		z80.Run(14 + 19)

		assert.Equal(t, byte(0xAB), z80.mem.Read(9))
		assert.Equal(t, 0, z80.t)
	}
}

func Test_LD_SP_HL(t *testing.T) {
	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x20, 0x30, prefix, ld_sp_hl}}
		z80 := NewZ80(mem)
		z80.Run(14 + 10)

		assert.Equal(t, uint16(0x3020), z80.reg.SP)
		assert.Equal(t, 0, z80.t)
	}
}

func Test_LD_mIXY_n(t *testing.T) {
	for _, prefix := range []byte{useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x02, 0x00, prefix, ld_mhl_n, 0x07, 0xAB, nop, 0x00}}
		z80 := NewZ80(mem)
		z80.Run(14 + 19)

		assert.Equal(t, byte(0xAB), z80.mem.Read(9))
		assert.Equal(t, 0, z80.t)
	}
}

func Test_LD_mm_A(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x9F, ld_mm_a, 0x06, 0x00, nop, 0x00}}
	z80 := NewZ80(mem)
	z80.Run(7 + 13)

	assert.Equal(t, z80.reg.A, mem.Cells[6])
	assert.Equal(t, 0, z80.t)
}

func Test_LD_A_mm(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_mm, 0x04, 0x00, nop, 0xDE}}
	z80 := NewZ80(mem)
	z80.Run(13)

	assert.Equal(t, byte(0xDE), z80.reg.A)
	assert.Equal(t, 0, z80.t)
}

func Test_LD_BC_A(t *testing.T) {
	var n byte = 0x76
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, n, ld_bc_nn, 0x07, 0x00, ld_bc_a, nop, 0x00}}
	z80 := NewZ80(mem)
	z80.Run(7 + 10 + 7)

	assert.Equal(t, n, z80.mem.Read(7))
	assert.Equal(t, 0, z80.t)
}

func Test_LD_DE_A(t *testing.T) {
	var n byte = 0x76
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, n, ld_de_nn, 0x07, 0x00, ld_de_a, nop, 0x00}}
	z80 := NewZ80(mem)
	z80.Run(7 + 10 + 7)

	assert.Equal(t, n, z80.mem.Read(7))
	assert.Equal(t, 0, z80.t)
}

func Test_LD_A_BC(t *testing.T) {
	var n byte = 0x76
	mem := &memory.BasicMemory{Cells: []byte{ld_bc_nn, 0x05, 0x00, ld_a_bc, nop, n}}
	z80 := NewZ80(mem)
	z80.Run(10 + 7)

	assert.Equal(t, n, z80.reg.A)
	assert.Equal(t, 0, z80.t)
}

func Test_LD_A_DE(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_de_nn, 0x05, 0x00, ld_a_de, nop, 0x76}}
	z80 := NewZ80(mem)
	z80.Run(10 + 7)

	assert.Equal(t, byte(0x76), z80.reg.A)
	assert.Equal(t, 0, z80.t)
}

func Test_LD_A_R(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{prefix_ed, ld_a_r}}
	z80 := NewZ80(mem)
	z80.reg.F = fH | fN | fC
	z80.Run(9)

	assert.Equal(t, byte(0x02), z80.reg.A)
	assert.Equal(t, fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x7E, prefix_ed, ld_r_a, prefix_ed, ld_a_r}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.iff2 = true
	z80.Run(7 + 9 + 9)

	assert.Equal(t, byte(0x00), z80.reg.A)
	assert.Equal(t, fZ|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_LD_A_I(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x05, prefix_ed, ld_a_i}}
	z80 := NewZ80(mem)
	z80.reg.F = fH | fN | fC
	z80.iff2 = true
	z80.Run(7 + 9)

	assert.Equal(t, byte(0), z80.reg.A)
	assert.Equal(t, fZ|fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_LD_R_A(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x85, prefix_ed, ld_r_a}}
	z80 := NewZ80(mem)
	z80.Run(7 + 9)

	assert.Equal(t, byte(0x85), z80.reg.R)
	assert.Equal(t, 0, z80.t)
}

func Test_LD_I_A(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x85, prefix_ed, ld_i_a}}
	z80 := NewZ80(mem)
	z80.Run(7 + 9)

	assert.Equal(t, byte(0x85), z80.reg.I)
	assert.Equal(t, 0, z80.t)
}

func Test_LD_R_n(t *testing.T) {
	var a, b, c, d, e, h, l, ixh, ixl, iyh, iyl byte = 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	mem := &memory.BasicMemory{
		Cells: []byte{
			ld_a_n, a, ld_b_n, b, ld_c_n, c, ld_d_n, d, ld_e_n, e, ld_h_n, h, ld_l_n, l,
			useIX, ld_h_n, ixh, useIX, ld_l_n, ixl, useIY, ld_h_n, iyh, useIY, ld_l_n, iyl}}
	z80 := NewZ80(mem)
	z80.Run(7 + 7 + 7 + 7 + 7 + 7 + 7 + 11 + 11 + 11 + 11)

	assert.Equal(t, a, z80.reg.A)
	assert.Equal(t, b, z80.reg.B)
	assert.Equal(t, c, z80.reg.C)
	assert.Equal(t, d, z80.reg.D)
	assert.Equal(t, e, z80.reg.E)
	assert.Equal(t, h, z80.reg.H)
	assert.Equal(t, l, z80.reg.L)
	assert.Equal(t, ixh, z80.reg.IXH)
	assert.Equal(t, ixl, z80.reg.IXL)
	assert.Equal(t, iyh, z80.reg.IYH)
	assert.Equal(t, iyl, z80.reg.IYL)
	assert.Equal(t, 0, z80.t)
}

func Test_LD_R_R(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{ld_a_n, 0x56, ld_b_a, ld_c_b, ld_d_c, ld_e_d, ld_h_e, ld_l_h, ld_a_n, 0,
			useIX, ld_h_b, useIX, ld_l_h, useIY, ld_l_b, useIY, ld_l_e, useIY, ld_a_l}}
	z80 := NewZ80(mem)
	z80.Run(7 + 4 + 4 + 4 + 4 + 4 + 4 + 7 + 8 + 8 + 8 + 8 + 8)

	assert.Equal(t, byte(0x56), z80.reg.A)
	assert.Equal(t, byte(0x56), z80.reg.B)
	assert.Equal(t, byte(0x56), z80.reg.C)
	assert.Equal(t, byte(0x56), z80.reg.D)
	assert.Equal(t, byte(0x56), z80.reg.E)
	assert.Equal(t, byte(0x56), z80.reg.H)
	assert.Equal(t, byte(0x56), z80.reg.L)
	assert.Equal(t, byte(0x56), z80.reg.IXH)
	assert.Equal(t, byte(0x56), z80.reg.IXL)
	assert.Equal(t, 0, z80.t)
}

func Test_LD_R_HL(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, ld_a_hl, ld_l_hl, nop, 0xA7}}
	z80 := NewZ80(mem)
	z80.Run(10 + 7 + 7)

	assert.Equal(t, byte(0xA7), z80.reg.A)
	assert.Equal(t, byte(0xA7), z80.reg.L)
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x01, 0x00, prefix, ld_l_hl, 0x07, nop, 0xA7}}
		z80 := NewZ80(mem)
		z80.Run(14 + 19)

		assert.Equal(t, byte(0xA7), z80.reg.L)
		assert.Equal(t, 0, z80.t)
	}
}

func Test_LD_HL_R(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_d_n, 0x99, ld_hl_nn, 0x07, 0x00, ld_hl_d, nop, 0x00}}
	z80 := NewZ80(mem)
	z80.Run(7 + 10 + 7)

	assert.Equal(t, byte(0x99), z80.mem.Read(7))
	assert.Equal(t, 0, z80.t)

	for _, prefix := range []byte{useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{ld_d_n, 0x99, prefix, ld_hl_nn, 0x07, 0x00, prefix, ld_hl_d, 0x03, nop, 0x00}}
		z80 := NewZ80(mem)
		z80.Run(7 + 14 + 19)

		assert.Equal(t, byte(0x99), z80.mem.Read(10))
		assert.Equal(t, 0, z80.t)
	}
}

func Test_CPL(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x5B, cpl}}
	z80 := NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(7 + 4)

	assert.Equal(t, byte(0xA4), z80.reg.A)
	assert.Equal(t, fH|fY|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_SCF(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{scf}}
	z80 := NewZ80(mem)
	z80.reg.F = fS | fZ | fH | fP | fN
	z80.Run(4)

	assert.Equal(t, fS|fZ|fY|fX|fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_CCF(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ccf}}
	z80 := NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(4)

	assert.Equal(t, fS|fZ|fY|fH|fX|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	z80.reg.F = fZ | fN
	z80.Run(4)

	assert.Equal(t, fZ|fY|fX|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	z80.reg.F = fZ | fN | fC
	z80.Run(4)

	assert.Equal(t, fZ|fY|fH|fX, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_RLCA(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x55, rlca}}
	z80 := NewZ80(mem)
	z80.reg.F = fH | fN
	z80.Run(7 + 4)

	assert.Equal(t, byte(0xAA), z80.reg.A)
	assert.Equal(t, fY|fX, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	mem.Cells[1] = 0xAA
	z80.reg.F = fH | fN
	z80.Run(7 + 4)

	assert.Equal(t, byte(0x55), z80.reg.A)
	assert.Equal(t, fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	mem.Cells[1] = 0x00
	z80.reg.F = fH | fN
	z80.Run(7 + 4)

	assert.Equal(t, byte(0x00), z80.reg.A)
	assert.Equal(t, fNONE, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	mem.Cells[1] = 0xFF
	z80.reg.F = fH | fN
	z80.Run(7 + 4)

	assert.Equal(t, byte(0xFF), z80.reg.A)
	assert.Equal(t, fY|fX|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_RRCA(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x55, rrca}}
	z80 := NewZ80(mem)
	z80.reg.F = fH | fN
	z80.Run(7 + 4)

	assert.Equal(t, byte(0xAA), z80.reg.A)
	assert.Equal(t, fY|fX|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	mem.Cells[1] = 0xAA
	z80.reg.F = fH | fN
	z80.Run(7 + 4)

	assert.Equal(t, byte(0x55), z80.reg.A)
	assert.Equal(t, fNONE, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	mem.Cells[1] = 0x00
	z80.reg.F = fH | fN
	z80.Run(7 + 4)

	assert.Equal(t, byte(0x00), z80.reg.A)
	assert.Equal(t, fNONE, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	mem.Cells[1] = 0xFF
	z80.reg.F = fH | fN
	z80.Run(7 + 4)

	assert.Equal(t, byte(0xFF), z80.reg.A)
	assert.Equal(t, fY|fX|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_RLA(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x80, rla}}
	z80 := NewZ80(mem)
	z80.reg.F = fH | fN | fC
	z80.Run(7 + 4)

	assert.Equal(t, byte(0x01), z80.reg.A)
	assert.Equal(t, fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x55, rla}}
	z80 = NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(7 + 4)

	assert.Equal(t, byte(0xab), z80.reg.A)
	assert.Equal(t, fS|fZ|fY|fX|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x88, rla, ld_b_a, rla}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(7 + 4 + 4 + 4)

	assert.Equal(t, byte(0x10), z80.reg.B)
	assert.Equal(t, byte(0x21), z80.reg.A)
	assert.Equal(t, fY, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_RRA(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x80, rra}}
	z80 := NewZ80(mem)
	z80.reg.F = fH | fN | fC
	z80.Run(7 + 4)

	assert.Equal(t, byte(0xC0), z80.reg.A)
	assert.Equal(t, fNONE, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x55, rra}}
	z80 = NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(7 + 4)

	assert.Equal(t, byte(0xAA), z80.reg.A)
	assert.Equal(t, fS|fZ|fY|fX|fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x89, rra, ld_b_a, rra}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(7 + 4 + 4 + 4)

	assert.Equal(t, byte(0x44), z80.reg.B)
	assert.Equal(t, byte(0xA2), z80.reg.A)
	assert.Equal(t, fY, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_RLD(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x7A, ld_hl_nn, 0x08, 0x00, prefix_ed, rld, nop, 0x31}}
	z80 := NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(7 + 10 + 18)

	assert.Equal(t, byte(0x73), z80.reg.A)
	assert.Equal(t, byte(0x1A), z80.mem.Read(8))
	assert.Equal(t, fY|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x0F, ld_hl_nn, 0x08, 0x00, prefix_ed, rld, nop, 0x0A}}
	z80 = NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(7 + 10 + 18)

	assert.Equal(t, byte(0), z80.reg.A)
	assert.Equal(t, byte(0xAF), z80.mem.Read(8))
	assert.Equal(t, fZ|fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_RRD(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x84, ld_hl_nn, 0x08, 0x00, prefix_ed, rrd, nop, 0x20}}
	z80 := NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(7 + 10 + 18)

	assert.Equal(t, byte(0x80), z80.reg.A)
	assert.Equal(t, byte(0x42), z80.mem.Read(8))
	assert.Equal(t, fS|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x03, ld_hl_nn, 0x08, 0x00, prefix_ed, rrd, nop, 0x60}}
	z80 = NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(7 + 10 + 18)

	assert.Equal(t, byte(0), z80.reg.A)
	assert.Equal(t, byte(0x36), z80.mem.Read(8))
	assert.Equal(t, fZ|fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_DAA(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x9A, daa}}
	z80 := NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(7 + 4)

	assert.Equal(t, byte(0), z80.reg.A)
	assert.Equal(t, fZ|fH|fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x99, daa}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(7 + 4)

	assert.Equal(t, byte(0x99), z80.reg.A)
	assert.Equal(t, fS|fX|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x8F, daa}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(7 + 4)

	assert.Equal(t, byte(0x95), z80.reg.A)
	assert.Equal(t, fS|fH|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x8F, daa}}
	z80 = NewZ80(mem)
	z80.reg.F = fN
	z80.Run(7 + 4)

	assert.Equal(t, byte(0x89), z80.reg.A)
	assert.Equal(t, fS|fX|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xCA, daa}}
	z80 = NewZ80(mem)
	z80.reg.F = fN
	z80.Run(7 + 4)

	assert.Equal(t, byte(0x64), z80.reg.A)
	assert.Equal(t, fY|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xC5, daa}}
	z80 = NewZ80(mem)
	z80.reg.F = fH | fN
	z80.Run(7 + 4)

	assert.Equal(t, byte(0x5F), z80.reg.A)
	assert.Equal(t, fH|fX|fP|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xCA, daa}}
	z80 = NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(7 + 4)

	assert.Equal(t, byte(0x64), z80.reg.A)
	assert.Equal(t, fY|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_DJNZ(t *testing.T) {
	var b byte = 0x20
	var o int8 = -3
	mem := &memory.BasicMemory{Cells: []byte{ld_b_n, b, ld_a_n, 0, inc_a, djnz, byte(o)}}
	z80 := NewZ80(mem)
	z80.Run(7 + 7 + 4 + 0x1F*(4+13) + 8)

	assert.Equal(t, b, z80.reg.A)
	assert.Equal(t, byte(0), z80.reg.B)
	assert.Equal(t, 0, z80.t)

	b = 0xFF
	o = 1
	mem = &memory.BasicMemory{Cells: []byte{ld_b_n, b, ld_a_n, 1, djnz, byte(o), halt, inc_a, jr_o, 0xFA}}
	z80 = NewZ80(mem)
	z80.Run(7 + 7 + 0xFE*(13+4+12) + 8)

	assert.Equal(t, b, z80.reg.A)
	assert.Equal(t, byte(0), z80.reg.B)
	assert.Equal(t, 0, z80.t)
}

func Test_JR_o(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{jr_o, 3, ld_c_n, 0x11, halt, ld_d_n, 0x22}}
	z80 := NewZ80(mem)
	z80.Run(12 + 7)

	assert.Equal(t, byte(0), z80.reg.C)
	assert.Equal(t, byte(0x22), z80.reg.D)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{jr_o, 6, halt, ld_c_n, 0x11, ld_b_n, 0x33, nop, jr_o, 0xF9}}
	z80 = NewZ80(mem)
	z80.Run(12 + 12 + 7 + 7)

	assert.Equal(t, byte(0x33), z80.reg.B)
	assert.Equal(t, byte(0x11), z80.reg.C)
	assert.Equal(t, 0, z80.t)
}

func Test_JR_Z_o(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 2, dec_a, jr_z_o, 0x02, ld_b_n, 0xab}}
	z80 := NewZ80(mem)
	z80.Run(7 + 4 + 7 + 7)

	assert.Equal(t, byte(0xab), z80.reg.B)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 1, dec_a, jr_z_o, 0x02, ld_b_n, 0xab}}
	z80 = NewZ80(mem)
	z80.Run(7 + 4 + 12)

	assert.Equal(t, byte(0), z80.reg.B)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 1, dec_a, jr_z_o, 0xFD}}
	z80 = NewZ80(mem)
	z80.Run(7 + 4 + 12 + 4 + 7)

	assert.Equal(t, byte(0xFF), z80.reg.A)
	assert.Equal(t, 0, z80.t)
}

func Test_JR_NZ_o(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 2, dec_a, jr_nz_o, 0x02, ld_b_n, 0xab}}
	z80 := NewZ80(mem)
	z80.Run(7 + 4 + 12)

	assert.Equal(t, byte(0), z80.reg.B)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 1, dec_a, jr_nz_o, 0x02, ld_b_n, 0xab}}
	z80 = NewZ80(mem)
	z80.Run(7 + 4 + 7 + 7)

	assert.Equal(t, byte(0xab), z80.reg.B)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 2, dec_a, jr_nz_o, 0xFD}}
	z80 = NewZ80(mem)
	z80.Run(7 + 4 + 12 + 4 + 7)

	assert.Equal(t, byte(0), z80.reg.A)
	assert.Equal(t, 0, z80.t)
}

func Test_JR_C(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0xFF, ld_b_n, 0xAB, dec_a, add_a_n, 1, jr_c, 2, ld_b_a}}
	z80 := NewZ80(mem)
	z80.Run(7 + 7 + 4 + 7 + 7 + 4)

	assert.Equal(t, byte(0xFF), z80.reg.B)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_b_n, 0xAB, scf, jr_c, 1, halt, ld_b_a}}
	z80 = NewZ80(mem)
	z80.Run(7 + 4 + 12 + 4)

	assert.Equal(t, byte(0xFF), z80.reg.B)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{xor_a, dec_a, add_a_n, 1, jr_c, 0xFC}}
	z80 = NewZ80(mem)
	z80.Run(4 + 4 + 7 + 12 + 7 + 7)

	assert.Equal(t, byte(1), z80.reg.A)
	assert.Equal(t, 0, z80.t)
}

func Test_JR_NC_o(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0xFF, ld_b_n, 0xAB, inc_a, add_a_n, 1, jr_nc_o, 1, ld_b_a}}
	z80 := NewZ80(mem)
	z80.Run(7 + 7 + 4 + 7 + 12)

	assert.Equal(t, byte(0xAB), z80.reg.B)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_b_n, 0xAB, scf, ccf, jr_nc_o, 1, halt, ld_b_a}}
	z80 = NewZ80(mem)
	z80.Run(7 + 4 + 4 + 12 + 4)

	assert.Equal(t, byte(0xFF), z80.reg.B)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0xFE, add_a_n, 1, jr_nc_o, 0xFC}}
	z80 = NewZ80(mem)
	z80.Run(7 + 7 + 12 + 7 + 7)

	assert.Equal(t, byte(0), z80.reg.A)
	assert.Equal(t, 0, z80.t)
}

func Test_JP_nn(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{jp_nn, 0x06, 0x00, ld_a_n, 0xAA, halt, ld_a_n, 0x55}}
	z80 := NewZ80(mem)
	z80.Run(10 + 7)

	assert.Equal(t, byte(0x55), z80.reg.A)
	assert.Equal(t, 0, z80.t)
}

func Test_JP_HL(t *testing.T) {
	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		mem := &memory.BasicMemory{Cells: []byte{prefix, ld_hl_nn, 0x09, 0x00, prefix, jp_hl, ld_a_n, 0xAA, halt, ld_a_n, 0x55}}
		z80 := NewZ80(mem)
		z80.Run(14 + 8 + 7)

		assert.Equal(t, byte(0x55), z80.reg.A)
		assert.Equal(t, 0, z80.t)
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
			Cells: []byte{test.jp, 0x06, 0x00, ld_a_n, 0xAA, halt, ld_a_n, 0x55},
		}

		z80 := NewZ80(mem)
		z80.reg.F = test.flag
		z80.Run(10 + 7)

		assert.Equal(t, byte(test.expected), z80.reg.A)
		assert.Equal(t, 0, z80.t)
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

		z80 := NewZ80(mem)
		z80.reg.F = test.flag
		if test.expected == 0x55 {
			z80.Run(10 + 17 + 7)
		} else {
			z80.Run(10 + 10 + 7)
		}

		assert.Equal(t, byte(test.expected), z80.reg.A)
		if z80.reg.A == 0x55 {
			assert.Equal(t, uint16(0x0E), z80.reg.SP)
			assert.Equal(t, byte(0), mem.Cells[15])
			assert.Equal(t, byte(0x06), mem.Cells[14])
		} else {
			assert.Equal(t, uint16(0x10), z80.reg.SP)
		}
		assert.Equal(t, 0, z80.t)
	}
}

func Test_RET(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_sp_nn, 0x0A, 0x00, ret, ld_a_n, 0xAA, halt, ld_a_n, 0x55, nop, 0x07, 0x00}}
	z80 := NewZ80(mem)
	z80.Run(10 + 10 + 7)

	assert.Equal(t, byte(0x55), z80.reg.A)
	assert.Equal(t, 0, z80.t)
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
		z80 := NewZ80(mem)
		z80.reg.F = test.flag
		if test.expected == 0x55 {
			z80.Run(10 + 11 + 7)
		} else {
			z80.Run(10 + 5 + 7)
		}

		assert.Equal(t, byte(test.expected), z80.reg.A)
		assert.Equal(t, 0, z80.t)
	}
}

func Test_RETN_RETI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_sp_nn, 0x0B, 0x00, prefix_ed, retn, ld_a_n, 0xAA, halt, ld_a_n, 0x55, halt, 0x08, 0x00}}
	z80 := NewZ80(mem)
	z80.iff2 = true
	z80.Run(10 + 14 + 7)

	assert.Equal(t, byte(0x55), z80.reg.A)
	assert.Equal(t, true, z80.iff1)
	assert.Equal(t, true, z80.iff2)
	assert.Equal(t, 0, z80.t)
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
			rst_28h, rst_30h, rst_38h, ld_b_n, 0x55, nop, 0, 0,
		},
	}
	z80 := NewZ80(mem)
	z80.reg.PC = 0x40
	z80.Run(10 + 8*(11+7+10) + 7)

	assert.Equal(t, byte(0xFF), z80.reg.A)
	assert.Equal(t, byte(0x55), z80.reg.B)
	assert.Equal(t, 0, z80.t)
}

func Test_PUSH_rr(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{ld_sp_nn, 0x23, 0x00, ld_a_n, 0x98,
			ld_bc_nn, 0x34, 0x12, ld_de_nn, 0x35, 0x13, ld_hl_nn, 0x36, 0x14,
			push_af, push_bc, push_de, push_hl, useIX, push_hl, useIY, push_hl, nop,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	z80 := NewZ80(mem)
	z80.reg.F = fS | fC
	z80.Run(10 + 7 + 10 + 10 + 10 + 11 + 11 + 11 + 11 + 15 + 15)

	assert.Equal(t, z80.mem.Read(34), z80.reg.A)
	assert.Equal(t, z80.mem.Read(33), z80.reg.F)
	assert.Equal(t, z80.mem.Read(32), z80.reg.B)
	assert.Equal(t, z80.mem.Read(31), z80.reg.C)
	assert.Equal(t, z80.mem.Read(30), z80.reg.D)
	assert.Equal(t, z80.mem.Read(29), z80.reg.E)
	assert.Equal(t, z80.mem.Read(28), z80.reg.H)
	assert.Equal(t, z80.mem.Read(27), z80.reg.L)
	assert.Equal(t, z80.mem.Read(26), z80.reg.IXH)
	assert.Equal(t, z80.mem.Read(25), z80.reg.IXL)
	assert.Equal(t, z80.mem.Read(24), z80.reg.IYH)
	assert.Equal(t, z80.mem.Read(23), z80.reg.IYL)
	assert.Equal(t, 0, z80.t)
}

func Test_POP_rr(t *testing.T) {
	mem := &memory.BasicMemory{
		Cells: []byte{ld_sp_nn, 0x0C, 0x00, pop_af, pop_bc, pop_de, pop_hl, useIX, pop_hl, useIY, pop_hl, nop,
			0x43, 0x21, 0x44, 0x22, 0x45, 0x23, 0x46, 0x24, 0x47, 0x25, 0x48, 0x26},
	}
	z80 := NewZ80(mem)
	z80.Run(10 + 10 + 10 + 10 + 10 + 14 + 14)

	assert.Equal(t, byte(0x21), z80.reg.A)
	assert.Equal(t, byte(0x43), z80.reg.F)
	assert.Equal(t, byte(0x22), z80.reg.B)
	assert.Equal(t, byte(0x44), z80.reg.C)
	assert.Equal(t, byte(0x23), z80.reg.D)
	assert.Equal(t, byte(0x45), z80.reg.E)
	assert.Equal(t, byte(0x24), z80.reg.H)
	assert.Equal(t, byte(0x46), z80.reg.L)
	assert.Equal(t, byte(0x25), z80.reg.IXH)
	assert.Equal(t, byte(0x47), z80.reg.IXL)
	assert.Equal(t, byte(0x26), z80.reg.IYH)
	assert.Equal(t, byte(0x48), z80.reg.IYL)
	assert.Equal(t, uint16(0x18), z80.reg.SP)
	assert.Equal(t, 0, z80.t)
}

func Test_IN_A_n(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x23, in_a_n, 0x01}}
	z80 := NewZ80(mem)
	z80.Run(7 + 11)

	assert.Equal(t, byte(0xFF), z80.reg.A)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	z80.IOBus = &TestIOBus{
		read: func(hi, lo byte) byte {
			if hi == 0x23 && lo == 0x01 {
				return 0xA5
			}
			return 0
		},
	}
	z80.Run(7 + 11)

	assert.Equal(t, byte(0xA5), z80.reg.A)
	assert.Equal(t, 0, z80.t)
}

func Test_IN_R_C(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_bc_nn, 0x23, 0x01, ld_d_n, 0x01, prefix_ed, in_d_c}}
	z80 := NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(10 + 7 + 12)

	assert.Equal(t, byte(0xFF), z80.reg.D)
	assert.Equal(t, fS|fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	z80.Reset()
	z80.reg.F = fNONE
	z80.IOBus = &TestIOBus{
		read: func(hi, lo byte) byte {
			if hi == 0x01 && lo == 0x23 {
				return 0
			}
			return 0xA5
		},
	}

	z80.Run(10 + 7 + 12)

	assert.Equal(t, byte(0), z80.reg.D)
	assert.Equal(t, fZ|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_OUT_n_A(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_a_n, 0x23, out_n_a, 0x01}}
	z80 := NewZ80(mem)
	z80.IOBus = &TestIOBus{
		write: func(hi, lo, data byte, ts int) {
			assert.Equal(t, byte(0x23), hi)
			assert.Equal(t, byte(0x01), lo)
			assert.Equal(t, byte(0x23), data)
		},
	}
	z80.Run(10 + 7 + 12)
	assert.Equal(t, 0, z80.t)
}

func Test_OUT_C_R(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_bc_nn, 0x11, 0x22, ld_h_n, 0x33, prefix_ed, out_c_h}}
	z80 := NewZ80(mem)
	z80.IOBus = &TestIOBus{
		write: func(hi, lo, data byte, ts int) {
			assert.Equal(t, byte(0x22), hi)
			assert.Equal(t, byte(0x11), lo)
			assert.Equal(t, byte(0x33), data)
		},
	}
	z80.Run(10 + 7 + 12)
	assert.Equal(t, 0, z80.t)
}

func Test_RLC_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x55, prefix_cb, rlc_r | rE}}
	z80 := NewZ80(mem)
	z80.reg.F = fZ | fH | fN | fC
	z80.Run(7 + 8)

	assert.Equal(t, byte(0xAA), z80.reg.E)
	assert.Equal(t, fS|fY|fX|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_d_n, 0xAA, prefix_cb, rlc_r | rD}}
	z80 = NewZ80(mem)
	z80.reg.F = fZ | fH | fN | fC
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x55), z80.reg.D)
	assert.Equal(t, fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x00, prefix_cb, rlc_r | rA}}
	z80 = NewZ80(mem)
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x00), z80.reg.A)
	assert.Equal(t, fZ|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_b_n, 0x80, prefix_cb, rlc_r | rB}}
	z80 = NewZ80(mem)
	z80.reg.F = fC
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x01), z80.reg.B)
	assert.Equal(t, fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_cb, rlc_r | 0b110, nop, 0x01}}
	z80 = NewZ80(mem)
	z80.Run(10 + 15)

	assert.Equal(t, byte(0x02), z80.mem.Read(0x06))
	assert.Equal(t, fNONE, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{useIX, ld_hl_nn, 0x04, 0x00, useIX, prefix_cb, 0x05, rlc_r | 0b110, nop, 0x01}}
	z80 = NewZ80(mem)
	z80.Run(14 + 23)

	assert.Equal(t, byte(0x02), z80.mem.Read(0x09))
	assert.Equal(t, fNONE, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_RRC_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x55, prefix_cb, rrc_r | rE}}
	z80 := NewZ80(mem)
	z80.reg.F = fZ | fH | fN
	z80.Run(7 + 8)

	assert.Equal(t, byte(0xAA), z80.reg.E)
	assert.Equal(t, fS|fY|fX|fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_d_n, 0xAA, prefix_cb, rrc_r | rD}}
	z80 = NewZ80(mem)
	z80.reg.F = fS | fZ | fH | fN | fC
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x55), z80.reg.D)
	assert.Equal(t, fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x00, prefix_cb, rrc_r | rA}}
	z80 = NewZ80(mem)
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x00), z80.reg.A)
	assert.Equal(t, fZ|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_b_n, 0x80, prefix_cb, rrc_r | rB}}
	z80 = NewZ80(mem)
	z80.reg.F = fC
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x40), z80.reg.B)
	assert.Equal(t, fNONE, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_cb, rrc_r | 0b110, nop, 0x01}}
	z80 = NewZ80(mem)
	z80.Run(10 + 15)

	assert.Equal(t, byte(0x80), z80.mem.Read(0x06))
	assert.Equal(t, fS|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_RL_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x55, prefix_cb, rl_r | rE}}
	z80 := NewZ80(mem)
	z80.reg.F = fZ | fH | fN | fC
	z80.Run(7 + 8)

	assert.Equal(t, byte(0xAB), z80.reg.E)
	assert.Equal(t, fS|fY|fX, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_d_n, 0xAA, prefix_cb, rl_r | rD}}
	z80 = NewZ80(mem)
	z80.reg.F = fZ | fH | fN | fC
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x55), z80.reg.D)
	assert.Equal(t, fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x80, prefix_cb, rl_r | rA}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x00), z80.reg.A)
	assert.Equal(t, fZ|fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_b_n, 0x80, prefix_cb, rl_r | rB}}
	z80 = NewZ80(mem)
	z80.reg.F = fC
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x01), z80.reg.B)
	assert.Equal(t, fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_cb, rl_r | 0b110, halt, 0x81}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(10 + 15)

	assert.Equal(t, byte(0x02), z80.mem.Read(0x06))
	assert.Equal(t, fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_RR_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x55, prefix_cb, rr_r | rE}}
	z80 := NewZ80(mem)
	z80.reg.F = fZ | fH | fN | fC
	z80.Run(7 + 8)

	assert.Equal(t, byte(0xAA), z80.reg.E)
	assert.Equal(t, fS|fY|fX|fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_d_n, 0xAA, prefix_cb, rr_r | rD}}
	z80 = NewZ80(mem)
	z80.reg.F = fZ | fH | fN | fC
	z80.Run(7 + 8)

	assert.Equal(t, byte(0xD5), z80.reg.D)
	assert.Equal(t, fS, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_a_n, 0x01, prefix_cb, rr_r | rA}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x00), z80.reg.A)
	assert.Equal(t, fZ|fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_b_n, 0x80, prefix_cb, rr_r | rB}}
	z80 = NewZ80(mem)
	z80.reg.F = fC
	z80.Run(7 + 8)

	assert.Equal(t, byte(0xC0), z80.reg.B)
	assert.Equal(t, fS|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_cb, rr_r | 0b110, halt, 0x81}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(10 + 15)

	assert.Equal(t, byte(0x40), z80.mem.Read(0x06))
	assert.Equal(t, fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_SLA_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x55, prefix_cb, sla_r | rE}}
	z80 := NewZ80(mem)
	z80.reg.F = fZ | fH | fN | fC
	z80.Run(7 + 8)

	assert.Equal(t, byte(0xAA), z80.reg.E)
	assert.Equal(t, fS|fY|fX|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_d_n, 0xAA, prefix_cb, sla_r | rD}}
	z80 = NewZ80(mem)
	z80.reg.F = fZ | fH | fN
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x54), z80.reg.D)
	assert.Equal(t, fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_SRA_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x85, prefix_cb, sra_r | rE}}
	z80 := NewZ80(mem)
	z80.reg.F = fZ | fH | fN
	z80.Run(7 + 8)

	assert.Equal(t, byte(0xC2), z80.reg.E)
	assert.Equal(t, fS|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_SLL_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x95, prefix_cb, sll_r | rE}}
	z80 := NewZ80(mem)
	z80.reg.F = fS | fZ | fH | fN | fC
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x2B), z80.reg.E)
	assert.Equal(t, fY|fX|fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_SRL_r(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_h_n, 0x85, prefix_cb, srl_r | rH}}
	z80 := NewZ80(mem)
	z80.reg.F = fS | fH | fN
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x42), z80.reg.H)
	assert.Equal(t, fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{useIX, ld_hl_nn, 0x00, 0x00, useIX, prefix_cb, 0x08, srl_r | rB, 0x85}}
	z80 = NewZ80(mem)
	z80.reg.F = fS | fH | fN
	z80.Run(14 + 23)

	assert.Equal(t, byte(0x42), z80.reg.B)
	assert.Equal(t, byte(0x42), z80.mem.Read(8))
	assert.Equal(t, fP|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_BIT_b(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_e_n, 0x40, prefix_cb, bit_b | rE | bit_6}}
	z80 := NewZ80(mem)
	z80.reg.F = fZ | fN | fC
	z80.Run(7 + 8)

	assert.Equal(t, fH|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_l_n, 0xFE, prefix_cb, bit_b | rL | bit_0}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(7 + 8)

	assert.Equal(t, fZ|fY|fH|fX|fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_cb, bit_b | 0b110 | bit_2, nop, 0xFD}}
	z80 = NewZ80(mem)
	z80.reg.F = fZ | fN
	z80.Run(10 + 12)

	assert.Equal(t, fH, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{useIY, ld_hl_nn, 0x07, 0x00, useIY, prefix_cb, 0x02, bit_b | 0b110 | bit_2, nop, 0xFD}}
	z80 = NewZ80(mem)
	z80.reg.F = fZ | fN
	z80.Run(14 + 20)

	assert.Equal(t, fH, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_RES_b(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_d_n, 0xFF, prefix_cb, res_b | rD | bit_7}}
	z80 := NewZ80(mem)
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x7F), z80.reg.D)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_cb, res_b | 0b110 | bit_2, nop, 0xFF}}
	z80 = NewZ80(mem)
	z80.reg.F = fZ | fN
	z80.Run(10 + 15)

	assert.Equal(t, byte(0xFB), z80.mem.Read(0x06))
	assert.Equal(t, 0, z80.t)
}

func Test_SET_b(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{ld_d_n, 0x00, prefix_cb, set_b | rD | bit_7}}
	z80 := NewZ80(mem)
	z80.Run(7 + 8)

	assert.Equal(t, byte(0x80), z80.reg.D)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{ld_hl_nn, 0x06, 0x00, prefix_cb, set_b | 0b110 | bit_2, nop, 0x00}}
	z80 = NewZ80(mem)
	z80.reg.F = fZ | fN
	z80.Run(10 + 15)

	assert.Equal(t, byte(0x04), z80.mem.Read(0x06))
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{useIY, ld_hl_nn, 0x06, 0x00, useIY, prefix_cb, 0x03, set_b | bit_2, nop, 0x00}}
	z80 = NewZ80(mem)
	z80.reg.F = fZ | fN
	z80.Run(14 + 23)

	assert.Equal(t, byte(0x04), z80.mem.Read(0x09))
	assert.Equal(t, byte(0x04), z80.reg.B)
	assert.Equal(t, 0, z80.t)
}

func Test_LD_IXY_nn(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		useIX, ld_hl_nn, 0x06, 0x01, useIY, ld_hl_nn, 0x07, 0x02,
		ld_hl_nn, 0x08, 0x03}}
	z80 := NewZ80(mem)
	z80.Run(14 + 14 + 10)

	assert.Equal(t, byte(0x01), z80.reg.IXH)
	assert.Equal(t, byte(0x06), z80.reg.IXL)
	assert.Equal(t, byte(0x02), z80.reg.IYH)
	assert.Equal(t, byte(0x07), z80.reg.IYL)
	assert.Equal(t, byte(0x03), z80.reg.H)
	assert.Equal(t, byte(0x08), z80.reg.L)
	assert.Equal(t, 0, z80.t)
}

func Test_LDI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0C, 0x00, ld_de_nn, 0x0D, 0x00, ld_bc_nn, 0x01, 0x00,
		prefix_ed, ldi, nop, 0xA5, 0x00}}
	z80 := NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(10 + 10 + 10 + 16)

	assert.Equal(t, byte(0xA5), z80.mem.Read(0x0D))
	assert.Equal(t, uint16(0), z80.reg.BC())
	assert.Equal(t, uint16(0x0E), z80.reg.DE())
	assert.Equal(t, uint16(0x0D), z80.reg.HL())
	assert.Equal(t, fS|fZ|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0C, 0x00, ld_de_nn, 0x0D, 0x00, ld_bc_nn, 0x02, 0x00,
		prefix_ed, ldi, nop, 0xA5, 0x00}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(10 + 10 + 10 + 16)

	assert.Equal(t, byte(0xA5), z80.mem.Read(0x0D))
	assert.Equal(t, uint16(1), z80.reg.BC())
	assert.Equal(t, uint16(0x0E), z80.reg.DE())
	assert.Equal(t, uint16(0x0D), z80.reg.HL())
	assert.Equal(t, fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_LDIR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0C, 0x00, ld_de_nn, 0x0F, 0x00, ld_bc_nn, 0x03, 0x00,
		prefix_ed, ldir, nop, 0x88, 0x36, 0xA5, 0x00, 0x00, 0x00}}
	z80 := NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(10 + 10 + 10 + 21 + 21 + 16)

	assert.Equal(t, byte(0x88), z80.mem.Read(0x0F))
	assert.Equal(t, byte(0x36), z80.mem.Read(0x10))
	assert.Equal(t, byte(0xA5), z80.mem.Read(0x11))
	assert.Equal(t, uint16(0), z80.reg.BC())
	assert.Equal(t, uint16(0x12), z80.reg.DE())
	assert.Equal(t, uint16(0x0F), z80.reg.HL())
	assert.Equal(t, fS|fZ|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_CPI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0B, 0x00, ld_bc_nn, 0x03, 0x00, ld_a_n, 0x88,
		prefix_ed, cpi, nop, 0x88}}
	z80 := NewZ80(mem)
	z80.reg.F = fC
	z80.Run(10 + 10 + 7 + 16)

	assert.Equal(t, uint16(0x02), z80.reg.BC())
	assert.Equal(t, uint16(0x0C), z80.reg.HL())
	assert.Equal(t, fZ|fP|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0B, 0x00, ld_bc_nn, 0x01, 0x00, ld_a_n, 0x88,
		prefix_ed, cpi, nop, 0x89}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(10 + 10 + 7 + 16)

	assert.Equal(t, uint16(0x00), z80.reg.BC())
	assert.Equal(t, uint16(0x0C), z80.reg.HL())
	assert.Equal(t, fS|fY|fH|fX|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_CPIR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0B, 0x00, ld_bc_nn, 0xFF, 0x00, ld_a_n, 0x88,
		prefix_ed, cpir, nop, 0x02, 0x04, 0x80, 0x88, 0x90}}
	z80 := NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(10 + 10 + 7 + 21 + 21 + 21 + 16)

	assert.Equal(t, uint16(0xFB), z80.reg.BC())
	assert.Equal(t, uint16(0x0F), z80.reg.HL())
	assert.Equal(t, fZ|fP|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_INI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x09, 0x00, ld_bc_nn, 0x34, 0x01,
		prefix_ed, ini, nop, 0x00}}
	z80 := NewZ80(mem)
	z80.reg.F = fC
	z80.IOBus = &TestIOBus{
		read: func(hi, lo byte) byte {
			if hi == 0x01 && lo == 0x34 {
				return 0x5E
			}
			return 0
		},
	}
	z80.Run(10 + 10 + 16)

	assert.Equal(t, uint16(0x34), z80.reg.BC())
	assert.Equal(t, byte(0x5E), z80.mem.Read((9)))
	assert.Equal(t, uint16(0x0A), z80.reg.HL())
	assert.Equal(t, fZ|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_INIR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x09, 0x00, ld_bc_nn, 0x34, 0x05,
		prefix_ed, inir, nop, 0x00, 0x00, 0x00, 0x00, 0x00}}
	z80 := NewZ80(mem)
	z80.reg.F = fC
	z80.IOBus = &TestIOBus{
		read: func(hi, lo byte) byte {
			if lo == 0x34 {
				return hi + 0x20
			}
			return 0
		},
	}
	z80.Run(10 + 10 + 21 + 21 + 21 + 21 + 16)

	assert.Equal(t, uint16(0x34), z80.reg.BC())
	assert.Equal(t, byte(0x25), z80.mem.Read((9)))
	assert.Equal(t, byte(0x24), z80.mem.Read((10)))
	assert.Equal(t, byte(0x23), z80.mem.Read((11)))
	assert.Equal(t, byte(0x22), z80.mem.Read((12)))
	assert.Equal(t, byte(0x21), z80.mem.Read((13)))
	assert.Equal(t, uint16(0x0E), z80.reg.HL())
	assert.Equal(t, fZ|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_OUTI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x09, 0x00, ld_bc_nn, 0x34, 0x01,
		prefix_ed, outi, nop, 0x87}}
	z80 := NewZ80(mem)
	z80.reg.F = fC
	z80.IOBus = &TestIOBus{
		write: func(hi, lo, data byte, ts int) {
			assert.Equal(t, byte(0), hi)
			assert.Equal(t, byte(0x34), lo)
			assert.Equal(t, byte(0x87), data)
		},
	}
	z80.Run(10 + 10 + 16)

	assert.Equal(t, uint16(0x34), z80.reg.BC())
	assert.Equal(t, uint16(0x0A), z80.reg.HL())
	assert.Equal(t, fZ|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_OTIR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x09, 0x00, ld_bc_nn, 0x34, 0x04,
		prefix_ed, otir, nop, 0x87, 0x88, 0x89, 0x8A}}
	z80 := NewZ80(mem)
	z80.reg.F = fC
	z80.IOBus = &TestIOBus{
		write: func(hi, lo, data byte, ts int) {
			assert.Equal(t, byte(0x34), lo)
			assert.Equal(t, byte(0x87+0x03-hi), data)
		},
	}
	z80.Run(10 + 10 + 21 + 21 + 21 + 16)

	assert.Equal(t, uint16(0x34), z80.reg.BC())
	assert.Equal(t, uint16(0x0D), z80.reg.HL())
	assert.Equal(t, fZ|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_LDD(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0C, 0x00, ld_de_nn, 0x0D, 0x00, ld_bc_nn, 0x01, 0x00,
		prefix_ed, ldd, nop, 0xA5, 0x00}}
	z80 := NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(10 + 10 + 10 + 16)

	assert.Equal(t, byte(0xA5), z80.mem.Read(0x0D))
	assert.Equal(t, uint16(0x00), z80.reg.BC())
	assert.Equal(t, uint16(0x0C), z80.reg.DE())
	assert.Equal(t, uint16(0x0B), z80.reg.HL())
	assert.Equal(t, fS|fZ|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0C, 0x00, ld_de_nn, 0x0D, 0x00, ld_bc_nn, 0x02, 0x00,
		prefix_ed, ldd, nop, 0xA5, 0x00}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(10 + 10 + 10 + 16)

	assert.Equal(t, byte(0xA5), z80.mem.Read(0x0D))
	assert.Equal(t, uint16(0x01), z80.reg.BC())
	assert.Equal(t, uint16(0x0C), z80.reg.DE())
	assert.Equal(t, uint16(0x0B), z80.reg.HL())
	assert.Equal(t, fP, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_LDDR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0E, 0x00, ld_de_nn, 0x11, 0x00, ld_bc_nn, 0x03, 0x00,
		prefix_ed, lddr, nop, 0x88, 0x36, 0xA5, 0x00, 0x00, 0x00}}
	z80 := NewZ80(mem)
	z80.reg.F = fALL
	z80.Run(10 + 10 + 10 + 21 + 21 + 16)

	assert.Equal(t, byte(0x88), z80.mem.Read(0x0F))
	assert.Equal(t, byte(0x36), z80.mem.Read(0x10))
	assert.Equal(t, byte(0xA5), z80.mem.Read(0x11))
	assert.Equal(t, uint16(0x00), z80.reg.BC())
	assert.Equal(t, uint16(0x0E), z80.reg.DE())
	assert.Equal(t, uint16(0x0B), z80.reg.HL())
	assert.Equal(t, fS|fZ|fY|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_CPD(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0B, 0x00, ld_bc_nn, 0x03, 0x00, ld_a_n, 0x88,
		prefix_ed, cpd, nop, 0x88}}
	z80 := NewZ80(mem)
	z80.reg.F = fC
	z80.Run(10 + 10 + 7 + 16)

	assert.Equal(t, uint16(0x02), z80.reg.BC())
	assert.Equal(t, uint16(0x0A), z80.reg.HL())
	assert.Equal(t, fZ|fP|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)

	mem = &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0B, 0x00, ld_bc_nn, 0x01, 0x00, ld_a_n, 0x88,
		prefix_ed, cpd, nop, 0x89}}
	z80 = NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(10 + 10 + 7 + 16)

	assert.Equal(t, uint16(0x00), z80.reg.BC())
	assert.Equal(t, uint16(0x0A), z80.reg.HL())
	assert.Equal(t, fS|fY|fH|fX|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_CPDR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0F, 0x00, ld_bc_nn, 0xFF, 0x00, ld_a_n, 0x88,
		prefix_ed, cpdr, nop, 0x02, 0x04, 0x80, 0x88, 0x90}}
	z80 := NewZ80(mem)
	z80.reg.F = fNONE
	z80.Run(10 + 10 + 7 + 21 + 16)

	assert.Equal(t, uint16(0xFD), z80.reg.BC())
	assert.Equal(t, uint16(0x0D), z80.reg.HL())
	assert.Equal(t, fZ|fP|fN, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_IND(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x09, 0x00, ld_bc_nn, 0x34, 0x01,
		prefix_ed, ind, nop, 0x00}}
	z80 := NewZ80(mem)
	z80.reg.F = fC
	z80.IOBus = &TestIOBus{
		read: func(hi, lo byte) byte {
			if hi == 0x01 && lo == 0x34 {
				return 0x5E
			}
			return 0
		},
	}
	z80.Run(10 + 10 + 16)

	assert.Equal(t, uint16(0x34), z80.reg.BC())
	assert.Equal(t, byte(0x5E), z80.mem.Read((9)))
	assert.Equal(t, uint16(0x08), z80.reg.HL())
	assert.Equal(t, fZ|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_INDR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0D, 0x00, ld_bc_nn, 0x34, 0x05,
		prefix_ed, indr, nop, 0x00, 0x00, 0x00, 0x00, 0x00}}
	z80 := NewZ80(mem)
	z80.reg.F = fC
	z80.IOBus = &TestIOBus{
		read: func(hi, lo byte) byte {
			if lo == 0x34 {
				return hi + 0x20
			}
			return 0
		},
	}
	z80.Run(10 + 10 + 21 + 21 + 21 + 21 + 16)

	assert.Equal(t, uint16(0x34), z80.reg.BC())
	assert.Equal(t, byte(0x21), z80.mem.Read((9)))
	assert.Equal(t, byte(0x22), z80.mem.Read((10)))
	assert.Equal(t, byte(0x23), z80.mem.Read((11)))
	assert.Equal(t, byte(0x24), z80.mem.Read((12)))
	assert.Equal(t, byte(0x25), z80.mem.Read((13)))
	assert.Equal(t, uint16(0x08), z80.reg.HL())
	assert.Equal(t, fZ|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_OUTD(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x09, 0x00, ld_bc_nn, 0x34, 0x01,
		prefix_ed, outd, nop, 0x87}}
	z80 := NewZ80(mem)
	z80.reg.F = fC
	z80.IOBus = &TestIOBus{
		write: func(hi, lo, data byte, ts int) {
			assert.Equal(t, byte(0), hi)
			assert.Equal(t, byte(0x34), lo)
			assert.Equal(t, byte(0x87), data)
		},
	}
	z80.Run(10 + 10 + 16)

	assert.Equal(t, uint16(0x34), z80.reg.BC())
	assert.Equal(t, uint16(0x08), z80.reg.HL())
	assert.Equal(t, fZ|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_OTDR(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x0C, 0x00, ld_bc_nn, 0x34, 0x04,
		prefix_ed, otdr, nop, 0x87, 0x88, 0x89, 0x8A}}
	z80 := NewZ80(mem)
	z80.reg.F = fC
	z80.IOBus = &TestIOBus{
		write: func(hi, lo, data byte, ts int) {
			assert.Equal(t, byte(0x34), lo)
			assert.Equal(t, byte(0x8A-(0x03-hi)), data)
		},
	}
	z80.Run(10 + 10 + 21 + 21 + 21 + 16)

	assert.Equal(t, uint16(0x34), z80.reg.BC())
	assert.Equal(t, uint16(0x08), z80.reg.HL())
	assert.Equal(t, fZ|fN|fC, z80.reg.F)
	assert.Equal(t, 0, z80.t)
}

func Test_invalidPrefix(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{useIX, useIX, useIX, ld_l_n, 0x01, useIY, useIY, ld_l_n, 0x02}}
	z80 := NewZ80(mem)
	z80.Run(4 + 4 + 11 + 4 + 11)

	assert.Equal(t, byte(0x01), z80.reg.IXL)
	assert.Equal(t, byte(0x02), z80.reg.IYL)
	assert.Equal(t, 0, z80.t)
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

	z80 := NewZ80(&memory.BasicMemory{})
	for _, test := range tests {
		z80.reg.F = test.flags
		result := z80.shouldJump(test.code)

		assert.Equal(t, test.expected, result)
	}
}

func Test_getHL(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{
		ld_hl_nn, 0x34, 0x12, useIX, ld_hl_nn, 0x45, 0x23, nop, 0x00, 0x02, 0xFE}}
	z80 := NewZ80(mem)
	z80.Run(10 + 14 + 4)

	hl := z80.getHL()
	assert.Equal(t, uint16(0x1234), hl)

	hl = z80.getHL()
	assert.Equal(t, uint16(0x1234), hl)

	z80.reg.prefix = useIX
	hl = z80.getHL()
	assert.Equal(t, uint16(0x2345), hl)

	hl = z80.getHL()
	assert.Equal(t, uint16(0x2347), hl)

	hl = z80.getHL()
	assert.Equal(t, uint16(0x2343), hl)
}

func Test_NMI(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{0x00, 0x00, 0x00, 0x00}}
	z80 := NewZ80(mem)
	z80.reg.PC = 0x1234
	z80.reg.SP = 0x04
	z80.halt, z80.iff1 = true, true

	z80.NMI()
	assert.Equal(t, uint16(0x02), z80.reg.SP)
	assert.Equal(t, uint16(0x66), z80.reg.PC)
	assert.Equal(t, byte(0x12), mem.Read(0x03))
	assert.Equal(t, byte(0x34), mem.Read(0x02))
	assert.Equal(t, false, z80.halt)
	assert.Equal(t, false, z80.iff1)
	assert.Equal(t, true, z80.iff2)
}

func Test_INT(t *testing.T) {
	mem := &memory.BasicMemory{Cells: []byte{0x00, 0x00, 0x00, 0x00}}
	z80 := NewZ80(mem)
	z80.reg.PC = 0x1234
	z80.reg.SP = 0x04
	z80.iff1 = false

	z80.INT(0)
	assert.Equal(t, uint16(0x04), z80.reg.SP)
	assert.Equal(t, uint16(0x1234), z80.reg.PC)

	z80.halt, z80.iff1, z80.iff2 = true, true, true
	z80.im = im1
	z80.INT(0)
	assert.Equal(t, uint16(0x02), z80.reg.SP)
	assert.Equal(t, uint16(0x38), z80.reg.PC)
	assert.Equal(t, byte(0x12), mem.Read(0x03))
	assert.Equal(t, byte(0x34), mem.Read(0x02))
	assert.Equal(t, false, z80.halt)
	assert.Equal(t, false, z80.iff1)
	assert.Equal(t, false, z80.iff2)

	z80.halt, z80.iff1, z80.iff2 = true, true, true
	z80.im = im2
	z80.reg.I = 0x23
	z80.INT(0x45)
	assert.Equal(t, uint16(0x00), z80.reg.SP)
	assert.Equal(t, uint16(0xFFFF), z80.reg.PC)
	assert.Equal(t, byte(0x00), mem.Read(0x01))
	assert.Equal(t, byte(0x38), mem.Read(0x00))
	assert.Equal(t, false, z80.halt)
	assert.Equal(t, false, z80.iff1)
	assert.Equal(t, false, z80.iff2)
}

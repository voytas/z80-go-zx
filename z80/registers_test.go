package z80

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getR(t *testing.T) {
	r := newRegisters()
	r.A, r.B, r.C, r.D, r.E, r.H, r.L = 1, 2, 3, 4, 5, 6, 7
	assert.Equal(t, &r.A, r.raw[rA])
	assert.Equal(t, &r.B, r.raw[rB])
	assert.Equal(t, &r.C, r.raw[rC])
	assert.Equal(t, &r.D, r.raw[rD])
	assert.Equal(t, &r.E, r.raw[rE])
	assert.Equal(t, &r.H, r.raw[rH])
	assert.Equal(t, &r.L, r.raw[rL])
}

func Test_getRR(t *testing.T) {
	r := newRegisters()
	r.B, r.C, r.D, r.E, r.H, r.L = 2, 3, 4, 5, 6, 7
	r.IXH, r.IXL, r.IYH, r.IYL = 8, 9, 0xA, 0xB

	assert.Equal(t, uint16(0x0203), r.BC())
	assert.Equal(t, uint16(0x0405), r.DE())
	assert.Equal(t, uint16(0x0607), r.HL())
	assert.Equal(t, uint16(0x0809), r.IX())
	assert.Equal(t, uint16(0x0A0B), r.IY())
}

func Test_setRR(t *testing.T) {
	r := newRegisters()
	r.SetBC(0x1122)
	r.SetDE(0x3344)
	r.SetHL(0x5566)

	assert.Equal(t, byte(0x11), r.B)
	assert.Equal(t, byte(0x22), r.C)
	assert.Equal(t, byte(0x33), r.D)
	assert.Equal(t, byte(0x44), r.E)
	assert.Equal(t, byte(0x55), r.H)
	assert.Equal(t, byte(0x66), r.L)
}

func Test_getReg(t *testing.T) {
	r := newRegisters()
	r.A, r.B, r.C, r.D, r.E, r.H, r.L = 0x01, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08
	r.IXH, r.IXL = 0x09, 0x0A
	r.IYH, r.IYL = 0x0B, 0x0C

	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		for _, reg := range []byte{rA, rB, rC, rD, rE, rH, rL} {
			r.prefix = prefix
			result := *r.r(reg)
			switch reg {
			case rA:
				assert.Equal(t, r.A, result)
			case rB:
				assert.Equal(t, r.B, result)
			case rC:
				assert.Equal(t, r.C, result)
			case rD:
				assert.Equal(t, r.D, result)
			case rE:
				assert.Equal(t, r.E, result)
			case rH:
				switch prefix {
				case useIX:
					assert.Equal(t, r.IXH, result)
				case useIY:
					assert.Equal(t, r.IYH, result)
				default:
					assert.Equal(t, r.H, result)
				}
			case rL:
				switch prefix {
				case useIX:
					assert.Equal(t, r.IXL, result)
				case useIY:
					assert.Equal(t, r.IYL, result)
				default:
					assert.Equal(t, r.L, result)
				}
			}
		}
	}
}

func Test_setReg(t *testing.T) {
	for _, prefix := range []byte{noPrefix, useIX, useIY} {
		for _, reg := range []byte{rA, rB, rC, rD, rE, rH, rL} {
			var val byte = 0x76
			r := newRegisters()
			r.prefix = prefix
			r.setR(reg, val)

			switch reg {
			case rA:
				assert.Equal(t, val, r.A)
			case rB:
				assert.Equal(t, val, r.B)
			case rC:
				assert.Equal(t, val, r.C)
			case rD:
				assert.Equal(t, val, r.D)
			case rE:
				assert.Equal(t, val, r.E)
			case rH:
				switch prefix {
				case useIX:
					assert.Equal(t, val, r.IXH)
				case useIY:
					assert.Equal(t, val, r.IYH)
				default:
					assert.Equal(t, val, r.H)
				}
			case rL:
				switch prefix {
				case useIX:
					assert.Equal(t, val, r.IXL)
				case useIY:
					assert.Equal(t, val, r.IYL)
				default:
					assert.Equal(t, val, r.L)
				}
			}
		}
	}
}

func Test_IncR(t *testing.T) {
	r := newRegisters()

	for i := 0; i < 127; i++ {
		r.IncR()
		assert.EqualValues(t, i+1, r.R)
	}

	for i := 0; i < 127; i++ {
		r.IncR()
		assert.EqualValues(t, i, r.R)
	}

	r.R = 0x80
	for i := 0; i < 127; i++ {
		r.IncR()
		assert.EqualValues(t, i+1+0x80, r.R)
	}

	for i := 0; i < 127; i++ {
		r.IncR()
		assert.EqualValues(t, i+0x80, r.R)
	}
}

package z80

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getR(t *testing.T) {
	r := newRegisters()
	r.A, r.B, r.C, r.D, r.E, r.H, r.L = 1, 2, 3, 4, 5, 6, 7
	assert.Equal(t, &r.A, r.raw[r_A])
	assert.Equal(t, &r.B, r.raw[r_B])
	assert.Equal(t, &r.C, r.raw[r_C])
	assert.Equal(t, &r.D, r.raw[r_D])
	assert.Equal(t, &r.E, r.raw[r_E])
	assert.Equal(t, &r.H, r.raw[r_H])
	assert.Equal(t, &r.L, r.raw[r_L])
}

func Test_getRR(t *testing.T) {
	r := newRegisters()
	r.B, r.C, r.D, r.E, r.H, r.L = 2, 3, 4, 5, 6, 7

	assert.Equal(t, uint16(0x0203), r.BC())
	assert.Equal(t, uint16(0x0405), r.DE())
	assert.Equal(t, uint16(0x0607), r.HL())
}

func Test_setRR(t *testing.T) {
	r := newRegisters()
	r.setBC(0x1122)
	r.setDE(0x3344)
	r.setHL(0x5566)

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
		for _, reg := range []byte{r_A, r_B, r_C, r_D, r_E, r_H, r_L} {
			r.prefix = prefix
			result := *r.r(reg)
			switch reg {
			case r_A:
				assert.Equal(t, r.A, result)
			case r_B:
				assert.Equal(t, r.B, result)
			case r_C:
				assert.Equal(t, r.C, result)
			case r_D:
				assert.Equal(t, r.D, result)
			case r_E:
				assert.Equal(t, r.E, result)
			case r_H:
				switch prefix {
				case useIX:
					assert.Equal(t, r.IXH, result)
				case useIY:
					assert.Equal(t, r.IYH, result)
				default:
					assert.Equal(t, r.H, result)
				}
			case r_L:
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
		for _, reg := range []byte{r_A, r_B, r_C, r_D, r_E, r_H, r_L} {
			var val byte = 0x76
			r := newRegisters()
			r.prefix = prefix
			r.setR(reg, val)

			switch reg {
			case r_A:
				assert.Equal(t, val, r.A)
			case r_B:
				assert.Equal(t, val, r.B)
			case r_C:
				assert.Equal(t, val, r.C)
			case r_D:
				assert.Equal(t, val, r.D)
			case r_E:
				assert.Equal(t, val, r.E)
			case r_H:
				switch prefix {
				case useIX:
					assert.Equal(t, val, r.IXH)
				case useIY:
					assert.Equal(t, val, r.IYH)
				default:
					assert.Equal(t, val, r.H)
				}
			case r_L:
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

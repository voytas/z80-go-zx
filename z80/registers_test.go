package z80

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getR(t *testing.T) {
	r := newRegisters()
	r.A, r.B, r.C, r.D, r.E, r.H, r.L = 1, 2, 3, 4, 5, 6, 7
	assert.Equal(t, &r.A, r.getR(r_A))
	assert.Equal(t, &r.B, r.getR(r_B))
	assert.Equal(t, &r.C, r.getR(r_C))
	assert.Equal(t, &r.D, r.getR(r_D))
	assert.Equal(t, &r.E, r.getR(r_E))
	assert.Equal(t, &r.H, r.getR(r_H))
	assert.Equal(t, &r.L, r.getR(r_L))
}

func Test_getRR(t *testing.T) {
	r := newRegisters()
	r.B, r.C, r.D, r.E, r.H, r.L = 2, 3, 4, 5, 6, 7

	assert.Equal(t, word(0x0203), r.getBC())
	assert.Equal(t, word(0x0405), r.getDE())
	assert.Equal(t, word(0x0607), r.getHL())
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
	r.IX, r.IY = 0x090A, 0x0B0C

	for _, prefix := range []byte{prefix_none, prefix_ix, prefix_iy} {
		for _, reg := range []byte{r_A, r_B, r_C, r_D, r_E, r_H, r_L} {
			result := r.getReg(reg, prefix)
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
				case prefix_ix:
					assert.Equal(t, byte(r.IX>>8), result)
				case prefix_iy:
					assert.Equal(t, byte(r.IY>>8), result)
				default:
					assert.Equal(t, r.H, result)
				}
			case r_L:
				switch prefix {
				case prefix_ix:
					assert.Equal(t, byte(r.IX), result)
				case prefix_iy:
					assert.Equal(t, byte(r.IY), result)
				default:
					assert.Equal(t, r.L, result)
				}
			}
		}
	}
}

func Test_setReg(t *testing.T) {
	for _, prefix := range []byte{prefix_none, prefix_ix, prefix_iy} {
		for _, reg := range []byte{r_A, r_B, r_C, r_D, r_E, r_H, r_L} {
			var val byte = 0x76
			r := newRegisters()
			r.setReg(reg, prefix, val)

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
				case prefix_ix:
					assert.Equal(t, val, byte(r.IX>>8))
				case prefix_iy:
					assert.Equal(t, val, byte(r.IY>>8))
				default:
					assert.Equal(t, val, r.H)
				}
			case r_L:
				switch prefix {
				case prefix_ix:
					assert.Equal(t, val, byte(r.IX))
				case prefix_iy:
					assert.Equal(t, val, byte(r.IY))
				default:
					assert.Equal(t, val, r.L)
				}
			}
		}
	}
}

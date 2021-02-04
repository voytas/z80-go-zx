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

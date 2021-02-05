package z80

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BasicMemory(t *testing.T) {
	var tests = []struct {
		addr     word
		expected byte
	}{
		{0x0000, 0x55}, {0x3FFF, 0x55}, {0x4000, 0xAA}, {0x7FFF, 0xAA}, {0x8000, 0xFF}, {0xFFFF, 0xFF},
	}

	mem := BasicMemory{
		ramStart: 0x4000,
		cells:    make([]byte, 0x8000),
	}

	for i := 0; i < len(mem.cells); i++ {
		mem.cells[i] = 0x55
	}
	for _, test := range tests {
		mem.write(test.addr, 0xAA)
		result := *mem.read((test.addr))

		assert.Equal(t, byte(test.expected), result)
	}

	mem.write(0x4000, 0xFF)
	v1 := mem.read(0x4000)
	assert.Equal(t, byte(0xFF), *v1)

	*v1 = 0xAB
	v2 := *mem.read(0x4000)
	assert.Equal(t, byte(0xAB), v2)
}

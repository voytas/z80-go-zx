package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BasicMemory(t *testing.T) {
	var tests = []struct {
		addr     uint16
		expected byte
	}{
		{0x0000, 0x55}, {0x3FFF, 0x55}, {0x4000, 0xAA}, {0x7FFF, 0xAA}, {0x8000, 0xFF}, {0xFFFF, 0xFF},
	}

	mem := BasicMemory{
		ramStart: 0x4000,
		Cells:    make([]byte, 0x8000),
	}

	for i := 0; i < len(mem.Cells); i++ {
		mem.Cells[i] = 0x55
	}
	for _, test := range tests {
		mem.Write(test.addr, 0xAA)
		result := mem.Read(test.addr)

		assert.Equal(t, byte(test.expected), result)
	}
}
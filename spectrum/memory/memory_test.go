package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voytas/z80-go-zx/spectrum/machine"
	"github.com/voytas/z80-go-zx/z80"
)

func Test_ContentedRead(t *testing.T) {
	mem, err := NewMem48k("../rom/48.rom")
	assert.Nil(t, err)

	cpu := z80.NewZ80(mem)
	cpu.Reg.PC = 25000
	cpu.Reg.SetHL(26000)
	cpu.Reg.A = 0x34
	*mem.Cells[25000] = 0x77 // ld (hl),a
	*mem.Cells[25001] = 0x76 // halt
	contendedStates = machine.ZX48k.ContentionTable
	cpu.Trap = func() {
		cpu.TC.Current = 14335
		cpu.TC.Total = 14335
	}
	mem.TC = cpu.TC

	cpu.Run(14335 + 10)

	assert.Equal(t, cpu.Reg.A, *mem.Cells[26000])
}

package z80

import (
	"github.com/voytas/z80-go-zx/z80/dasm"
)

func (cpu *CPU) debug(opcode byte) string {
	if cpu.reg.prefix == noPrefix || opcode == useIX || opcode == useIY {
		return dasm.Decode(cpu.reg.PC-1, cpu.mem)
	}
	return ""
}

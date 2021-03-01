package z80

import (
	"fmt"

	"github.com/voytas/z80-go-zx/z80/dasm"
)

func (cpu *CPU) debug(opcode byte) {
	if cpu.reg.prefix == noPrefix || opcode == useIX || opcode == useIY {
		fmt.Println(dasm.Decode(cpu.reg.PC-1, cpu.mem))
	}
}

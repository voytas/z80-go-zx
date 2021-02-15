package z80

import (
	"fmt"

	"github.com/voytas/z80-go-zx/z80/dasm"
)

func (cpu *CPU) debug() {
	s := dasm.Decode(cpu.PC, cpu.mem)
	fmt.Println(s)
}

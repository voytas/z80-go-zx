package debugger

import (
	"fmt"

	"github.com/voytas/z80-go-zx/z80/dasm"
	"github.com/voytas/z80-go-zx/z80/memory"
)

// Very basic console output for debugging the opcodes
func Debug(opcode, prefix byte, PC uint16, mem memory.Memory) {
	if prefix == 0 || opcode == 0xDD || opcode == 0xFD {
		fmt.Println(dasm.Decode(PC-1, mem))
	}
}

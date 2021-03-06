package debugger

import (
	"fmt"

	"github.com/voytas/z80-go-zx/z80"
	"github.com/voytas/z80-go-zx/z80/dasm"
)

// Very basic console output for debugging the opcodes
func Debug(prefix, opcode byte, PC uint16, mem z80.Memory) {
	if prefix == 0 || opcode == 0xDD || opcode == 0xFD {
		fmt.Println(dasm.Decode(PC-1, mem))
	}
}

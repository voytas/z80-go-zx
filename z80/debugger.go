package z80

import (
	"fmt"

	"github.com/voytas/z80-go-zx/z80/dasm"
)

func (z80 *Z80) debug(opcode byte) {
	if z80.reg.prefix == noPrefix || opcode == useIX || opcode == useIY {
		fmt.Println(dasm.Decode(z80.reg.PC-1, z80.mem))
	}
}

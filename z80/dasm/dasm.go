package dasm

import (
	"fmt"
	"strings"

	"github.com/voytas/z80-go-zx/z80"
)

type instruction struct {
	mnemonic string
	size     int   // instruction size in bytes
	args     []int // index of arguments positions
}

// Decode current opcode into mnemonic. This is very basic and simple
// implementation, just a helper for debugging any issues.
func Decode(addr uint16, mem z80.Memory) string {
	opcode := mem.Read(addr)
	var inst *instruction
	switch opcode {
	case 0xCB:
		opcode = mem.Read(addr + 1)
		inst = cbInstructions[opcode]
	case 0xDD:
		opcode = mem.Read(addr + 1)
		if opcode == 0xCB {
			inst = ddcbInstructions[mem.Read(addr+3)]
		} else {
			inst = ddInstructions[opcode]
		}
	case 0xED:
		opcode = mem.Read(addr + 1)
		inst = edInstructions[opcode]
	case 0xFD:
		opcode = mem.Read(addr + 1)
		if opcode == 0xCB {
			inst = fdcbInstructions[mem.Read(addr+3)]
		} else {
			inst = fdInstructions[opcode]
		}
	default:
		inst = primaryInstructions[opcode]
	}

	if inst == nil {
		inst = &instruction{mnemonic: "[INVALID]", size: 2}
	}

	var bytes []byte
	for i := 0; i < inst.size; i++ {
		bytes = append(bytes, mem.Read(addr+uint16(i)))
	}
	s := fmt.Sprintf("%04X: ", addr) + fmtBytes(bytes)

	mnemonic := inst.mnemonic
	for i, arg := range inst.args {
		mnemonic = strings.ReplaceAll(mnemonic, fmt.Sprintf("$%v", i+1), fmt.Sprintf("%02X", bytes[arg]))
	}

	return s + mnemonic
}

func fmtBytes(bytes []byte) string {
	s := ""
	for _, b := range bytes {
		s += fmt.Sprintf("%02X ", b)
	}

	return fmt.Sprintf("%-13s ", s)
}

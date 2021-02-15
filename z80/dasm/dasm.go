package dasm

import (
	"fmt"
	"strings"

	"github.com/voytas/z80-go-zx/z80/memory"
)

type instruction struct {
	mnemonic string
	args     int // number of bytes to read to get arguments (0, 1 or 2)
	relative bool
}

func Decode(addr uint16, mem memory.Memory) string {
	opcode := mem.Read(addr)
	s := fmt.Sprintf("%04X ", addr)
	bytes := append([]byte{}, opcode)

	var inst *instruction
	switch opcode {
	case 0xCB:
		addr += 1
		opcode = mem.Read(addr)
		bytes = append(bytes, opcode)
		inst = cbInstructions[opcode]
	case 0xDD:
		addr += 1
		opcode = mem.Read(addr)
		inst = ddInstructions[opcode]
		if inst != nil {
			bytes = append(bytes, opcode)
		}
	case 0xED:
		opcode = mem.Read(addr)
		inst = edInstructions[opcode]
		if inst != nil {
			bytes = append(bytes, opcode)
		}
	case 0xFD:
		addr += 1
		opcode = mem.Read(addr)
		inst = fdInstructions[opcode]
		if inst != nil {
			bytes = append(bytes, opcode)
		}
	default:
		inst = primaryInstructions[opcode]
	}

	var mnemonic string
	if inst == nil {
		mnemonic = "[NOP]" // invalid opcode works like NOP
	} else if inst.args == 1 {
		v := mem.Read(addr + 1)
		bytes = append(bytes, v)
		mnemonic = strings.ReplaceAll(inst.mnemonic, "$1", fmt.Sprintf("%02X", v))
	} else if inst.args == 2 {
		l := mem.Read(addr + 1)
		bytes = append(bytes, l)
		h := mem.Read(addr + 2)
		bytes = append(bytes, h)
		mnemonic = strings.ReplaceAll(
			strings.ReplaceAll(inst.mnemonic, "$1", fmt.Sprintf("%02X", h)), "$2", fmt.Sprintf("%02X", l))
	} else {
		mnemonic = inst.mnemonic
	}

	s += fmtBytes(bytes) + mnemonic
	return s
}

func fmtBytes(bytes []byte) string {
	s := ""
	for _, b := range bytes {
		s += fmt.Sprintf("%02X ", b)
	}

	return fmt.Sprintf("%-12s ", s)
}

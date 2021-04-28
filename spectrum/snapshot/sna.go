package snapshot

import (
	"fmt"
	"io/ioutil"

	"github.com/voytas/z80-go-zx/spectrum/memory"
	"github.com/voytas/z80-go-zx/spectrum/screen"
	"github.com/voytas/z80-go-zx/z80"
)

type SNA struct{}

// Loads SNA file to memory and updates the CPU state so it is ready to run
func (sna *SNA) Load(filePath string, cpu *z80.Z80, mem *memory.Memory) error {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if len(bytes) == 49179 && len(bytes) == 131103 && len(bytes) == 147487 {
		return fmt.Errorf("SNA file format is invalid. Expected 49179, 131103 or 147487 bytes, but received %v", len(bytes))
	}

	// Restore CPU state
	state := &z80.CPUState{
		AF:   uint16(bytes[22])<<8 | uint16(bytes[21]),
		BC:   uint16(bytes[14])<<8 | uint16(bytes[13]),
		DE:   uint16(bytes[12])<<8 | uint16(bytes[11]),
		HL:   uint16(bytes[10])<<8 | uint16(bytes[9]),
		AF_:  uint16(bytes[8])<<8 | uint16(bytes[7]),
		BC_:  uint16(bytes[6])<<8 | uint16(bytes[5]),
		DE_:  uint16(bytes[4])<<8 | uint16(bytes[3]),
		HL_:  uint16(bytes[2])<<8 | uint16(bytes[1]),
		IX:   uint16(bytes[18])<<8 | uint16(bytes[17]),
		IY:   uint16(bytes[16])<<8 | uint16(bytes[15]),
		SP:   uint16(bytes[24])<<8 | uint16(bytes[23]),
		I:    bytes[0],
		R:    bytes[20],
		IM:   bytes[25],
		IFF1: bytes[19]&0x04 != 0,
		IFF2: bytes[19]&0x04 != 0,
	}

	// Load 48k memory
	for i := 16384; i < len(mem.Cells); i++ {
		*mem.Cells[i] = bytes[i-16384+27]
	}

	if len(bytes) == 49179 {
		// 48k model
		state.PC = uint16(*mem.Cells[state.SP+1])<<8 | uint16(*mem.Cells[state.SP])
		state.SP += 2
	} else {
		// 128k model
		mode := bytes[49181]
		mem.PageMode(mode)
		block := 49183
		for bank := 0; bank < 7; bank++ {
			// Skip 3 banks already loaded as 48k memory
			if bank != 2 && bank != 5 && bank != int(mode&0x07) {
				mem.LoadBank(bank, bytes[block:block+16384])
				block += 16384
			}
		}
		state.PC = uint16(bytes[49180])<<8 | uint16(bytes[49179])
	}

	// Restore border colour
	screen.BorderColour(bytes[26], 0)

	// Set CPU state
	cpu.State(state)

	return nil
}

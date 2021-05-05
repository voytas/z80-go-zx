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
func (sna *SNA) Load(file string, cpu *z80.Z80, mem *memory.Memory) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if len(data) == 49179 && len(data) == 131103 && len(data) == 147487 {
		return fmt.Errorf("SNA file format is invalid. Expected 49179, 131103 or 147487 bytes, but received %v", len(data))
	}

	// Restore CPU state
	state := &z80.CPUState{
		AF:   uint16(data[22])<<8 | uint16(data[21]),
		BC:   uint16(data[14])<<8 | uint16(data[13]),
		DE:   uint16(data[12])<<8 | uint16(data[11]),
		HL:   uint16(data[10])<<8 | uint16(data[9]),
		AF_:  uint16(data[8])<<8 | uint16(data[7]),
		BC_:  uint16(data[6])<<8 | uint16(data[5]),
		DE_:  uint16(data[4])<<8 | uint16(data[3]),
		HL_:  uint16(data[2])<<8 | uint16(data[1]),
		IX:   uint16(data[18])<<8 | uint16(data[17]),
		IY:   uint16(data[16])<<8 | uint16(data[15]),
		SP:   uint16(data[24])<<8 | uint16(data[23]),
		I:    data[0],
		R:    data[20],
		IM:   data[25],
		IFF1: data[19]&0x04 != 0,
		IFF2: data[19]&0x04 != 0,
	}

	// Load 48k memory
	for i := 16384; i < len(mem.Cells); i++ {
		*mem.Cells[i] = data[i-16384+27]
	}

	if len(data) == 49179 {
		// 48k model
		state.PC = uint16(*mem.Cells[state.SP+1])<<8 | uint16(*mem.Cells[state.SP])
		state.SP += 2
	} else {
		// 128k model
		mode := data[49181]
		mem.PageMode(mode)
		block := 49183
		for bank := 0; bank < 7; bank++ {
			// Skip 3 banks already loaded as 48k memory
			if bank != 2 && bank != 5 && bank != int(mode&0x07) {
				mem.LoadBank(bank, data[block:block+16384])
				block += 16384
			}
		}
		state.PC = uint16(data[49180])<<8 | uint16(data[49179])
	}

	// Restore border colour
	screen.BorderColour(data[26], 0)

	// Set CPU state
	cpu.State(state)

	return nil
}

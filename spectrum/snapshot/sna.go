package snapshot

import (
	"io/ioutil"

	"github.com/voytas/z80-go-zx/spectrum/screen"
	"github.com/voytas/z80-go-zx/z80"
)

// Loads SNA file to memory and updates the CPU state so it is ready to run
func LoadSNA(filePath string, cpu *z80.Z80, mem []byte) error {
	sna, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Prepare CPU state
	state := z80.CPUState{
		AF:   uint16(sna[22])<<8 | uint16(sna[21]),
		BC:   uint16(sna[14])<<8 | uint16(sna[13]),
		DE:   uint16(sna[12])<<8 | uint16(sna[11]),
		HL:   uint16(sna[10])<<8 | uint16(sna[9]),
		AF_:  uint16(sna[8])<<8 | uint16(sna[7]),
		BC_:  uint16(sna[6])<<8 | uint16(sna[5]),
		DE_:  uint16(sna[4])<<8 | uint16(sna[3]),
		HL_:  uint16(sna[2])<<8 | uint16(sna[1]),
		IX:   uint16(sna[18])<<8 | uint16(sna[17]),
		IY:   uint16(sna[16])<<8 | uint16(sna[15]),
		SP:   uint16(sna[24])<<8 | uint16(sna[23]),
		I:    sna[0],
		R:    sna[20],
		IM:   sna[25],
		IFF1: sna[19]&0x04 != 0,
		IFF2: sna[19]&0x04 != 0,
	}

	// Fill the memory
	for i := 16384; i < len(mem); i++ {
		mem[i] = sna[i-16384+27]
	}

	// Restore border colour
	screen.BorderColour(sna[26], 0)

	// Simulate RETN
	state.PC = uint16(mem[state.SP+1])<<8 | uint16(mem[state.SP])
	state.SP += 2

	// Set CPU state
	cpu.State(&state)

	return nil
}

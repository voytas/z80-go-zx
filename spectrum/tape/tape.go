package tape

import (
	"path/filepath"
	"strings"

	"github.com/voytas/z80-go-zx/spectrum/memory"
	"github.com/voytas/z80-go-zx/z80"
)

const (
	noTapeFile = 0
	tapFile    = 1
	tzxFile    = 2
)

// Type reader
type TapeReader interface {
	NextBlock() *TapeBlock
}

// Represents a tape
type Tape struct {
	reader TapeReader
}

// Represents a tape block
type TapeBlock struct {
	flag     byte   // 00=header, FF=data
	data     []byte // block data
	checksum byte   // checksum
}

// Handles fast loading (block) when load routine is executed.
// Only works if standard ROM routine is used.
func (t *Tape) Load(cpu *z80.Z80, mem *memory.Memory) {
	block := t.reader.NextBlock()
	// No data to load
	if block == nil {
		return
	}

	// Check if running Load or Verify (CF = 1 or CF = 0)
	if cpu.Reg.F_&z80.FC == 0 {
		return
	}

	// Checksum (xor of all bytes)
	checksum := block.flag

	addr := cpu.Reg.IX() // start address
	len := cpu.Reg.DE()  // block length

	// Check expected block type is correct
	if cpu.Reg.A_ == block.flag {
		var i uint16
		for i = 0; i < len; i++ {
			mem.Write(addr+i, block.data[i])
			checksum ^= block.data[i]
		}
		checksum ^= block.checksum
	}

	// Update registers as per normal load and return from the load routine
	cpu.Reg.D, cpu.Reg.E = 0, 0
	cpu.Reg.IXH, cpu.Reg.IXL = byte((addr+len)<<8), byte(addr+len)&0xFF
	cpu.Reg.A = checksum
	cpu.Reg.PC = 0x05E0
}

// Loads a tape file
func (t *Tape) LoadFile(file string) error {
	var err error = nil
	switch getTapeType(file) {
	case tapFile:
		t.reader, err = newTAPReader(file)
	case tzxFile:
		t.reader, err = newTZXReader(file)
	}

	return err
}

// Checks if specified file is *.tap or *.tzx
func (t *Tape) IsTape(file string) bool {
	ft := getTapeType(file)
	return ft == tapFile || ft == tzxFile
}

// Returns the type of the tape file
func getTapeType(file string) int {
	switch strings.ToLower(filepath.Ext(file)) {
	case ".tap":
		return tapFile
	case ".tzx":
		return tzxFile
	default:
		return noTapeFile
	}
}

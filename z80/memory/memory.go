package memory

import "log"

// Memory provides required methods for the emulated memory implementation
type Memory interface {
	// Read memory at the specified address
	Read(addr uint16) byte
	// Write memory at the specified address
	Write(addr uint16, value byte)
}

// BasicMemory provides simple memory implementation supporting ROM and RAM functionality
type BasicMemory struct {
	Cells    []byte
	ramStart uint16 // address of the first RAM cell, lower address is treated as ROM, read-only
}

func (m *BasicMemory) Read(addr uint16) byte {
	// Check if we are reading outside available memory
	if addr >= uint16(len(m.Cells)) {
		return 0xFF
	}
	return m.Cells[addr]
}

func (m *BasicMemory) Write(addr uint16, value byte) {
	// Only write where RAM is available
	if addr >= m.ramStart && addr < uint16(len(m.Cells)) {
		m.Cells[addr] = value
	}
}

func (m *BasicMemory) Configure(cells []byte, ramStart int) {
	if len(cells) > 65536 {
		log.Fatal("Total memory size cannot exceed 65536 byte")
	}
	m.Cells = cells
	m.ramStart = uint16(ramStart)
}

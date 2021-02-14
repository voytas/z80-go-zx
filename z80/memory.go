package z80

import "log"

// Memory provides required methods for the emulated memory implementation
type Memory interface {
	// read memory at the specified address
	read(addr word) byte
	// write memory at the specified address
	write(addr word, value byte)
}

// BasicMemory provides simple memory implementation supporting ROM and RAM functionality
type BasicMemory struct {
	cells    []byte
	ramStart word // address of the first RAM cell, lower address is treated as ROM, read-only
}

func (m *BasicMemory) read(addr word) byte {
	// Check if we are reading outside available memory
	if addr >= word(len(m.cells)) {
		return 0xFF
	}
	return m.cells[addr]
}

func (m *BasicMemory) write(addr word, value byte) {
	// Only write where RAM is available
	if addr >= m.ramStart && addr < word(len(m.cells)) {
		m.cells[addr] = value
	}
}

func (m *BasicMemory) Configure(cells []byte, ramStart int) {
	if len(cells) > 65536 {
		log.Fatal("Total memory size cannot exceed 65536 byte")
	}
	m.cells = cells
	m.ramStart = word(ramStart)
}

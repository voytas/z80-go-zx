package z80

// Memory provides required methods for the emulated memory implementation
type Memory interface {
	// Read memory at the specified address
	read(addr word) *byte
	// Write memory at the specified address
	write(addr word, value byte)
}

// BasicMemory provides simple memory implementation supporting ROM and RAM functionality
type BasicMemory struct {
	cells    []byte
	ramStart word // address of the first RAM cell, lower address is treated as ROM, read-only
}

var _invalid_cell byte = 0xFF

func (m *BasicMemory) read(addr word) *byte {
	// Check if we are reading outside available memory
	if addr >= word(len(m.cells)) {
		return &_invalid_cell
	}
	return &m.cells[addr]
}

func (m *BasicMemory) write(addr word, value byte) {
	// Only write where RAM is available
	if addr >= m.ramStart && addr < word(len(m.cells)) {
		m.cells[addr] = value
	}
}

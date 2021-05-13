package memory

// Memory provides required methods for the emulated memory implementation
type Memory interface {
	// Read memory at the specified address
	Read(addr uint16) byte
	// Write memory at the specified address
	Write(addr uint16, value byte)
}

// BasicMemory provides the most basic memory read/write functionality
type BasicMemory struct {
	Cells []byte
}

func (m *BasicMemory) Read(addr uint16) byte {
	if int(addr) >= len(m.Cells) {
		return 0xFF
	}
	return m.Cells[addr]
}

func (m *BasicMemory) Write(addr uint16, value byte) {
	m.Cells[addr] = value
}

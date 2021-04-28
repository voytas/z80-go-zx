package memory

import (
	"io/ioutil"

	"github.com/voytas/z80-go-zx/z80"
)

// Memory mode: 48k or 128k
const (
	mode48k  = 1
	mode128k = 2
)

// Represents 16k memory bank
type Bank [0x4000]byte

type Memory struct {
	Screen     *Bank         // current screen bank
	Cells      []*byte       // memory as a single array of 65536 bytes
	TC         *z80.TCounter // T state counter
	banks      [8]Bank       // 8 memory banks
	rom48      Bank          // ROM 1 (48k)
	rom128     Bank          // ROM 2 (128k)
	active     [4]*Bank      // currently active banks
	pgDisabled bool          // paging disabled until next reset
	mode       int
}

// Creates a new memory for 48k model
func NewMem48k(romPath string) (*Memory, error) {
	m := &Memory{mode: mode48k}
	err := m.load48ROM(romPath)
	if err != nil {
		return nil, err
	}

	m.Cells = make([]*byte, 0x10000)
	m.copyBank(0, &m.rom48)
	m.copyBank(0x4000, &m.banks[5])
	m.copyBank(0x8000, &m.banks[2])
	m.copyBank(0xC000, &m.banks[0])

	m.pgDisabled = true
	m.Screen = &m.banks[5]

	buildContentionIndex(mode48k)

	return m, nil
}

// Creates a new memory for 128k model with paging
func NewMem128k(rom1Path, rom2Path string) (*Memory, error) {
	m := &Memory{mode: mode128k}
	err := m.load128ROM(rom1Path, rom2Path)
	if err != nil {
		return nil, err
	}

	m.Cells = make([]*byte, 0x10000)
	m.copyBank(0, &m.rom128)
	m.copyBank(0x4000, &m.banks[5])
	m.copyBank(0x8000, &m.banks[2])
	m.copyBank(0xC000, &m.banks[0])

	m.active[0] = &m.rom128   // ROM
	m.active[1] = &m.banks[5] // Screen 1
	m.active[2] = &m.banks[2] // Not pageable
	m.active[3] = &m.banks[0] // RAM 0-7

	m.Screen = m.active[1]

	buildContentionIndex(mode128k)

	return m, nil
}

// Reads a value from the memory address
func (m *Memory) Read(addr uint16) byte {
	m.addContention(addr)
	return *m.Cells[addr]
}

// Writes a value to the memory address
func (m *Memory) Write(addr uint16, value byte) {
	if addr >= 0x4000 && addr <= 0xFFFF {
		*m.Cells[addr] = value
	}
	m.addContention(addr)
}

// Sets the paging mode for 128k model
func (m *Memory) PageMode(mode byte) {
	if m.pgDisabled {
		return
	}

	// Disable paging until next reset
	m.pgDisabled = mode&0b00100000 != 0

	// ROM bank selection
	if mode&0b00010000 != 0 {
		if m.active[0] != &m.rom48 {
			// 48k ROM select
			m.copyBank(0, &m.rom48)
			m.active[0] = &m.rom48
		}
	} else if m.active[0] != &m.rom128 {
		// 128k ROM select
		m.copyBank(0, &m.rom128)
		m.active[0] = &m.rom128
	}

	// Screen bank selection - does not swap memory bank
	if mode&0b00001000 != 0 {
		if m.Screen != &m.banks[7] {
			// second screen select (bank 7)
			m.Screen = &m.banks[7]
		}
	} else if m.Screen != &m.banks[5] {
		// normal screen select (bank 5)
		m.Screen = &m.banks[5]
	}

	// RAM bank selection
	bank := mode & 0x07
	if m.active[3] != &m.banks[bank] {
		m.copyBank(0xC000, &m.banks[bank])
		m.active[3] = &m.banks[bank]
	}
}

// Loads the specified memory bank with data
func (m *Memory) LoadBank(page int, data []byte) {
	for i := 0; i < len(data); i++ {
		m.banks[page][i] = data[i]
	}
}

// Copies the memory bank to the specified address
func (m *Memory) copyBank(addr int, src *Bank) {
	for i := 0; i < len(src); i++ {
		m.Cells[addr+i] = &src[i]
	}
}

// Loads ROMs for 128k model emulation
func (m *Memory) load128ROM(rom1Path, rom2Path string) error {
	rom, err := ioutil.ReadFile(rom1Path)
	if err != nil {
		return err
	}
	for i := 0; i < len(m.rom128); i++ {
		m.rom128[i] = rom[i]
	}

	rom, err = ioutil.ReadFile(rom2Path)
	if err != nil {
		return err
	}
	for i := 0; i < len(m.rom48); i++ {
		m.rom48[i] = rom[i]
	}

	return nil
}

// Loads ROMs for 48k model emulation
func (m *Memory) load48ROM(romPath string) error {
	rom, err := ioutil.ReadFile(romPath)
	if err != nil {
		return err
	}
	for i := 0; i < len(m.rom48); i++ {
		m.rom48[i] = rom[i]
	}

	return nil
}

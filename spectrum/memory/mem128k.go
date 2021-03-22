package memory

import "io/ioutil"

type bank [0x4000]byte

type PageableMemory interface {
	PageMode(mode byte)
}

type Mem128k struct {
	Cells    []*byte
	banks    [8]bank
	rom48    bank
	rom128   bank
	active   [4]*bank
	disabled bool
}

func NewMem128k(rom1Path, rom2Path string) (*Mem128k, error) {
	m := &Mem128k{}
	err := m.loadROM(rom1Path, rom2Path)
	if err != nil {
		return nil, err
	}

	m.Cells = make([]*byte, 0x10000)
	m.copyBank(0, &m.rom128)
	m.copyBank(0x4000, &m.banks[5])
	m.copyBank(0xC000, &m.banks[0])

	m.active[0] = &m.rom128   // ROM 1 or 2
	m.active[1] = &m.banks[5] // Screen 1 or 2
	m.active[2] = &m.banks[2] // Not swappable
	m.active[3] = &m.banks[0] // RAM 0-7

	return m, nil
}

func (m *Mem128k) Read(addr uint16) byte {
	// m.addContendedState(addr)
	return *m.Cells[addr]
}

func (m *Mem128k) Write(addr uint16, value byte) {
	if addr >= ramStart && addr <= 0xFFFF {
		*m.Cells[addr] = value
	}
	//m.addContendedState(addr)
}

func (m *Mem128k) PageMode(mode byte) {
	if m.disabled {
		return
	}

	// Disable paging until next reset
	m.disabled = mode&0b00100000 != 0

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

	// Screen bank selection
	if mode&0b00010000 != 0 {
		if m.active[1] != &m.banks[7] {
			// second screen select (bank 7)
			m.copyBank(0x4000, &m.banks[7])
			m.active[1] = &m.banks[7]
		}
	} else if m.active[1] != &m.banks[5] {
		// normal screen select (bank 5)
		m.copyBank(0x4000, &m.banks[5])
		m.active[1] = &m.banks[5]
	}

	// RAM bank selection
	bank := mode & 0x07
	if m.active[3] != &m.banks[bank] {
		m.copyBank(0xC000, &m.banks[bank])
		m.active[3] = &m.banks[bank]
	}
}

// Copies the memory bank to the specified address
func (m *Mem128k) copyBank(addr int, src *bank) {
	for i := 0; i < len(src); i++ {
		m.Cells[addr+i] = &src[i]
	}
}

func (m *Mem128k) loadROM(rom1Path, rom2Path string) error {
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

package memory

import (
	"io/ioutil"

	"github.com/voytas/z80-go-zx/z80"
)

const ramStart = 0x4000

var contendedStates []byte // index of extra states per each

// Some memory addresses have slower access (extra T states), because of ULA priority
// memory access during the screen draw. So in order to accurately calculate T states
// we need to add extra states when specific memory addresses are accessed.
// https://sinclair.wiki.zxnet.co.uk/wiki/Contended_memory
type Mem48k struct {
	Cells []*byte
	TC    *z80.TCounter
}

// Creates a memory using specified rom file and memory size.
func NewMem48k(romPath string) (*Mem48k, error) {
	mem := &Mem48k{}

	rom, err := ioutil.ReadFile(romPath)
	if err != nil {
		return nil, err
	}

	mem.Cells = make([]*byte, 0x10000)
	for i := 0; i < len(rom); i++ {
		*mem.Cells[i] = rom[i]
	}

	return mem, err
}

func init() {
	createContendedIndex()
}

// For fastest access, create an index of all possible states
// where memory access can be contended. It uses array for
// better performance, although memory usage is higher.
func createContendedIndex() {
	contendedStates = make([]byte, 14559+192*224)
	delays := []byte{6, 5, 4, 3, 2, 1, 0, 0}

	end := 14463
	for t := 14335; t < end; t += len(delays) {
		for di, dt := range delays {
			if t+di < end {
				contendedStates[t+di] = dt
			}
		}
	}
	end = len(contendedStates)
	for t := 14559; t < end; t += len(delays) {
		for di, dt := range delays {
			if t+di < end {
				contendedStates[t+di] = dt
			}
		}
	}
}

// Add extra states if memory address is contended
func (m *Mem48k) addContendedState(addr uint16) {
	if addr >= 0x4000 && addr <= 0x7fff && m.TC.Current < len(contendedStates) {
		m.TC.Add(int(contendedStates[m.TC.Current]))
	}
}

// Read a value from the memory address.
func (m *Mem48k) Read(addr uint16) byte {
	m.addContendedState(addr)
	return *m.Cells[addr]
}

// Write a value to the memory address.
func (m *Mem48k) Write(addr uint16, value byte) {
	if addr >= ramStart && addr <= 0xFFFF {
		*m.Cells[addr] = value
	}
	m.addContendedState(addr)
}

package memory

import (
	"io/ioutil"

	"github.com/voytas/z80-go-zx/spectrum/emulator/state"
)

const ramStart = 0x4000

var ramEnd uint16          // specifies last available RAM address
var contendedStates []byte // index of extra states per each

// Some memory addresses have slower access (extra t-states), because of ULA priority
// memory access during the screen draw. So in order to accurately calculate t-states
// we need to add extra states when specific memory addresses are accessed.
// https://sinclair.wiki.zxnet.co.uk/wiki/Contended_memory
type ContendedMemory struct {
	Cells []byte
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
func addContendedState(addr uint16) {
	if addr >= 0x4000 && addr <= 0x7fff && *state.Current < len(contendedStates) {
		*state.Current += int(contendedStates[*state.Current])
	}
}

// Read a value from the memory address.
func (m *ContendedMemory) Read(addr uint16) byte {
	if addr >= uint16(len(m.Cells)) {
		return 0xFF
	}
	addContendedState(addr)

	return m.Cells[addr]
}

// Write a value to the memory address.
func (m *ContendedMemory) Write(addr uint16, value byte) {
	if addr >= ramStart && addr < ramEnd {
		m.Cells[addr] = value
	}
	addContendedState(addr)
}

// Creates a memory using specified rom file and memory size.
func NewMemory(romPath string, size uint16) (*ContendedMemory, error) {
	mem := &ContendedMemory{}
	mem.Cells = make([]byte, size)

	rom, err := ioutil.ReadFile(romPath)
	if err != nil {
		return nil, err
	}

	copy(mem.Cells, rom)
	ramEnd = uint16(size)

	return mem, err
}

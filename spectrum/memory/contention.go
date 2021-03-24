package memory

// Array containing pre-calculated extra contention T states
// for each applicable T state
var contendedStates []byte

func buildContentionIndex(mode int) {
	var start, states int
	switch mode {
	case mode48k:
		start = 14335 // starting contention state
		states = 224  // number of T states per line
	case mode128k:
		start = 14361 // starting contention state
		states = 228  // number of T states per line
	}

	contendedStates = make([]byte, start+192*states)
	delays := []byte{6, 5, 4, 3, 2, 1, 0, 0}

	for line := 0; line < 192; line++ {
		t := start + line*states
		for i := 0; i < 128; i += len(delays) {
			for di, dt := range delays {
				contendedStates[t+i+di] = dt
			}
		}

	}
}

// Add extra states if memory address is contended
func (m *Memory) addContention(addr uint16) {
	if addr >= 0x4000 && addr <= 0x7fff && m.TC.Current < len(contendedStates) {
		m.TC.Add(int(contendedStates[m.TC.Current]))
	}
}

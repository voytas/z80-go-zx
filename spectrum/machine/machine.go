package machine

type Machine struct {
	Clock           float32 // Clock im MHz
	FrameStates     int     // Number of frames to draw the screen
	ROM1Path        string  // Path to the ROM file 1
	ROM2Path        string  // Path to the ROM file 2 (128k only)
	ContentionTable []byte  // Contention table that provides extra states for given state
}

var ZX48k = &Machine{
	Clock:           3.5,
	FrameStates:     69888,
	ROM1Path:        "./spectrum/rom/48.rom",
	ContentionTable: buildContentionIndex(14335, 224),
}

var ZX128k = &Machine{
	Clock:           3.5469,
	FrameStates:     70908,
	ROM1Path:        "./spectrum/rom/128-0.rom",
	ROM2Path:        "./spectrum/rom/128-1.rom",
	ContentionTable: buildContentionIndex(14361, 228),
}

// Builds the contention table using starting contention state
// and number of T states per line
func buildContentionIndex(start, states int) []byte {
	cs := make([]byte, start+192*states+192)
	delays := []byte{6, 5, 4, 3, 2, 1, 0, 0}

	for line := 0; line < 192; line++ {
		t := start + line*states
		for i := 0; i < 128; i += len(delays) {
			for di, dt := range delays {
				cs[t+i+di] = dt
			}
		}
	}

	return cs
}

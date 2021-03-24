package machine

type Machine struct {
	Clock       float32 // Clock im MHz
	FrameStates int     // Number of frames to draw the screen
	ROM1Path    string  // Path to the ROM file 1
	ROM2Path    string  // Path to the ROM file 2 (128k only)
}

var ZX48k = &Machine{
	Clock:       3.5,
	FrameStates: 69888,
	ROM1Path:    "./spectrum/rom/48.rom",
}

var ZX128k = &Machine{
	Clock:       3.5469,
	FrameStates: 70908,
	ROM1Path:    "./spectrum/rom/128-0.rom",
	ROM2Path:    "./spectrum/rom/128-1.rom",
}

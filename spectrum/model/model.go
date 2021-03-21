package model

type Model struct {
	Memory      uint16  // Total memory in bytes (ROM + RAM)
	Clock       float32 // Clock im MHz
	FrameStates int     // Number of frames to draw the screen
	ROMPath     string  // Path to the ROM file
}

var ZX48k = Model{
	Memory:      0xFFFF,
	Clock:       3.5,
	FrameStates: 69888,
	ROMPath:     "./spectrum/rom/48k.rom",
}

var Current Model = ZX48k

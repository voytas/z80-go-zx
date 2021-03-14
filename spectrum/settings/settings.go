package settings

type Settings struct {
	Memory      uint16 // Total memory in bytes (ROM + RAM)
	FrameStates int    // Number of frames to draw the screen
	ROMPath     string // Path to the ROM file
}

var ZX48k = Settings{
	Memory:      0xFFFF,
	FrameStates: 69888,
	ROMPath:     "./spectrum/rom/48k.rom",
}

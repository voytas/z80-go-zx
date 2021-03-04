package emulator

import "github.com/voytas/z80-go-zx/spectrum/emulator/keyboard"

type iobus struct{}

func (bus *iobus) Read(hi, lo byte) byte {
	if lo == 0xFE {
		return keyboard.GetKeyPortValue(hi)
	}
	return 0xFF
}

func (bus *iobus) Write(hi, lo, data byte) {
}

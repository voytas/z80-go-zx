package emulator

import "github.com/voytas/z80-go-zx/spectrum/emulator/keyboard"

type ioBus struct {
	PortFE byte
}

func (bus *ioBus) Read(hi, lo byte) byte {
	if lo == 0xFE {
		return keyboard.GetKeyPortValue(hi)
	}
	return 0xFF
}

func (bus *ioBus) Write(hi, lo, data byte) {
	if lo == 0xFE {
		bus.PortFE = data
	}
}

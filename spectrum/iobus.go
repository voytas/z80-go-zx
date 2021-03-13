package spectrum

import (
	"github.com/voytas/z80-go-zx/spectrum/keyboard"
	"github.com/voytas/z80-go-zx/spectrum/screen"
)

type ioBus struct {
	tCount *int
}

func (bus *ioBus) Read(hi, lo byte) byte {
	if lo == 0xFE {
		return keyboard.GetKeyPortValue(hi)
	}
	return 0xFF
}

func (bus *ioBus) Write(hi, lo, data byte) {
	if lo == 0xFE {
		screen.AddBorderState(data, *bus.tCount)
	}
}

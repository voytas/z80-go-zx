package spectrum

import (
	"github.com/voytas/z80-go-zx/spectrum/keyboard"
	"github.com/voytas/z80-go-zx/spectrum/screen"
	"github.com/voytas/z80-go-zx/spectrum/sound"
	"github.com/voytas/z80-go-zx/z80"
)

type ioBus struct {
	tc     *z80.TCounter
	beeper *sound.Beeper
}

func NewBus() (*ioBus, error) {
	beeper, err := sound.NewBeeper()
	if err != nil {
		return nil, err
	}

	bus := &ioBus{
		beeper: beeper,
	}

	return bus, nil
}

func (b *ioBus) Read(hi, lo byte) byte {
	if lo == 0xFE {
		return keyboard.GetKeyPortValue(hi)
	}
	return 0xFF
}

func (b *ioBus) Write(hi, lo, data byte) {
	if lo&0x01 == 0 {
		screen.BorderColour(data, b.tc.Current)
		b.beeper.Beep(data, b.tc.Total)
	}
}

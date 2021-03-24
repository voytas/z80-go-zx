package bus

import (
	"github.com/voytas/z80-go-zx/spectrum/keyboard"
	"github.com/voytas/z80-go-zx/spectrum/machine"
	"github.com/voytas/z80-go-zx/spectrum/memory"
	"github.com/voytas/z80-go-zx/spectrum/screen"
	"github.com/voytas/z80-go-zx/spectrum/sound"
	"github.com/voytas/z80-go-zx/z80"
)

type Bus struct {
	tc     *z80.TCounter
	beeper *sound.Beeper
	mem    *memory.Memory
}

func NewBus(machine *machine.Machine, tc *z80.TCounter, mem *memory.Memory) (*Bus, error) {
	beeper, err := sound.NewBeeper(machine.Clock)
	if err != nil {
		return nil, err
	}

	bus := &Bus{
		beeper: beeper,
		mem:    mem,
		tc:     tc,
	}

	return bus, nil
}

func (b *Bus) Read(hi, lo byte) byte {
	if lo == 0xFE {
		return keyboard.GetKeyPortValue(hi)
	}
	return 0xFF
}

func (b *Bus) Write(hi, lo, data byte) {
	if hi&0x80 == 0 && lo&0x02 == 0 {
		// 128k memory port 0x7FFD is decoded, hardware will
		// respond to any port address with bits 1 and 15 reset
		b.mem.PageMode(data)
	} else if lo&0x01 == 0 {
		// ULA
		screen.BorderColour(data, b.tc.Current)
		b.beeper.Beep(data, b.tc.Total)
	}
}

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
	ay     *sound.AY8910
	mem    *memory.Memory
}

func NewBus(machine *machine.Machine, tc *z80.TCounter, mem *memory.Memory) (*Bus, error) {
	beeper, err := sound.NewBeeper(machine.Clock)
	if err != nil {
		return nil, err
	}

	return &Bus{
		beeper: beeper,
		ay:     sound.NewAY8910(),
		mem:    mem,
		tc:     tc,
	}, nil
}

func (b *Bus) Read(hi, lo byte) byte {
	if lo == 0xFE {
		return keyboard.GetKeyPortValue(hi)
	}
	return 0xFF
}

func (b *Bus) Write(hi, lo, data byte) {
	if hi&0x80 == 0 && lo&0x02 == 0 {
		// Memory page select 128k (port 0x7FFD is decoded as: A15=0, A1=0
		b.mem.PageMode(data)
	} else if hi&0xC0 == 0xC0 && lo&0x02 == 0x00 {
		// AY register select (port 0xFFFD is decoded as: A15=1, A14=1, A1=0
		b.ay.SelectReg(data & 0x0F)
	} else if hi&0x80 == 0x80 && lo&0x02 == 0x00 {
		// AY write data (port 0xBFFD is decoded as: A15=1, A1=0
		b.ay.WriteReg(data, b.tc.Total)
	} else if lo&0x01 == 0 {
		// ULA (port 0xFE is decoded as: A0=0)
		screen.BorderColour(data, b.tc.Current)
		b.beeper.Beep(data, b.tc.Total)
	}
}

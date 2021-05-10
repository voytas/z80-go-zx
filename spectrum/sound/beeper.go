package sound

import (
	"io"

	"github.com/hajimehoshi/oto"
)

type Beeper struct {
	ear     byte
	ctx     *oto.Context
	player  *oto.Player
	lastT   int64 // last T state
	lastA   byte  // last Amplitude
	samples chan byte
}

const (
	beeperAmplitudeLo     = 0      // low amplitude value
	beeperAmplitudeHi     = 255    // high amplitude value
	beeperMaxDuration     = 214042 // BEEP x,-60
	beeperStatesPerSample = 8
)

// Create a new instance of the Beeper
func NewBeeper(clock float32) (*Beeper, error) {
	sampleRate := int(clock * 1000000 / beeperStatesPerSample)

	ctx, err := oto.NewContext(sampleRate, 1, 1, 16384)
	if err != nil {
		return nil, err
	}

	b := &Beeper{
		ctx:     ctx,
		player:  ctx.NewPlayer(),
		samples: make(chan byte, sampleRate),
	}

	go func() {
		_, _ = io.Copy(b.player, b)
	}()

	return b, nil
}

func (b *Beeper) Close() {
	b.ctx.Close()
}

func (b *Beeper) Read(buf []byte) (int, error) {
	for i := 0; i < len(buf); i++ {
		select {
		case a := <-b.samples:
			buf[i] = a
		default:
			buf[i] = b.lastA
		}
	}

	return len(buf), nil
}

// Process beeper change at T state
func (b *Beeper) Beep(val byte, t int64) {
	on := false
	ear := val & 0x10
	if b.ear != ear {
		b.ear = ear
		on = true
	}

	if on {
		duration := t - b.lastT
		b.lastT = t
		if duration > 0 && duration <= beeperMaxDuration {
			if b.lastA == beeperAmplitudeLo {
				b.lastA = beeperAmplitudeHi
			} else {
				b.lastA = beeperAmplitudeLo
			}

			length := int(duration / beeperStatesPerSample)
			for i := 0; i < length; i++ {
				b.samples <- b.lastA
			}
		}
	}
}

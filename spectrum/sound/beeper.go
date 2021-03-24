package sound

import (
	"io"

	"github.com/hajimehoshi/oto"
)

type Beeper struct {
	ear    byte
	ctx    *oto.Context
	player *oto.Player
	lastT  int64 // last T state
	lastA  byte  // last Amplitude
}

const (
	amplitudeLo     = 0      // low amplitude value
	amplitudeHi     = 255    // high amplitude value
	maxDuration     = 214042 // BEEP x,-60
	statesPerSample = 8
)

var samples chan byte

// Create a new instance of the Beeper
func NewBeeper(clock float32) (*Beeper, error) {
	sampleRate := int(clock * 1000000 / statesPerSample)

	ctx, err := oto.NewContext(sampleRate, 1, 1, 16384)
	if err != nil {
		return nil, err
	}

	beeper := &Beeper{
		ctx:    ctx,
		player: ctx.NewPlayer(),
	}

	samples = make(chan byte, sampleRate)

	go func() {
		_, _ = io.Copy(beeper.player, beeper)
	}()

	return beeper, nil
}

func (b *Beeper) Close() {
	b.ctx.Close()
}

func (b *Beeper) Read(buf []byte) (int, error) {
	for i := 0; i < len(buf); i++ {
		select {
		case a := <-samples:
			buf[i] = a
		default:
			buf[i] = b.lastA
		}
	}

	return len(buf), nil
}

// Process beeper change at T state
func (b *Beeper) Beep(value byte, t int64) {
	on := false
	ear := value & 0x10
	if b.ear != ear {
		b.ear = ear
		on = true
	}

	if on {
		duration := t - b.lastT
		b.lastT = t
		if duration > 0 && duration <= maxDuration {
			if b.lastA == amplitudeLo {
				b.lastA = amplitudeHi
			} else {
				b.lastA = amplitudeLo
			}

			length := int(duration / statesPerSample)
			for i := 0; i < length; i++ {
				samples <- b.lastA
			}
		}
	}
}

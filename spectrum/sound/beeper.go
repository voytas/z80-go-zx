package sound

import (
	"container/list"
	"io"
	"log"

	"github.com/hajimehoshi/oto"
)

type Beeper struct {
	ear    byte
	ctx    *oto.Context
	on     bool
	count  int
	lastTC int
	queue  *list.List
}

const (
	sampleRate = 44100
)

// Create a new instance of the Beeper
func NewBeeper() (*Beeper, error) {
	ctx, err := oto.NewContext(sampleRate, 1, 1, 4096)
	if err != nil {
		return nil, err
	}

	beeper := &Beeper{
		ctx:   ctx,
		queue: list.New(),
	}

	player := ctx.NewPlayer()
	go func() {
		_, _ = io.Copy(player, beeper)
	}()

	return beeper, nil
}

func (b *Beeper) Close() {
	b.ctx.Close()
}

func (b *Beeper) Read(buf []byte) (int, error) {
	for i := 0; i < len(buf); i++ {
		if b.on {
			buf[i] = 255
		} else {
			buf[i] = 0
		}
	}
	return len(buf), nil
}

func (b *Beeper) Beep(value byte, tc int64) {
	log.Printf("Beep ear on=%v tc=%v", value&0x10, tc)

	//on := false
	ear := value & 0x10
	if b.ear != ear {
		b.ear = ear
		// log.Printf("Beep ear %v %v", ear, tc)
		//on = true
	}

	//b.queue.PushBack()

	// t := tc - b.lastTC
	// if t < 0 {
	// 	t += model.Current.FrameStates
	// }

	// b.on = on
	//b.count = int(model.Current.Clock*1000000/float32(sampleRate)) * t
	//log.Printf("Beep ear on=%v duration=%v", on, t)
}

// 1 frame = 69888 states, 50.08 per s
// sample rate = 44100

// 69888 * 50.08 = 1s = 3499991.04T for 1 s
//  3.5 mil * 1s
// 1s = 44100
// test: 6686 T states
// 79.36487619047619 T states per 1 sample

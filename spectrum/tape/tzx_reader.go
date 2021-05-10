package tape

import (
	"reflect"

	"github.com/voytas/z80-go-zx/spectrum/tape/tzx"
)

type tzxReader struct {
	tzx   *tzx.Tzx
	index int
}

func newTZXReader(file string) (*tzxReader, error) {
	tzx, err := tzx.Load(file)
	if err != nil {
		return nil, err
	}

	return &tzxReader{
		tzx: tzx,
	}, nil
}

func (t *tzxReader) NextBlock() *TapeBlock {
	if t.index >= len(t.tzx.Blocks) {
		return nil
	}

	block := t.tzx.Blocks[t.index]
	for block != nil {
		name := reflect.TypeOf(block).Elem().Name()
		t.index += 1
		switch name {
		case "StandardSpeedDataBlock":
			b := block.(*tzx.StandardSpeedDataBlock)
			return &TapeBlock{
				flag:     b.Data[0],
				data:     b.Data[1 : b.Length-1],
				checksum: b.Data[b.Length-1],
			}
		default:
			if t.index < len(t.tzx.Blocks) {
				block = t.tzx.Blocks[t.index]
			} else {
				return nil
			}
		}
	}

	return nil
}

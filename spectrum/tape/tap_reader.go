package tape

import (
	"io/ioutil"

	"github.com/voytas/z80-go-zx/spectrum/helpers"
)

type tapReader struct {
	reader *helpers.BinaryReader
}

func newTAPReader(file string) (*tapReader, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return &tapReader{
		reader: helpers.NewBinaryReader(data),
	}, nil
}

func (t *tapReader) NextBlock() *TapeBlock {
	len := t.reader.ReadWord()
	data := t.reader.ReadBytes(int(len))
	if t.reader.Eof {
		return nil
	}

	return &TapeBlock{
		flag:     data[0],
		data:     data[1 : len-1],
		checksum: data[len-1],
	}
}

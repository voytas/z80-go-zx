package tape

type tapReader struct {
	reader *BinaryReader
}

func newTAPReader(data []byte) *tapReader {
	return &tapReader{
		reader: NewBinaryReader(data),
	}
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

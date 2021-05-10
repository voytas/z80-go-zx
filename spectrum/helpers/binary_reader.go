package helpers

type BinaryReader struct {
	pos  int
	data []byte
	Eof  bool
}

// Create new instance of the binary reader
func NewBinaryReader(data []byte) *BinaryReader {
	r := &BinaryReader{
		pos:  0,
		data: data,
	}
	return r
}

// Read a single byte value and advance the reader to the next postion
func (r *BinaryReader) ReadByte() byte {
	if len(r.data)-r.pos < 1 {
		r.Eof = true
		return 0
	}
	v := r.data[r.pos]
	r.pos += 1
	return v
}

// Read a single word value and advance the reader to the next postion
func (r *BinaryReader) ReadWord() uint16 {
	if len(r.data)-r.pos < 2 {
		r.Eof = true
		return 0
	}
	v := uint16(r.data[r.pos]) | uint16(r.data[r.pos+1])<<8
	r.pos += 2
	return v
}

// Read a single double word value and advance the reader to the next postion
func (r *BinaryReader) ReadDWord() uint32 {
	if len(r.data)-r.pos < 4 {
		r.Eof = true
		return 0
	}
	v := uint32(r.data[r.pos]) | uint32(r.data[r.pos+1])<<8 | uint32(r.data[r.pos+2])<<16 | uint32(r.data[r.pos+3])<<24
	r.pos += 4
	return v
}

// Read an array of byte values and advance the reader to the next postion
func (r *BinaryReader) ReadBytes(count int) []byte {
	if len(r.data)-r.pos < count {
		r.Eof = true
		return nil
	}
	v := r.data[r.pos : r.pos+count]
	r.pos += count
	return v
}

// Read an array of word values and advance the reader to the next postion
func (r *BinaryReader) ReadWords(count int) []uint16 {
	v := make([]uint16, count)
	for i := 0; i < count*2; i += 2 {
		v[i/2] = uint16(r.data[r.pos+i]) + uint16(r.data[r.pos+i+1])<<8
	}
	r.pos += count * 2
	return v
}

// Read a string value and advance the reader to the next postion
func (r *BinaryReader) ReadString(len int) string {
	v := ""
	for i := 0; i < len; i++ {
		v += string(r.data[i])
	}
	r.pos += len
	return v
}

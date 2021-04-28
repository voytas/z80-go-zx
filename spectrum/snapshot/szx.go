package snapshot

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io/ioutil"
	"log"

	"github.com/voytas/z80-go-zx/spectrum/memory"
	"github.com/voytas/z80-go-zx/spectrum/screen"
	"github.com/voytas/z80-go-zx/z80"
)

const (
	zxstmid_16k       = 0
	zxstmid_48k       = 1
	zxstmid_128k      = 2
	zxstrf_compressed = 1
)

type SZX struct {
	offset int
	bytes  []byte
	state  *z80.CPUState
}

type szxBlock struct {
	id   string
	size int
	data []byte
}

func (szx *SZX) Load(filePath string, cpu *z80.Z80, mem *memory.Memory) error {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	szx.bytes = bytes

	if szx.dwToId(bytes[0:4]) != "ZXST" {
		return errors.New("Not a valid SZX file")
	}

	machine := bytes[6]
	switch machine {
	case zxstmid_16k:
	case zxstmid_48k:
	case zxstmid_128k:
	default:
		return errors.New("Snapshot is for not supported model")
	}

	for {
		block := szx.readNextBlock()
		if block == nil {
			break
		}

		switch block.id {
		case "Z80R":
			szx.processZ80(block)
		case "AY00":
			szx.processAY(block)
		case "RAMP":
			err := szx.processRAMPage(block, mem)
			if err != nil {
				return errors.New("Error reading memory block")
			}
		case "SPCR":
			szx.processULA(block, mem)
		}

		log.Print(block.id)
	}

	cpu.State(szx.state)

	return nil
}

// Read the next block from the snapshot. Returns nil if there are no more blocks.
func (szx *SZX) readNextBlock() *szxBlock {
	if szx.offset == 0 {
		szx.offset = 8
	}

	if szx.offset >= len(szx.bytes) {
		return nil
	}

	size := szx.dwToInt(szx.bytes[szx.offset+4 : szx.offset+8])
	block := &szxBlock{
		id:   szx.dwToId(szx.bytes[szx.offset : szx.offset+4]),
		size: size,
		data: szx.bytes[szx.offset+8 : szx.offset+8+size],
	}
	szx.offset += 8 + size

	return block
}

// Converts DW to string identifier
func (szx *SZX) dwToId(dw []byte) string {
	return string(dw[0]) + string(dw[1]) + string(dw[2]) + string(dw[3])
}

// Converts DW to integer
func (szx *SZX) dwToInt(dw []byte) int {
	return int(dw[0]) + int(dw[1])<<8 + int(dw[2])>>16 + int(dw[3])<<24
}

// ZXSTZ80REGS block
func (szx *SZX) processZ80(block *szxBlock) {
	szx.state = &z80.CPUState{
		AF:   uint16(block.data[0]) | uint16(block.data[1])<<8,
		BC:   uint16(block.data[2]) | uint16(block.data[3])<<8,
		DE:   uint16(block.data[4]) | uint16(block.data[5])<<8,
		HL:   uint16(block.data[6]) | uint16(block.data[7])<<8,
		AF_:  uint16(block.data[8]) | uint16(block.data[9])<<8,
		BC_:  uint16(block.data[10]) | uint16(block.data[11])<<8,
		DE_:  uint16(block.data[12]) | uint16(block.data[13])<<8,
		HL_:  uint16(block.data[14]) | uint16(block.data[15])<<8,
		IX:   uint16(block.data[16]) | uint16(block.data[17])<<8,
		IY:   uint16(block.data[18]) | uint16(block.data[19])<<8,
		SP:   uint16(block.data[20]) | uint16(block.data[21])<<8,
		PC:   uint16(block.data[22]) | uint16(block.data[23])<<8,
		I:    block.data[24],
		R:    block.data[25],
		IFF1: block.data[26] != 0,
		IFF2: block.data[27] != 0,
		IM:   block.data[28],
	}
}

// ZXSTAYBLOCK page
func (szx *SZX) processAY(block *szxBlock) {
	// TODO: set AY state
}

// ZXSTRAMPAGE block
func (szx *SZX) processRAMPage(block *szxBlock, mem *memory.Memory) error {
	page := int(block.data[2])
	if block.data[0]&zxstrf_compressed != 0 {
		r, err := zlib.NewReader(bytes.NewReader(block.data[3:]))
		if err != nil {
			return err
		}
		defer r.Close()
		data, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		mem.LoadBank(page, data)
	} else {
		mem.LoadBank(page, block.data[3:])
	}

	return nil
}

// ZXSTSPECREGS block
func (szx *SZX) processULA(block *szxBlock, mem *memory.Memory) {
	screen.BorderColour(block.data[0], 0)
	mem.PageMode(block.data[1])
}

package exerciser

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/voytas/z80-go-zx/z80"
	"github.com/voytas/z80-go-zx/z80/memory"
)

type ioBus struct{}

func (bus *ioBus) Read(hi, lo byte) byte {
	return 0xFF
}
func (bus *ioBus) Write(hi, lo, data byte) {
	if lo == 5 {
		ch := string(data)
		fmt.Print(ch)
	}
}

func Run(program string) {
	fmt.Printf("Running %s\n", program)
	content, err := ioutil.ReadFile(program)
	if err != nil {
		log.Fatal(err)
	}

	// Initialise 64kB memory with loaded .COM file stating at 100h
	cells := make([]byte, 0xFFFF)
	for i := 0; i < len(content); i++ {
		cells[0x100+i] = content[i]
	}

	// Setup boot
	boot := []byte{
		0xC3, 0x00, 0xF1, // jp 0xF100
		0x00, 0x00, // nop, nop
		0xC3, 0x00, 0xF0, // jp 0xF000
	}
	copy(cells, boot)
	// Setup print routine and self modify address to halt when test is finished since it does jp 0
	// The test code is using BDOS function 2 (C_WRITE) &  BDOS function 9 (C_WRITESTR)
	bdos := []byte{
		0x79,       // ld a,c
		0xFE, 0x02, // cp 2
		0x20, 0x04, // jr nz, +4
		0x7B,       // ld a,e
		0xD3, 0x05, // out(5),a
		0xC9,       // ret
		0xFE, 0x09, // cp 9
		0xC0,       // ret nz
		0x1A,       // ld a,(de)
		0xFE, 0x24, // cp '$'
		0xC8,       // ret z
		0xD3, 0x05, // out(5),a
		0x13,       // inc de
		0x18, 0xF7, // jr, -9
		0xED, 0x59, 0xC9, // out (c),e ret
		0x100:            0x00, // nop - will be modified to halt
		0x21, 0x00, 0xF1, // ld hl,0xF100,
		0x36, 0x76, // ld (hl),0x76 - halt
		0x31, 0x00, 0xF0, // ld sp,0xF000
		0xC3, 0x00, 0x01, // jp 0x100
	}
	for i, b := range bdos {
		cells[0xF000+i] = b
	}

	mem := memory.BasicMemory{
		Cells: cells,
	}
	z80 := z80.NewZ80(&mem)
	z80.IOBus = &ioBus{}
	z80.Run(0)
	fmt.Println("")
}

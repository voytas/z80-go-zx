package exerciser

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/voytas/z80-go-zx/z80"
	"github.com/voytas/z80-go-zx/z80/memory"
)

func Run(program string) {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadFile(root + fmt.Sprintf("/exercises/%s", program))
	if err != nil {
		log.Fatal(err)
	}

	// Initialise 64kB memory with loaded .COM file stating at 100h
	cells := make([]byte, 0xFFFF)
	for i := 0; i < len(content); i++ {
		cells[0x100+i] = content[i]
	}

	// Setup print routine and self modify address to halt when test is finished since it does jp 0
	// The test code is using BDOS function 2 (C_WRITE) &  BDOS function 9 (C_WRITESTR)
	code := []byte{
		0x18, 0x1E, // jr 0x0020; nop
		0x05:       0x79, // ld a,c
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
		0x20: 0x31, 0xFF, 0xFF, // ld sp,0xFFFF
		0x21, 0x00, 0x00, // ld hl,0x0000
		0x36, 0x76, // ld (hl),0x76 - so we halt if jp 0
		0xC3, 0x00, 0x01, // jp 0x0100
	}
	copy(cells, code)

	mem := memory.BasicMemory{}
	mem.Configure(cells, 0)
	cpu := z80.NewCPU(&mem)
	cpu.OUT = func(hi, lo, data byte) {
		if lo == 5 {
			fmt.Print(string(data))
		}
	}
	cpu.Run()

}

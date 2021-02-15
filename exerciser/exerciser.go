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

	// Initialise 4096 kB memory with loaded .COM file stating at 100h
	cells := make([]byte, 0xFFFF)
	for i := 0; i < len(content); i++ {
		cells[0x100+i] = content[i]
	}

	cells[0], cells[1], cells[2] = 0x31, 0xFF, 0xFF // ld sp,0xFFFF
	cells[3], cells[4], cells[5] = 0xC3, 0x00, 0x01 // jp 0x100
	// TODO: Simulate print char
	//cells[5] = 0xC9 // RET

	mem := memory.BasicMemory{}
	mem.Configure(cells, 0x100)
	cpu := z80.NewCPU(&mem)
	cpu.Run()

}

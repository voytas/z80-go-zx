package main

import (
	"fmt"
	"go-zx-go/z80"
)

func main() {
	mem := &z80.Memory{
		Cells:    make([]byte, 0x10000),
		RAMStart: 1000,
	}
	cpu := z80.NewCPU(mem)

	fmt.Printf("value %v", cpu)
}

package exerciser

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/voytas/z80-go-zx/z80"
	"github.com/voytas/z80-go-zx/z80/memory"

	_ "embed"
)

//go:embed boot.com
var boot []byte

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
	// Load selected program
	content, err := ioutil.ReadFile(program)
	if err != nil {
		log.Fatal(err)
	}

	// Load bootstrapper
	// boot, err := ioutil.ReadFile("./boot.com")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Copy loaded program the the execution area
	for i := 0; i < len(content); i++ {
		boot[0x100+i] = content[i]
	}

	mem := memory.BasicMemory{
		Cells: boot,
	}
	z80 := z80.NewZ80(&mem)
	z80.IOBus = &ioBus{}
	z80.Run(0)
	fmt.Println("")
}

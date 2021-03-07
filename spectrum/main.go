package main

import (
	"github.com/voytas/z80-go-zx/spectrum/emulator"
	"github.com/voytas/z80-go-zx/spectrum/emulator/settings"
)

func main() {
	emulator.Run(settings.ZX48k)
}

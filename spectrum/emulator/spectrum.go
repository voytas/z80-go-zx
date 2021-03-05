package emulator

import (
	"io/ioutil"
	"log"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/voytas/z80-go-zx/spectrum/emulator/keyboard"
	"github.com/voytas/z80-go-zx/spectrum/emulator/screen"
	"github.com/voytas/z80-go-zx/z80"
	"github.com/voytas/z80-go-zx/z80/memory"
)

type Emulator struct {
	ioBus ioBus
	z80   *z80.Z80
	mem   []byte
}

func init() {
	runtime.LockOSThread()
}

func (emu *Emulator) Run() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(632, 504, "ZX Spectrum 48k", nil, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize gl bindings:", err)
	}

	window.SetKeyCallback(keyboard.Callback)

	gl.ClearColor(1, 1, 1, 1)
	gl.PixelZoom(4, 4)
	//gl.WindowPos2d(100, 100)

	emu.createSpectrum()

	// ZX Spectrum generates 50 interrupts per second
	ticker := time.NewTicker(20 * time.Millisecond)

	frame := 1
	for !window.ShouldClose() {
		emu.z80.INT(0)
		emu.z80.Run(69888)
		<-ticker.C

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		screen.LastBorderColour = emu.ioBus.PortFE
		scr := screen.RGBA(emu.mem, frame)
		gl.DrawPixels(
			screen.BorderLeft+256+screen.BorderRight,
			screen.BorderTop+192+screen.BorderBottom,
			gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&scr.Pix[0]))

		window.SwapBuffers()
		glfw.PollEvents()

		frame += 1
		if frame > 50 {
			frame = 1
		}
	}
}

func (emu *Emulator) createSpectrum() {
	rom, err := ioutil.ReadFile("./rom/48k.rom")
	if err != nil {
		log.Fatal(err)
	}

	emu.mem = make([]byte, 0xFFFF)
	copy(emu.mem, rom)
	mem := &memory.BasicMemory{}
	mem.Configure(emu.mem, 0x4000)

	emu.ioBus = ioBus{}
	emu.z80 = z80.NewZ80(mem)
	emu.z80.IOBus = &emu.ioBus
}

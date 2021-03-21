package spectrum

import (
	"log"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/voytas/z80-go-zx/spectrum/keyboard"
	"github.com/voytas/z80-go-zx/spectrum/memory"
	"github.com/voytas/z80-go-zx/spectrum/model"
	"github.com/voytas/z80-go-zx/spectrum/screen"
	"github.com/voytas/z80-go-zx/spectrum/snapshot"
	"github.com/voytas/z80-go-zx/z80"
)

type Emulator struct {
	bus *ioBus
	z80 *z80.Z80
	mem *memory.Mem48k
}

func init() {
	runtime.LockOSThread()
}

func Run(model model.Model, fileToLoad string) {
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

	gl.ClearColor(0, 0, 0, 1)
	gl.PixelZoom(4, -4)
	gl.RasterPos2d(-1, 1)

	emu, err := createEmulator(model, fileToLoad)
	if err != nil {
		log.Fatalln("failed to create emulator:", err)
	}

	// 50.08 Hz
	ticker := time.NewTicker(19968 * time.Microsecond)
	for !window.ShouldClose() {
		emu.z80.INT(0xFF)
		emu.z80.Run(model.FrameStates)
		<-ticker.C

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		scr := screen.Render(emu.mem.Cells)
		gl.DrawPixels(
			screen.BorderLeft+256+screen.BorderRight,
			screen.BorderTop+192+screen.BorderBottom,
			gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&scr.Pix[0]))

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func createEmulator(model model.Model, fileToLoad string) (*Emulator, error) {
	// Initialise memory
	mem, err := memory.NewMem48k(model.ROMPath, model.Memory)
	if err != nil {
		return nil, err
	}

	// Initialise IO bus (ports)
	bus, err := NewBus()
	if err != nil {
		return nil, err
	}

	// Initialise new CPU
	cpu := z80.NewZ80(mem)
	cpu.IOBus = bus

	// Share T state counter
	mem.TC = cpu.TC
	bus.tc = cpu.TC

	emu := &Emulator{
		mem: mem,
		bus: bus,
		z80: cpu,
	}

	if fileToLoad != "" {
		err = snapshot.LoadSNA(fileToLoad, emu.z80, mem.Cells)
		if err != nil {
			return nil, err
		}
	}

	return emu, nil
}

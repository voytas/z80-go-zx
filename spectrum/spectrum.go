package spectrum

import (
	"log"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/voytas/z80-go-zx/spectrum/bus"
	"github.com/voytas/z80-go-zx/spectrum/keyboard"
	"github.com/voytas/z80-go-zx/spectrum/machine"
	"github.com/voytas/z80-go-zx/spectrum/memory"
	"github.com/voytas/z80-go-zx/spectrum/screen"
	"github.com/voytas/z80-go-zx/spectrum/snapshot"
	"github.com/voytas/z80-go-zx/z80"
)

type Emulator struct {
	bus *bus.Bus
	z80 *z80.Z80
	mem *memory.Memory
}

func init() {
	runtime.LockOSThread()
}

func Run(machine *machine.Machine, fileToLoad string) {
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

	emu, err := createEmulator(machine, fileToLoad)
	if err != nil {
		log.Fatalln("failed to create emulator:", err)
	}

	freq := machine.Clock * 1000000 / float32(machine.FrameStates)
	ticker := time.NewTicker(time.Duration(1/freq*1000000) * time.Microsecond)
	defer ticker.Stop()

	for !window.ShouldClose() {
		emu.z80.Run(machine.FrameStates)
		emu.z80.INT(0xFF)
		<-ticker.C

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		scr := screen.Render(emu.mem.Screen)
		gl.DrawPixels(
			screen.BorderLeft+256+screen.BorderRight,
			screen.BorderTop+192+screen.BorderBottom,
			gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&scr.Pix[0]))

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func createEmulator(m *machine.Machine, fileToLoad string) (*Emulator, error) {
	// Initialise memory
	var mem *memory.Memory = nil
	var err error
	if m == machine.ZX48k {
		mem, err = memory.NewMem48k(m.ROM1Path)
	} else if m == machine.ZX128k {
		mem, err = memory.NewMem128k(m.ROM1Path, m.ROM2Path)
	} else {
		log.Fatal("Machine not supported")
	}
	if err != nil {
		return nil, err
	}

	// Initialise new CPU
	cpu := z80.NewZ80(mem)

	// Initialise IO bus (ports)
	bus, err := bus.NewBus(m, cpu.TC, mem)
	if err != nil {
		return nil, err
	}

	cpu.IOBus = bus

	// Share T state counter
	mem.TC = cpu.TC

	emu := &Emulator{
		mem: mem,
		bus: bus,
		z80: cpu,
	}

	if fileToLoad != "" {
		err = snapshot.Load(fileToLoad, emu.z80, mem)
		if err != nil {
			return nil, err
		}
	}

	return emu, nil
}

package emulator

import (
	"log"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/voytas/z80-go-zx/spectrum/emulator/keyboard"
	"github.com/voytas/z80-go-zx/spectrum/emulator/memory"
	"github.com/voytas/z80-go-zx/spectrum/emulator/screen"
	"github.com/voytas/z80-go-zx/spectrum/emulator/settings"
	"github.com/voytas/z80-go-zx/spectrum/emulator/snapshot"
	"github.com/voytas/z80-go-zx/spectrum/emulator/state"
	"github.com/voytas/z80-go-zx/z80"
)

type Emulator struct {
	ioBus  ioBus
	z80    *z80.Z80
	mem    *memory.ContendedMemory
	tCount *int
}

func init() {
	runtime.LockOSThread()
}

func Run(settings settings.Settings) {
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

	emu, err := createEmulator(settings)
	if err != nil {
		log.Fatalln("failed to create emulator:", err)
	}

	// ZX Spectrum generates 50 interrupts per second
	ticker := time.NewTicker(20 * time.Millisecond)

	frame := 1
	for !window.ShouldClose() {
		emu.z80.INT(0)
		emu.z80.Run(settings.FrameStates)
		<-ticker.C

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		screen.LastBorderColour = emu.ioBus.PortFE
		scr := screen.RGBA(emu.mem.Cells, frame)
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

func createEmulator(settings settings.Settings) (*Emulator, error) {
	mem, err := memory.NewMemory(settings.ROMPath, settings.Memory)
	if err != nil {
		return nil, err
	}

	emu := &Emulator{
		mem:   mem,
		ioBus: ioBus{},
		z80:   z80.NewZ80(mem),
	}
	emu.z80.IOBus = &emu.ioBus
	emu.tCount = &emu.z80.TCount
	state.Current = emu.tCount

	err = snapshot.LoadSNA("./games/Manic Miner.sna", emu.z80, mem.Cells)
	if err != nil {
		return nil, err
	}

	return emu, nil
}

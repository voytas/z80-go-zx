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

func init() {
	runtime.LockOSThread()
}

func Run() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(612, 494, "ZX Spectrum 48k", nil, nil)
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
	gl.WindowPos2d(100, 100)

	zx, mem := createSpectrum()
	//zx.INT(0)

	// Spectrum generates 50 interrupts per second and we draw our screen then
	ticker := time.NewTicker(20 * time.Millisecond)

	for !window.ShouldClose() {
		zx.INT(0)
		zx.Run(69888)
		<-ticker.C

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		scr := screen.RGBA(mem)
		gl.DrawPixels(256, 192, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&scr.Pix[0]))

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func createSpectrum() (*z80.Z80, []byte) {
	rom, err := ioutil.ReadFile("./rom/48k.rom")
	if err != nil {
		log.Fatal(err)
	}

	cells := make([]byte, 0xFFFF)
	copy(cells, rom)
	mem := &memory.BasicMemory{}
	mem.Configure(cells, 0x4000)

	zx := z80.NewZ80(mem)
	zx.IOBus = &iobus{}

	return zx, cells
}

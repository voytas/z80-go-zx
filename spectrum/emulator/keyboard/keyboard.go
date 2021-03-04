package keyboard

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Index of keyboard ports and values
var ports = map[byte]byte{
	0xFE: 0xFF, // Shift, Z X C V
	0xFD: 0xFF, // A S D F G
	0xFB: 0xFF, // Q W E R T
	0xF7: 0xFF, // 1 2 3 4 5
	0xEF: 0xFF, // 0 9 8 7 6
	0xDF: 0xFF, // P O I U Y
	0xBF: 0xFF, // Enter L K J H
	0x7F: 0xFF, // Space Sym M N B
}

// Index of keys and corresponding mask / port
var keyPorts = map[glfw.Key]struct {
	mask byte
	port byte
}{
	glfw.KeyLeftShift:  {mask: 0b00001, port: 0xFE},
	glfw.KeyZ:          {mask: 0b00010, port: 0xFE},
	glfw.KeyX:          {mask: 0b00100, port: 0xFE},
	glfw.KeyC:          {mask: 0b01000, port: 0xFE},
	glfw.KeyV:          {mask: 0b10000, port: 0xFE},
	glfw.KeyA:          {mask: 0b00001, port: 0xFD},
	glfw.KeyS:          {mask: 0b00010, port: 0xFD},
	glfw.KeyD:          {mask: 0b00100, port: 0xFD},
	glfw.KeyF:          {mask: 0b01000, port: 0xFD},
	glfw.KeyG:          {mask: 0b10000, port: 0xFD},
	glfw.KeyQ:          {mask: 0b00001, port: 0xFB},
	glfw.KeyW:          {mask: 0b00010, port: 0xFB},
	glfw.KeyE:          {mask: 0b00100, port: 0xFB},
	glfw.KeyR:          {mask: 0b01000, port: 0xFB},
	glfw.KeyT:          {mask: 0b10000, port: 0xFB},
	glfw.Key1:          {mask: 0b00001, port: 0xF7},
	glfw.Key2:          {mask: 0b00010, port: 0xF7},
	glfw.Key3:          {mask: 0b00100, port: 0xF7},
	glfw.Key4:          {mask: 0b01000, port: 0xF7},
	glfw.Key5:          {mask: 0b10000, port: 0xF7},
	glfw.Key0:          {mask: 0b00001, port: 0xEF},
	glfw.Key9:          {mask: 0b00010, port: 0xEF},
	glfw.Key8:          {mask: 0b00100, port: 0xEF},
	glfw.Key7:          {mask: 0b01000, port: 0xEF},
	glfw.Key6:          {mask: 0b10000, port: 0xEF},
	glfw.KeyP:          {mask: 0b00001, port: 0xDF},
	glfw.KeyO:          {mask: 0b00010, port: 0xDF},
	glfw.KeyI:          {mask: 0b00100, port: 0xDF},
	glfw.KeyU:          {mask: 0b01000, port: 0xDF},
	glfw.KeyY:          {mask: 0b10000, port: 0xDF},
	glfw.KeyEnter:      {mask: 0b00001, port: 0xBF},
	glfw.KeyL:          {mask: 0b00010, port: 0xBF},
	glfw.KeyK:          {mask: 0b00100, port: 0xBF},
	glfw.KeyJ:          {mask: 0b01000, port: 0xBF},
	glfw.KeyH:          {mask: 0b10000, port: 0xBF},
	glfw.KeySpace:      {mask: 0b00001, port: 0x7F},
	glfw.KeyRightShift: {mask: 0b00010, port: 0x7F},
	glfw.KeyM:          {mask: 0b00100, port: 0x7F},
	glfw.KeyN:          {mask: 0b01000, port: 0x7F},
	glfw.KeyB:          {mask: 0b10000, port: 0x7F},
}

// Returns a status of the keys for the specific port
func GetKeyPortValue(port byte) byte {
	val, ok := ports[port]
	if !ok {
		return 0xFF
	}
	return val
}

// OpenGL keyboard callback
func Callback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	kp, ok := keyPorts[key]
	if !ok {
		return
	}

	switch action {
	case glfw.Press:
		ports[kp.port] &= ^kp.mask
	case glfw.Release:
		ports[kp.port] |= kp.mask
	}
}

package keyboard

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Key half-rows and key statuses
var ports = []byte{
	0x01: 0xFF, // Shift Z X C V
	0x02: 0xFF, // A S D F G
	0x04: 0xFF, // Q W E R T
	0x08: 0xFF, // 1 2 3 4 5
	0x10: 0xFF, // 0 9 8 7 6
	0x20: 0xFF, // P O I U Y
	0x40: 0xFF, // Enter L K J H
	0x80: 0xFF, // Space Sym M N B
}

// Index of keys and corresponding mask / port
var keyPorts = map[glfw.Key]struct {
	mask byte
	port byte
}{
	glfw.KeyLeftShift:  {mask: 0b00001, port: 0x01},
	glfw.KeyZ:          {mask: 0b00010, port: 0x01},
	glfw.KeyX:          {mask: 0b00100, port: 0x01},
	glfw.KeyC:          {mask: 0b01000, port: 0x01},
	glfw.KeyV:          {mask: 0b10000, port: 0x01},
	glfw.KeyA:          {mask: 0b00001, port: 0x02},
	glfw.KeyS:          {mask: 0b00010, port: 0x02},
	glfw.KeyD:          {mask: 0b00100, port: 0x02},
	glfw.KeyF:          {mask: 0b01000, port: 0x02},
	glfw.KeyG:          {mask: 0b10000, port: 0x02},
	glfw.KeyQ:          {mask: 0b00001, port: 0x04},
	glfw.KeyW:          {mask: 0b00010, port: 0x04},
	glfw.KeyE:          {mask: 0b00100, port: 0x04},
	glfw.KeyR:          {mask: 0b01000, port: 0x04},
	glfw.KeyT:          {mask: 0b10000, port: 0x04},
	glfw.Key1:          {mask: 0b00001, port: 0x08},
	glfw.Key2:          {mask: 0b00010, port: 0x08},
	glfw.Key3:          {mask: 0b00100, port: 0x08},
	glfw.Key4:          {mask: 0b01000, port: 0x08},
	glfw.Key5:          {mask: 0b10000, port: 0x08},
	glfw.Key0:          {mask: 0b00001, port: 0x10},
	glfw.Key9:          {mask: 0b00010, port: 0x10},
	glfw.Key8:          {mask: 0b00100, port: 0x10},
	glfw.Key7:          {mask: 0b01000, port: 0x10},
	glfw.Key6:          {mask: 0b10000, port: 0x10},
	glfw.KeyP:          {mask: 0b00001, port: 0x20},
	glfw.KeyO:          {mask: 0b00010, port: 0x20},
	glfw.KeyI:          {mask: 0b00100, port: 0x20},
	glfw.KeyU:          {mask: 0b01000, port: 0x20},
	glfw.KeyY:          {mask: 0b10000, port: 0x20},
	glfw.KeyEnter:      {mask: 0b00001, port: 0x40},
	glfw.KeyL:          {mask: 0b00010, port: 0x40},
	glfw.KeyK:          {mask: 0b00100, port: 0x40},
	glfw.KeyJ:          {mask: 0b01000, port: 0x40},
	glfw.KeyH:          {mask: 0b10000, port: 0x40},
	glfw.KeySpace:      {mask: 0b00001, port: 0x80},
	glfw.KeyRightShift: {mask: 0b00010, port: 0x80},
	glfw.KeyM:          {mask: 0b00100, port: 0x80},
	glfw.KeyN:          {mask: 0b01000, port: 0x80},
	glfw.KeyB:          {mask: 0b10000, port: 0x80},
}

// Returns a status of the keys for the specific port.
// Port can also specify any key, for example if checking port 0x02
// it means any key except A-G, some games use this trick.
func GetKeyPortValue(port byte) byte {
	val := byte(0xFF)
	port = ^port
	for _, p := range []byte{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80} {
		if port&p == p {
			val &= ports[p]
		}
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

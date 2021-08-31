package keyboard

import (
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	KEY_NONE = iota
	KEY_SHIFT
	KEY_SYMBOL
	KEY_ENTER
	KEY_SPACE
	KEY_1
	KEY_2
	KEY_3
	KEY_4
	KEY_5
	KEY_6
	KEY_7
	KEY_8
	KEY_9
	KEY_0
	KEY_A
	KEY_B
	KEY_C
	KEY_D
	KEY_E
	KEY_F
	KEY_G
	KEY_H
	KEY_I
	KEY_J
	KEY_K
	KEY_L
	KEY_M
	KEY_N
	KEY_O
	KEY_P
	KEY_Q
	KEY_R
	KEY_S
	KEY_T
	KEY_U
	KEY_V
	KEY_W
	KEY_X
	KEY_Y
	KEY_Z
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

// Simulate key down with a delay
func KeyDown(key byte, delay time.Duration) {
	handleKey(key, true)
	time.Sleep(delay * time.Millisecond)
}

// Simulate key up with a delay
func KeyUp(key byte, delay time.Duration) {
	handleKey(key, false)
	time.Sleep(delay * time.Millisecond)
}

// Simulate key down and up with a delay
func KeyDownUp(key byte, delay time.Duration) {
	handleKey(key, true)
	time.Sleep(delay * time.Millisecond)
	handleKey(key, false)
}

func handleKey(key byte, down bool) {
	var update = func(port byte, mask byte) {
		if down {
			ports[port] &= mask // key down
		} else {
			ports[port] |= ^mask // key up
		}
	}

	switch key {
	case KEY_NONE:
		ports[0x01], ports[0x02], ports[0x04], ports[0x08] = 0xFF, 0xFF, 0xFF, 0xFF
		ports[0x10], ports[0x20], ports[0x40], ports[0x80] = 0xFF, 0xFF, 0xFF, 0xFF
	case KEY_SHIFT:
		update(0x01, 0b11111110)
	case KEY_SYMBOL:
		update(0x80, 0b11111101)
	case KEY_ENTER:
		update(0x40, 0b11111110)
	case KEY_SPACE:
		update(0x80, 0b11111110)
	case KEY_1:
		update(0x08, 0b11111110)
	case KEY_2:
		update(0x08, 0b11111101)
	case KEY_3:
		update(0x08, 0b11111011)
	case KEY_4:
		update(0x08, 0b11110111)
	case KEY_5:
		update(0x08, 0b11101111)
	case KEY_6:
		update(0x10, 0b11101111)
	case KEY_7:
		update(0x10, 0b11110111)
	case KEY_8:
		update(0x10, 0b11111011)
	case KEY_9:
		update(0x10, 0b11111101)
	case KEY_0:
		update(0x10, 0b11111110)
	case KEY_A:
		update(0x02, 0b11111110)
	case KEY_B:
		update(0x80, 0b11101111)
	case KEY_C:
		update(0x01, 0b11110111)
	case KEY_D:
		update(0x02, 0b11111011)
	case KEY_E:
		update(0x04, 0b11111011)
	case KEY_F:
		update(0x02, 0b11110111)
	case KEY_G:
		update(0x02, 0b11101111)
	case KEY_H:
		update(0x40, 0b11101111)
	case KEY_I:
		update(0x20, 0b11111011)
	case KEY_J:
		update(0x40, 0b11110111)
	case KEY_K:
		update(0x40, 0b11111011)
	case KEY_L:
		update(0x40, 0b11111101)
	case KEY_M:
		update(0x80, 0b11111011)
	case KEY_N:
		update(0x80, 0b11110111)
	case KEY_O:
		update(0x20, 0b11111101)
	case KEY_P:
		update(0x20, 0b11111110)
	case KEY_Q:
		update(0x04, 0b11111110)
	case KEY_R:
		update(0x04, 0b11110111)
	case KEY_S:
		update(0x02, 0b11111101)
	case KEY_T:
		update(0x04, 0b11101111)
	case KEY_U:
		update(0x20, 0b11110111)
	case KEY_V:
		update(0x01, 0b11101111)
	case KEY_W:
		update(0x04, 0b11111101)
	case KEY_X:
		update(0x01, 0b11111011)
	case KEY_Y:
		update(0x20, 0b11101111)
	case KEY_Z:
		update(0x01, 0b11111101)
	}
}

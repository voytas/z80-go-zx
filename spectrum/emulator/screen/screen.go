package screen

import (
	"image"
)

const (
	BorderTop    = 30
	BorderBottom = 30
	BorderRight  = 30
	BorderLeft   = 30
)

var LastBorderColour byte

// Address of each line on the screen (0-191), it is not linear
var lines = []int{
	0x4000, 0x4100, 0x4200, 0x4300, 0x4400, 0x4500, 0x4600, 0x4700, // Lines 0-7
	0x4020, 0x4120, 0x4220, 0x4320, 0x4420, 0x4520, 0x4620, 0x4720, // Lines 8-15
	0x4040, 0x4140, 0x4240, 0x4340, 0x4440, 0x4540, 0x4640, 0x4740, // Lines 16-23
	0x4060, 0x4160, 0x4260, 0x4360, 0x4460, 0x4560, 0x4660, 0x4760, // Lines 24-31
	0x4080, 0x4180, 0x4280, 0x4380, 0x4480, 0x4580, 0x4680, 0x4780, // Lines 32-39
	0x40A0, 0x41A0, 0x42A0, 0x43A0, 0x44A0, 0x45A0, 0x46A0, 0x47A0, // Lines 40-47
	0x40C0, 0x41C0, 0x42C0, 0x43C0, 0x44C0, 0x45C0, 0x46C0, 0x47C0, // Lines 48-55
	0x40E0, 0x41E0, 0x42E0, 0x43E0, 0x44E0, 0x45E0, 0x46E0, 0x47E0, // Lines 56-63
	0x4800, 0x4900, 0x4A00, 0x4B00, 0x4C00, 0x4D00, 0x4E00, 0x4F00, // Lines 64-71
	0x4820, 0x4920, 0x4A20, 0x4B20, 0x4C20, 0x4D20, 0x4E20, 0x4F20, // Lines 72-79
	0x4840, 0x4940, 0x4A40, 0x4B40, 0x4C40, 0x4D40, 0x4E40, 0x4F40, // Lines 80-87
	0x4860, 0x4960, 0x4A60, 0x4B60, 0x4C60, 0x4D60, 0x4E60, 0x4F60, // Lines 88-95
	0x4880, 0x4980, 0x4A80, 0x4B80, 0x4C80, 0x4D80, 0x4E80, 0x4F80, // Lines 96-103
	0x48A0, 0x49A0, 0x4AA0, 0x4BA0, 0x4CA0, 0x4DA0, 0x4EA0, 0x4FA0, // Lines 104-111
	0x48C0, 0x49C0, 0x4AC0, 0x4BC0, 0x4CC0, 0x4DC0, 0x4EC0, 0x4FC0, // Lines 112-119
	0x48E0, 0x49E0, 0x4AE0, 0x4BE0, 0x4CE0, 0x4DE0, 0x4EE0, 0x4FE0, // Lines 120-127
	0x5000, 0x5100, 0x5200, 0x5300, 0x5400, 0x5500, 0x5600, 0x5700, // Lines 128-135
	0x5020, 0x5120, 0x5220, 0x5320, 0x5420, 0x5520, 0x5620, 0x5720, // Lines 136-143
	0x5040, 0x5140, 0x5240, 0x5340, 0x5440, 0x5540, 0x5640, 0x5740, // Lines 144-151
	0x5060, 0x5160, 0x5260, 0x5360, 0x5460, 0x5560, 0x5660, 0x5760, // Lines 152-159
	0x5080, 0x5180, 0x5280, 0x5380, 0x5480, 0x5580, 0x5680, 0x5780, // Lines 160-167
	0x50A0, 0x51A0, 0x52A0, 0x53A0, 0x54A0, 0x55A0, 0x56A0, 0x57A0, // Lines 168-175
	0x50C0, 0x51C0, 0x52C0, 0x53C0, 0x54C0, 0x55C0, 0x56C0, 0x57C0, // Lines 176-183
	0x50E0, 0x51E0, 0x52E0, 0x53E0, 0x54E0, 0x55E0, 0x56E0, 0x57E0, // Lines 184-191
}

var bits = []byte{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01}

// Indexed array of ink colours
var inkColours = [][]byte{
	// Normal
	0b0000000: {0x00, 0x00, 0x00}, // Black
	0b0000001: {0x00, 0x00, 0xD7}, // Blue
	0b0000010: {0xD7, 0x00, 0x00}, // Red
	0b0000011: {0xD7, 0x00, 0xD7}, // Magenta
	0b0000100: {0x00, 0xD7, 0x00}, // Green
	0b0000101: {0x00, 0xD7, 0xD7}, // Cyan
	0b0000110: {0xD7, 0xD7, 0x00}, // Yellow
	0b0000111: {0xD7, 0xD7, 0xD7}, // White
	// Bright
	0b1000000: {0x00, 0x00, 0x00}, // Black
	0b1000001: {0x00, 0x00, 0xFF}, // Blue
	0b1000010: {0xFF, 0x00, 0x00}, // Red
	0b1000011: {0xFF, 0x00, 0xFF}, // Magenta
	0b1000100: {0x00, 0xFF, 0x00}, // green
	0b1000101: {0x00, 0xFF, 0xFF}, // cyan
	0b1000110: {0xFF, 0xFF, 0x00}, // yellow
	0b1000111: {0xFF, 0xFF, 0xFF}, // white
}

var borderColours = inkColours

// Indexed array of paper colours
var paperColours = [][]byte{
	// Normal
	0b0000000: {0x00, 0x00, 0x00}, // Black
	0b0001000: {0x00, 0x00, 0xD7}, // Blue
	0b0010000: {0xD7, 0x00, 0x00}, // Red
	0b0011000: {0xD7, 0x00, 0xD7}, // Magenta
	0b0100000: {0x00, 0xD7, 0x00}, // Green
	0b0101000: {0x00, 0xD7, 0xD7}, // Cyan
	0b0110000: {0xD7, 0xD7, 0x00}, // Yellow
	0b0111000: {0xD7, 0xD7, 0xD7}, // White
	// Bright
	0b1000000: {0x00, 0x00, 0x00}, // Black
	0b1001000: {0x00, 0x00, 0xFF}, // Blue
	0b1010000: {0xFF, 0x00, 0x00}, // Red
	0b1011000: {0xFF, 0x00, 0xFF}, // Magenta
	0b1100000: {0x00, 0xFF, 0x00}, // green
	0b1101000: {0x00, 0xFF, 0xFF}, // cyan
	0b1110000: {0xFF, 0xFF, 0x00}, // yellow
	0b1111000: {0xFF, 0xFF, 0xFF}, // white
}

// Renders the screen as RGBA image (upside down e.g. bottom left to top right)
func RGBA(mem []byte, frame int) *image.RGBA {
	width := BorderLeft + 256 + BorderRight
	height := BorderTop + 192 + BorderBottom
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	border := borderColours[LastBorderColour&0x07]

	// Main screen
	for line, addr := range lines {
		pixel := 4 * (width*(191-line+BorderBottom) + BorderLeft)

		// Left border
		for left := pixel - 4*BorderLeft; left < pixel; left += 4 {
			img.Pix[left] = border[0]
			img.Pix[left+1] = border[1]
			img.Pix[left+2] = border[2]
			img.Pix[left+3] = 0xFF
		}
		// Right border
		for right := pixel + 4*256; right < pixel+4*256+4*BorderLeft; right += 4 {
			img.Pix[right] = border[0]
			img.Pix[right+1] = border[1]
			img.Pix[right+2] = border[2]
			img.Pix[right+3] = 0xFF
		}

		for col := 0; col < 32; col++ {
			attr := mem[0x5800+32*(line/8)+col]
			cell := mem[addr+col]
			for _, bit := range bits {
				var colour []byte
				flash := attr&0x80 != 0 && frame >= 32
				on := cell&bit != 0
				if on != flash {
					colour = inkColours[attr&0b01000111]
				} else {
					colour = paperColours[attr&0b01111000]
				}
				img.Pix[pixel] = colour[0]
				img.Pix[pixel+1] = colour[1]
				img.Pix[pixel+2] = colour[2]
				img.Pix[pixel+3] = 0xFF
				pixel += 4
			}
		}
	}

	// Top border
	for pixel := 4 * width * (BorderBottom + 192); pixel < 4*width*(height); pixel += 4 {
		img.Pix[pixel] = border[0]
		img.Pix[pixel+1] = border[1]
		img.Pix[pixel+2] = border[2]
		img.Pix[pixel+3] = 0xFF
	}

	// Bottom border
	for pixel := 0; pixel < 4*width*BorderBottom; pixel += 4 {
		img.Pix[pixel] = border[0]
		img.Pix[pixel+1] = border[1]
		img.Pix[pixel+2] = border[2]
		img.Pix[pixel+3] = 0xFF
	}

	return img
}
package screen

import (
	"image"
)

var frame = 1                                            // current frame count
var img = image.NewRGBA(image.Rect(0, 0, width, height)) // rendered screen object

func init() {
	// Set alpha to FF, it won't change
	for px := 3; px < len(img.Pix); px += 4 {
		img.Pix[px] = 0xFF
	}
}

// Renders the screen as RGBA image
func Render(mem []byte) *image.RGBA {
	// Top border
	for pixel := 0; pixel < 4*width*BorderTop; pixel += 4 {
		border := findBorderColour(pixelT[pixel/4])
		img.Pix[pixel] = border[0]
		img.Pix[pixel+1] = border[1]
		img.Pix[pixel+2] = border[2]
	}

	// Main screen and left/right border
	for line, addr := range lines {
		px := 4 * (width*(line+BorderTop) + BorderLeft)

		// Left border
		for left := px - 4*BorderLeft; left < px; left += 4 {
			border := findBorderColour(pixelT[left/4])
			img.Pix[left] = border[0]
			img.Pix[left+1] = border[1]
			img.Pix[left+2] = border[2]
		}
		// Right border
		for right := px + 4*256; right < px+4*256+4*BorderRight; right += 4 {
			border := findBorderColour(pixelT[right/4])
			img.Pix[right] = border[0]
			img.Pix[right+1] = border[1]
			img.Pix[right+2] = border[2]
		}

		// Centre
		for col := 0; col < 32; col++ {
			attr := mem[0x5800+32*(line/8)+col]
			cell := mem[addr+col]
			for _, bit := range []byte{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01} {
				var colour []byte
				flash := attr&0x80 != 0 && frame >= 32
				on := cell&bit != 0
				if on != flash {
					colour = inkPalette[attr&0b01000111]
				} else {
					colour = paperPalette[attr&0b01111000]
				}
				img.Pix[px] = colour[0]
				img.Pix[px+1] = colour[1]
				img.Pix[px+2] = colour[2]
				px += 4
			}
		}
	}

	// Bottom border
	for px := 4 * width * (BorderTop + 192); px < 4*width*(height); px += 4 {
		border := findBorderColour(pixelT[px/4])
		img.Pix[px] = border[0]
		img.Pix[px+1] = border[1]
		img.Pix[px+2] = border[2]
	}

	// Can safely drop recorded states
	resetBorderStates()

	// Keep frame count for the "flash" attribute
	frame += 1
	if frame > 50 {
		frame = 1
	}

	return img
}

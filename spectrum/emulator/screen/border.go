package screen

const (
	BorderTop    = 30  // top border height (max 64)
	BorderBottom = 30  // bottom border height (max 56)
	BorderLeft   = 30  // left border width (max 48)
	BorderRight  = 30  // right border width (max 48)
	tPerLine     = 224 // number of states to render single screen line
)

var (
	borderStates    []*borderState
	nextBorderState *borderState
	lastBorderState borderState
	width           = BorderLeft + 256 + BorderRight
	height          = BorderTop + 192 + BorderBottom
	pixelT          []int // T state for each screen pixel
)

type borderState struct {
	colour byte
	tCount int
	index  int
}

func init() {
	// Initialise array containing each pixel T state value for quick access
	pixelT = make([]int, width*height)
	for line := 0; line < height; line++ {
		for px := 0; px < width; px++ {
			pixelT[line*width+px] = (64-BorderTop+line)*tPerLine + (48-BorderLeft+px)/2
		}
	}
}

func AddBorderState(colour byte, tCount int) {
	colour &= 0x07
	if lastBorderState.colour != colour {
		lastBorderState.colour = colour
		borderStates = append(borderStates, &borderState{
			colour: colour,
			tCount: tCount,
			index:  len(borderStates),
		})

	}
}

func findBorderColour(t int) []byte {
	if len(borderStates) == 0 {
		return borderPalette[lastBorderState.colour]
	}

	start := 0
	if nextBorderState != nil {
		if t < nextBorderState.tCount {
			return borderPalette[lastBorderState.colour]
		}
		start = nextBorderState.index
	}

	var state *borderState
	for i := start; i < len(borderStates); i++ {
		if borderStates[i].tCount > t {
			nextBorderState = borderStates[i]
			break
		}
		state = borderStates[i]
	}

	if state != nil {
		lastBorderState = *state
	}

	return borderPalette[lastBorderState.colour]
}

func resetBorderStates() {
	if len(borderStates) > 0 {
		lastBorderState = *borderStates[len(borderStates)-1]
	}
	borderStates = nil
	nextBorderState = nil
}

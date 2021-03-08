package screen

const (
	BorderTop    = 30
	BorderBottom = 30
	BorderRight  = 30
	BorderLeft   = 30
)

var lastBorderColour byte
var borderStates []*borderState

type borderState struct {
	colour byte
	tCount int
}

func AddBorderState(colour byte, tCount int) {
	lastBorderColour = colour & 0x07
	borderStates = append(borderStates, &borderState{
		colour: colour,
		tCount: tCount,
	})
}

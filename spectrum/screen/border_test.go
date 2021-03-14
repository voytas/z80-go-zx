package screen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findBorderColour(t *testing.T) {
	lastBorderState.colour = 5

	for tCount := range []int{0, 3847, 5678, 63626} {
		c := findBorderColour(tCount)
		assert.Equal(t, c, borderPalette[5])
	}

	borderStates = append(borderStates, &borderState{colour: 2, tCount: 4333})
	c := findBorderColour(3847)
	assert.Equal(t, c, borderPalette[5])
	c = findBorderColour(4332)
	assert.Equal(t, c, borderPalette[5])
	c = findBorderColour(4333)
	assert.Equal(t, c, borderPalette[2])
	c = findBorderColour(13222)
	assert.Equal(t, c, borderPalette[2])

	nextBorderState = nil
	lastBorderState.colour = 5
	borderStates = append(borderStates, &borderState{colour: 1, tCount: 7363, index: 1})
	c = findBorderColour(3847)
	assert.Equal(t, c, borderPalette[5])
	c = findBorderColour(4332)
	assert.Equal(t, c, borderPalette[5])
	c = findBorderColour(4333)
	assert.Equal(t, c, borderPalette[2])
	c = findBorderColour(4334)
	assert.Equal(t, c, borderPalette[2])
	c = findBorderColour(7362)
	assert.Equal(t, c, borderPalette[2])
	c = findBorderColour(7363)
	assert.Equal(t, c, borderPalette[1])
	c = findBorderColour(7364)
	assert.Equal(t, c, borderPalette[1])

	nextBorderState = nil
	lastBorderState.colour = 5
	borderStates = append(borderStates, &borderState{colour: 3, tCount: 18222, index: 2})
	c = findBorderColour(3847)
	assert.Equal(t, c, borderPalette[5])
	c = findBorderColour(4332)
	assert.Equal(t, c, borderPalette[5])
	c = findBorderColour(4333)
	assert.Equal(t, c, borderPalette[2])
	c = findBorderColour(4334)
	assert.Equal(t, c, borderPalette[2])
	c = findBorderColour(7362)
	assert.Equal(t, c, borderPalette[2])
	c = findBorderColour(7363)
	assert.Equal(t, c, borderPalette[1])
	c = findBorderColour(7364)
	assert.Equal(t, c, borderPalette[1])
	c = findBorderColour(18222)
	assert.Equal(t, c, borderPalette[3])
}

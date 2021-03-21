package z80

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Add(t *testing.T) {
	tc := TCounter{}

	tc.Add(23)
	assert.EqualValues(t, 23, tc.Total)
	assert.EqualValues(t, 23, tc.Current)
	assert.Equal(t, false, tc.done())
}

func Test_Limit(t *testing.T) {
	tc := TCounter{}
	tc.limit(41)

	tc.Add(12)
	assert.EqualValues(t, 12, tc.Total)
	assert.EqualValues(t, 12, tc.Current)

	tc.Add(30)
	assert.EqualValues(t, 42, tc.Total)
	assert.EqualValues(t, 42, tc.Current)
	assert.Equal(t, true, tc.done())

	tc.limit(20)
	assert.EqualValues(t, 42, tc.Total)
	assert.EqualValues(t, 0, tc.Current)

	tc.Add(19)
	assert.EqualValues(t, 61, tc.Total)
	assert.EqualValues(t, 19, tc.Current)
	assert.Equal(t, true, tc.done())
}

func Test_Halt(t *testing.T) {
	tc := TCounter{}
	tc.limit(41)

	tc.halt()
	assert.EqualValues(t, 44, tc.Total)
	assert.EqualValues(t, 44, tc.Current)
	assert.Equal(t, true, tc.done())
}

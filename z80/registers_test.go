package z80

// func TestBC(t *testing.T) {
// 	r := registers{}
// 	assert.Equal(t, word(0x0000), r.BC())

// 	r = registers{C: 0x01}
// 	assert.Equal(t, word(0x0001), r.BC())

// 	r = registers{B: 0x01}
// 	assert.Equal(t, word(0x0100), r.BC())

// 	r = registers{B: 0xFF}
// 	assert.Equal(t, word(0xFF00), r.BC())

// 	r = registers{C: 0xFF}
// 	assert.Equal(t, word(0x00FF), r.BC())

// 	r = registers{B: 0x12, C: 0xF7}
// 	assert.Equal(t, word(0x12F7), r.BC())
// }

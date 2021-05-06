package tzx

// Converts little endian array to integer
func bytesToInt(data []byte) int {
	v := 0
	for i, b := range data {
		v += int(b) << (i * 8)
	}
	return v
}

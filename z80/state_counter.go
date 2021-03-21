package z80

// Provides T states counters
type TCounter struct {
	Total   int64 // total T states since boot or hard reset
	Current int   // T states for the current frame
	max     int   // T states to run
}

// Add (or subtract) specified number of T states
func (tc *TCounter) Add(t int) {
	tc.Total += int64(t)
	tc.Current += t
}

// Set the limit of T states to execute
func (tc *TCounter) limit(max int) {
	tc.max = max + tc.remaining()
	tc.Current = 0
}

// Checks whether the maximum of T states has been met or exceeded
func (tc *TCounter) done() bool {
	return tc.max != 0 && tc.Current >= tc.max
}

// Add T states for the halt operation
func (tc *TCounter) halt() {
	r := tc.remaining()
	if r%4 != 0 {
		r += 4 - r%4
	}
	tc.Add(r)
}

// Get the remaining T states
func (tc *TCounter) remaining() int {
	return tc.max - tc.Current
}

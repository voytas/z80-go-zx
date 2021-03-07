package z80

// Represents the state that might be needed by external components.
type State struct {
	// Keeps count of the executed t-states. This value can be updated
	// if some external factors impact the execution, for example
	// contended memory that adds extra t-states.
	TCount int
}

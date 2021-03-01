package exerciser

import (
	"testing"
)

func Benchmark_Run(t *testing.B) {
	for i := 0; i < 100; i++ {
		Run("/Users/voytas/Projects/go/src/z80-go-zx/exercises/prelim.com")
	}
}

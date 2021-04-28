package snapshot

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/voytas/z80-go-zx/spectrum/memory"
	"github.com/voytas/z80-go-zx/z80"
)

func Load(filePath string, cpu *z80.Z80, mem *memory.Memory) error {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".sna":
		sna := &SNA{}
		return sna.Load(filePath, cpu, mem)
	case ".szx":
		szx := &SZX{}
		return szx.Load(filePath, cpu, mem)
	default:
		return fmt.Errorf("File format not supported: %s", ext)
	}
}

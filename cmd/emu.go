package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/voytas/z80-go-zx/spectrum"
	"github.com/voytas/z80-go-zx/spectrum/machine"
)

var Model string

var emuCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "emu -m 48k|128k [snapshot.(sna|szx)]",
	Short: "Run ZX Spectrum emulator",
	Long: `
		Run ZX Spectrum emulator. You can optionally specify snapshot file to load,
		only SNA & SZX formats are supported.

		Supported models are 48k and 128k`,
	Run: func(cmd *cobra.Command, args []string) {
		var fileName string
		if len(args) > 0 {
			fileName = args[0]
		}
		var m *machine.Machine
		switch strings.TrimSpace(Model) {
		case "48k":
			m = machine.ZX48k
		case "128k":
			m = machine.ZX128k
		}
		spectrum.Run(m, fileName)
	},
}

func init() {
	emuCmd.Flags().StringVarP(&Model, "model", "m", "48k", "Model to run: 48k or 128k")
	rootCmd.AddCommand(emuCmd)
}

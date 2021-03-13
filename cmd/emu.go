package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voytas/z80-go-zx/spectrum"
	"github.com/voytas/z80-go-zx/spectrum/settings"
)

var emuCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "emu [file.sna]",
	Short: "Run ZX Spectrum emulator",
	Long: `
		Run ZX Spectrum emulator. You can optionally specify program to run,
		but only SNA format is currently supported`,
	Run: func(cmd *cobra.Command, args []string) {
		var fileName string
		if len(args) > 0 {
			fileName = args[0]
		}
		spectrum.Run(settings.ZX48k, fileName)
	},
}

func init() {
	rootCmd.AddCommand(emuCmd)
}

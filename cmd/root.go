package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "z80 command",
	Short: "Simple command line tool to execute tests.",
	Long: `
	This utility has been created to execute:

	Preliminary Z80 tests                    - prelim.com
	Documented Z80 instruction set exerciser - zexdoc.com
	All Z80 instruction set exerciser        - zexall.com

	or ZX Spectrum 48k/128k emulator.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

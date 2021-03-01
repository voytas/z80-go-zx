package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voytas/z80-go-zx/exerciser"
)

var exerciseCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "exercise program",
	Short: "Executes specified test program",
	Long: `The program can be for example prelim.com,
		zexdoc.com or zexall.com`,
	Run: func(cmd *cobra.Command, args []string) {
		exerciser.Run(args[0])
	},
}

func init() {
	rootCmd.AddCommand(exerciseCmd)
}

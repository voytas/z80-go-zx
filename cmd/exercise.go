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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exerciseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exerciseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

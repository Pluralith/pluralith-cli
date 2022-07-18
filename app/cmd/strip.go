package cmd

import (
	"pluralith/pkg/strip"

	"github.com/spf13/cobra"
)

// stripCmd represents the strip command
var stripCmd = &cobra.Command{
	Use:   "strip",
	Short: "Strip a given state file of secrets according to config",
	Long:  `Strip a given state file of secrets according to config`,
	Run: func(cmd *cobra.Command, args []string) {
		strip.StripAndHash()
	},
}

func init() {
	rootCmd.AddCommand(stripCmd)
}

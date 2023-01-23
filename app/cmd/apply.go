package cmd

import (
	"fmt"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

// applyOldCmd represents the applyOld command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run terraform apply and show changes in Pluralith",
	Long:  `Run terraform apply and show changes in Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintFormatted("✘", []string{"red", "bold"})
		fmt.Print(" The pluralith apply command is deprecated. Please use the pluralith graph command\n\n")
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

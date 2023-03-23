package cmd

import (
	"fmt"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

// planOldCmd represents the planOld command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run terraform plan and show changes in Pluralith",
	Long:  `Run terraform plan and show changes in Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintFormatted("✘", []string{"red", "bold"})
		fmt.Print(" The Pluralith plan command is deprecated. Please use the pluralith graph command\n\n")
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}

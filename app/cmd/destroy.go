package cmd

import (
	"fmt"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

// destroyOldCmd represents the destroyOld command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Run terraform destroy and show changes in Pluralith",
	Long:  `Run terraform destroy and show changes in Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintFormatted("✘", []string{"red", "bold"})
		fmt.Print(" The Pluralith destroy command is deprecated. Please use the pluralith graph command\n\n")
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}

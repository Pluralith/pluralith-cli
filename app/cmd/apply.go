package cmd

import (
	"fmt"
	"pluralith/pkg/terraform"

	"github.com/spf13/cobra"
)

// applyOldCmd represents the applyOld command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run terraform apply and show changes in Pluralith",
	Long:  `Run terraform apply and show changes in Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := terraform.RunTerraform("apply", cmd.Flags()); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

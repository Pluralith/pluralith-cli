package cmd

import (
	"fmt"
	"pluralith/pkg/terraform"

	"github.com/spf13/cobra"
)

// planOldCmd represents the planOld command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run terraform plan and show changes in Pluralith",
	Long:  `Run terraform plan and show changes in Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := terraform.RunTerraform("plan", args); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}

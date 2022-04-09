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
		varArgs, varErr := terraform.ConstructVarArgs(cmd.Flags())
		if varErr != nil {
			fmt.Println(varErr)
		}

		if planErr := terraform.RunTerraform("plan", varArgs); planErr != nil {
			fmt.Println(planErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.PersistentFlags().StringSlice("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	planCmd.PersistentFlags().StringSlice("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
}

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
		varArgs, varErr := terraform.ConstructVarArgs(cmd.Flags())
		if varErr != nil {
			fmt.Println(varErr)
		}

		if applyErr := terraform.RunTerraform("apply", varArgs); applyErr != nil {
			fmt.Println(applyErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.PersistentFlags().StringSlice("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	applyCmd.PersistentFlags().StringSlice("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
}

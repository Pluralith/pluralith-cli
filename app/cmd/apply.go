package cmd

import (
	"fmt"
	"pluralith/pkg/cost"
	"pluralith/pkg/terraform"

	"github.com/spf13/cobra"
)

// applyOldCmd represents the applyOld command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run terraform apply and show changes in Pluralith",
	Long:  `Run terraform apply and show changes in Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		tfArgs := terraform.ConstructTerraformArgs(cmd.Flags())
		costArgs := cost.ConstructInfracostArgs(cmd.Flags())

		if applyErr := terraform.RunTerraform("apply", tfArgs, costArgs); applyErr != nil {
			fmt.Println(applyErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.PersistentFlags().StringSlice("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	applyCmd.PersistentFlags().StringSlice("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
	applyCmd.PersistentFlags().String("cost-usage-file", "", "Path to an infracost usage file to be used for the cost breakdown")
	applyCmd.PersistentFlags().Bool("show-costs", false, "If we detect infracost we automatically run a cost breakdown and show it in the diagram. Use this flag to turn that off")
}

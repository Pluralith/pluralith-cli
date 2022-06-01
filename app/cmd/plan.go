package cmd

import (
	"fmt"
	"pluralith/pkg/cost"
	"pluralith/pkg/terraform"

	"github.com/spf13/cobra"
)

// planOldCmd represents the planOld command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run terraform plan and show changes in Pluralith",
	Long:  `Run terraform plan and show changes in Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		tfArgs := terraform.ConstructTerraformArgs(cmd.Flags())
		costArgs := cost.ConstructInfracostArgs(cmd.Flags())

		if planErr := terraform.RunTerraform("plan", tfArgs, costArgs); planErr != nil {
			fmt.Println(planErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.PersistentFlags().StringSlice("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	planCmd.PersistentFlags().StringSlice("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
	planCmd.PersistentFlags().String("cost-usage-file", "", "Path to an infracost usage file to be used for the cost breakdown")
	planCmd.PersistentFlags().Bool("show-costs", false, "If we detect infracost we automatically run a cost breakdown and show it in the diagram. Use this flag to turn that off")
}

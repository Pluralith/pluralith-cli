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
		costArgs, costErr := cost.ConstructInfracostArgs(cmd.Flags())
		if costErr != nil {
			fmt.Println(costErr)
			return
		}

		if applyErr := terraform.RunTerraform("apply", tfArgs, costArgs); applyErr != nil {
			fmt.Println(applyErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.PersistentFlags().String("plan-file", "", "Path to an execution plan binary file. If passed, this will skip a plan run under the hood.")
	applyCmd.PersistentFlags().StringArray("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	applyCmd.PersistentFlags().StringArray("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
	applyCmd.PersistentFlags().String("cost-usage-file", "", "Path to an infracost usage file to be used for the cost breakdown")
	applyCmd.PersistentFlags().Bool("show-costs", false, "If we detect infracost we automatically run a cost breakdown and show it in the diagram. Use this flag to turn that off")
}

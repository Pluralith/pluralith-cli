package cmd

import (
	"fmt"
	"pluralith/pkg/cost"
	"pluralith/pkg/terraform"

	"github.com/spf13/cobra"
)

// destroyOldCmd represents the destroyOld command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Run terraform destroy and show changes in Pluralith",
	Long:  `Run terraform destroy and show changes in Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		tfArgs := terraform.ConstructTerraformArgs(cmd.Flags())
		costArgs, costErr := cost.ConstructInfracostArgs(cmd.Flags())
		if costErr != nil {
			fmt.Println(costErr)
			return
		}

		if destroyErr := terraform.RunTerraform("destroy", tfArgs, costArgs); destroyErr != nil {
			fmt.Println(destroyErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
	destroyCmd.PersistentFlags().String("plan-file", "", "Path to an execution plan binary file. If passed, this will skip a plan run under the hood.")
	destroyCmd.PersistentFlags().StringSlice("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	destroyCmd.PersistentFlags().StringSlice("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
	destroyCmd.PersistentFlags().String("cost-usage-file", "", "Path to an infracost usage file to be used for the cost breakdown")
	destroyCmd.PersistentFlags().Bool("show-costs", false, "If we detect infracost we automatically run a cost breakdown and show it in the diagram. Use this flag to turn that off")
}

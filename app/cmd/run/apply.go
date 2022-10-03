package run

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RunApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run terraform apply and push updates to Pluralith",
	Long:  `Run terraform apply and push updates to Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ci apply")
		// tfArgs := terraform.ConstructTerraformArgs(cmd.Flags())
		// costArgs, costErr := cost.ConstructInfracostArgs(cmd.Flags())
		// if costErr != nil {
		// 	fmt.Println(costErr)
		// 	return
		// }

		// if applyErr := terraform.RunTerraform("apply", tfArgs, costArgs); applyErr != nil {
		// 	fmt.Println(applyErr)
		// }
	},
}

func init() {
	// runc.AddCommand(applyCmd)
	// runApplyCmd.PersistentFlags().String("plan-file", "", "Path to an execution plan binary file. If passed, this will skip a plan run under the hood.")
	// runApplyCmd.PersistentFlags().StringArray("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	// runApplyCmd.PersistentFlags().StringArray("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
	// runApplyCmd.PersistentFlags().String("cost-usage-file", "", "Path to an infracost usage file to be used for the cost breakdown")
	// runApplyCmd.PersistentFlags().Bool("show-costs", true, "If we detect infracost we automatically run a cost breakdown and show it in the diagram. Use this flag to turn that off")
}

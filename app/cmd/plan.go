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
		// showCosts, flagErr := cmd.Flags().GetBool("show-costs")
		// if flagErr != nil {
		// 	fmt.Println(flagErr)
		// }

		// if showCosts {
		// 	cost.CalculateCost()
		// }

		// fmt.Println(showCosts)

		if err := terraform.RunTerraform("plan", cmd.Flags()); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.PersistentFlags().Bool("show-costs", false, "Shows cost information in Pluralith desktop (courtesy of Infracost)")
}

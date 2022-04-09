package cmd

import (
	"fmt"
	"pluralith/pkg/terraform"

	"github.com/spf13/cobra"
)

// destroyOldCmd represents the destroyOld command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Run terraform destroy and show changes in Pluralith",
	Long:  `Run terraform destroy and show changes in Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		varArgs, varErr := terraform.ConstructVarArgs(cmd.Flags())
		if varErr != nil {
			fmt.Println(varErr)
		}

		if destroyErr := terraform.RunTerraform("destroy", varArgs); destroyErr != nil {
			fmt.Println(destroyErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
	destroyCmd.PersistentFlags().StringSlice("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	destroyCmd.PersistentFlags().StringSlice("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
}

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
		if err := terraform.RunTerraform("destroy", cmd.Flags()); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}

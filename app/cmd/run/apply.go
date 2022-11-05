package run

import (
	"fmt"
	"pluralith/pkg/ci"
	"pluralith/pkg/graph"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

var RunApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run terraform apply and push updates to Pluralith",
	Long:  `Run terraform apply and push updates to Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintFormatted("⠿", []string{"blue", "bold"})
		fmt.Println(" Initiating Apply Run ⇢ Posting To Pluralith Dashboard")

		// - - Prepare for Run - -
		tfArgs, costArgs, exportArgs, preErr := ci.PreRun(cmd.Flags())
		if preErr != nil {
			fmt.Println(preErr)
			return
		}

		// - - Generate Graph - -
		if graphErr := graph.GenerateGraph("apply", tfArgs, costArgs, exportArgs, true); graphErr != nil {
			fmt.Println(graphErr)
			return
		}

		// - - Post Graph - -
		if ciError := ci.PostGraph("apply", exportArgs); ciError != nil {
			fmt.Println(ciError)
			return
		}

		// - - Run Apply - -
		if eventErr := ci.PostEvents("apply", tfArgs, costArgs, exportArgs); eventErr != nil {
			fmt.Println(eventErr)
			return
		}
	},
}

func init() {}

package run

import (
	"fmt"
	"pluralith/pkg/backends"
	"pluralith/pkg/ci"
	"pluralith/pkg/graph"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

var RunPlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run terraform plan and push updates to Pluralith",
	Long:  `Run terraform plan and push updates to Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintFormatted("⠿", []string{"blue", "bold"})
		fmt.Println(" Initiating Plan Run ⇢ Posting To Pluralith Dashboard")

		// - - Prepare for Run - -
		tfArgs, costArgs, exportArgs, preErr := ci.PreRun(cmd.Flags())
		if preErr != nil {
			fmt.Println(preErr)
			return
		}

		// - - Generate Graph - -
		if graphErr := graph.GenerateGraph("plan", tfArgs, costArgs, exportArgs, true); graphErr != nil {
			fmt.Println(graphErr)
			return
		}

		// - - Post Graph - -
		if ciError := ci.PostGraph("plan", exportArgs); ciError != nil {
			fmt.Println(ciError)
			return
		}

		// - - Push Diagram to State Backend - -
		if exportArgs["sync-to-backend"] == true {
			if pushErr := backends.StoreInBackend(); pushErr != nil {
				fmt.Println(pushErr)
				return
			}
		}
	},
}

func init() {}

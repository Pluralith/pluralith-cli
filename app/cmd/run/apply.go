package run

import (
	"fmt"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/backends"
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
		preValid, tfArgs, costArgs, exportArgs, preErr := ci.PreRun(cmd.Flags())
		if preErr != nil {
			fmt.Println(preErr)
			return
		}
		if !preValid {
			return
		}

		// - - Generate Graph - -
		if graphErr := graph.GenerateGraph("apply", tfArgs, costArgs, exportArgs, true, false); graphErr != nil {
			fmt.Println(graphErr)
			return
		}

		// - - Post Graph - -
		if ciError := ci.PostGraph("apply", exportArgs); ciError != nil {
			fmt.Println(ciError)
			return
		}

		// - - Push Diagram to State Backend - -
		if exportArgs["sync-to-backend"] == true || auxiliary.StateInstance.PluralithConfig.Config.SyncToBackend {
			if pushErr := backends.SyncToBackend(); pushErr != nil {
				fmt.Println(pushErr)
				return
			}
		}

		// - - Run Apply - -
		if eventErr := ci.PostEvents("apply", tfArgs, costArgs, exportArgs); eventErr != nil {
			fmt.Println(eventErr)
			return
		}
	},
}

func init() {}

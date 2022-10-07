package run

import (
	"fmt"
	"pluralith/pkg/ci"
	"pluralith/pkg/graph"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

var RunDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Run terraform destroy and push updates to Pluralith",
	Long:  `Run terraform destroy and push updates to Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintFormatted("⠿", []string{"blue", "bold"})
		fmt.Println(" Initiating Destroy Run ⇢ Posting To Pluralith Dashboard")

		// - - Prepare for Run - -
		tfArgs, costArgs, exportArgs, preErr := ci.PreRun(cmd.Flags())
		if preErr != nil {
			fmt.Println(preErr)
		}

		// - - Generate Graph - -
		if graphErr := graph.GenerateGraph(tfArgs, costArgs, exportArgs, true); graphErr != nil {
			fmt.Println(graphErr)
		}

		// - - Post Graph - -
		if ciError := ci.PostGraph("destroy", exportArgs); ciError != nil {
			fmt.Println(ciError)
		}

		// - - Run Destroy - -
		if eventErr := ci.PostEvents("destroy", tfArgs, costArgs, exportArgs); eventErr != nil {
			fmt.Println(eventErr)
		}
	},
}

func init() {}

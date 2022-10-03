package run

import (
	"fmt"
	"pluralith/pkg/graph"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

var RunPlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run terraform plan and push updates to Pluralith",
	Long:  `Run terraform plan and push updates to Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		// Print UX head
		ux.PrintFormatted("⠿", []string{"blue", "bold"})
		fmt.Println(" Initiating Plan Run ⇢ Posting To Pluralith Dashboard")

		tfArgs, costArgs, exportArgs, preErr := PreRun(cmd.Flags())
		if preErr != nil {
			fmt.Println(preErr)
		}

		if graphErr := graph.RunGraph(tfArgs, costArgs, exportArgs, true); graphErr != nil {
			fmt.Println(graphErr)
		}

		if ciError := graph.HandleCIRun(exportArgs); ciError != nil {
			fmt.Println(ciError)
		}
	},
}

func init() {}

package cmd

import (
	"fmt"
	"pluralith/pkg/graph"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

// stripCmd represents the strip command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Run terraform plan locally and push diagram to Pluralith",
	Long:  `Run terraform plan locally and push diagram to Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintFormatted("⠿", []string{"blue", "bold"})
		fmt.Println(" Initiating Graph locally ⇢ Posting Diagram To Pluralith Dashboard")

		// - - Prepare for Run - -
		preValid, tfArgs, costArgs, exportArgs, preErr := graph.PreGraph(cmd.Flags())
		if preErr != nil {
			fmt.Println("Error: ", preErr)
			return
		}
		if !preValid {
			return
		}

		// - - Generate Graph - -
		if graphErr := graph.GenerateGraph("plan", tfArgs, costArgs, exportArgs, true, true); graphErr != nil {
			fmt.Println(graphErr)
			return
		}

		// Post diagram if local-only flag not set
		if !exportArgs["export-pdf"].(bool) {
			// - - Post Graph - -
			if ciError := graph.PostGraph("plan", exportArgs); ciError != nil {
				fmt.Println(ciError)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
	graphCmd.PersistentFlags().Bool("local-only", false, "Diagram will not be pushed to the Pluralith Dashboard. Instead, a PDF export of the Diagram gets generated locally")
	graphCmd.PersistentFlags().String("title", "Pluralith Diagram", "The title for your diagram, will be displayed in the PDF output")
	graphCmd.PersistentFlags().String("author", "", "The author/creator of the diagram, will be displayed in the PDF output")
	graphCmd.PersistentFlags().String("version", "", "The diagram version, will be displayed in the PDF output")
	graphCmd.PersistentFlags().String("out-dir", "", "The directory the diagram should be exported to")
	graphCmd.PersistentFlags().String("file-name", "", "The name of the exported PDF")
	graphCmd.PersistentFlags().Bool("show-changes", false, "Determines whether the exported diagram highlights changes made in the latest Terraform plan or outputs a general diagram of the infrastructure")
	graphCmd.PersistentFlags().Bool("show-drift", false, "Determines whether the exported diagram highlights resource drift detected by Terraform")
	graphCmd.PersistentFlags().Bool("show-costs", false, "Determines whether the exported diagram includes cost information")
	graphCmd.PersistentFlags().String("cost-mode", "delta", "Determines which costs are shown. Can be 'delta' or 'total'")
	graphCmd.PersistentFlags().String("cost-period", "month", "Determines over which period costs are aggregated. Can be 'hour' or 'month'")
	graphCmd.PersistentFlags().String("cost-usage-file", "", "Path to an infracost usage file to be used for the cost breakdown")
	graphCmd.PersistentFlags().String("plan-file", "", "Path to an execution plan binary file. If passed, this will skip a plan run under the hood.")
	graphCmd.PersistentFlags().String("plan-file-json", "", "Path to an execution plan json file. If passed, this will skip a plan run under the hood.")
	graphCmd.PersistentFlags().StringArray("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	graphCmd.PersistentFlags().StringArray("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
}

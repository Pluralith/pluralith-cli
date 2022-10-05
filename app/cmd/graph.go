package cmd

import (
	"fmt"
	"pluralith/pkg/auth"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/cost"
	"pluralith/pkg/graph"
	"pluralith/pkg/terraform"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

// stripCmd represents the strip command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Generate and export a diagram of the current plan state as a PDF",
	Long:  `Generate and export a diagram of the current plan state as a PDF`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check for CI environment -> If CI -> Prompt user to use `pluralith run`
		if auxiliary.StateInstance.IsCI {
			ux.PrintFormatted("✘", []string{"red", "bold"})
			fmt.Print(" CI environment detected ⇢ Use ")
			ux.PrintFormatted("'pluralith run'", []string{"red", "bold"})
			fmt.Println(" to generate a diagram and post your run")
			return
		}

		ux.PrintFormatted("⠿", []string{"blue", "bold"})
		fmt.Print(" Exporting Diagram\n\n")

		tfArgs := terraform.ConstructTerraformArgs(cmd.Flags())
		costArgs, costErr := cost.ConstructInfracostArgs(cmd.Flags())
		if costErr != nil {
			fmt.Println(costErr)
			return
		}

		exportArgs := graph.ConstructExportArgs(cmd.Flags())
		exportArgs["export-pdf"] = true // Always export pdf when running locally

		configValid, _, configErr := auth.VerifyConfig(true)
		if !configValid {
			return
		}
		if configErr != nil {
			fmt.Println(configErr)
		}

		if graphErr := graph.GenerateGraph(tfArgs, costArgs, exportArgs, false); graphErr != nil {
			fmt.Println(graphErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
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
	graphCmd.PersistentFlags().StringArray("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	graphCmd.PersistentFlags().StringArray("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
}

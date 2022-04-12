package cmd

import (
	"fmt"
	"pluralith/pkg/cost"
	"pluralith/pkg/graph"
	"pluralith/pkg/terraform"

	"github.com/spf13/cobra"
)

// stripCmd represents the strip command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Generate and export a diagram of the current plan state as a PDF",
	Long:  `Generate and export a diagram of the current plan state as a PDF`,
	Run: func(cmd *cobra.Command, args []string) {
		tfArgs, tfErr := terraform.ConstructTerraformArgs(cmd.Flags())
		if tfErr != nil {
			fmt.Println(tfErr)
		}

		costArgs, costErr := cost.ConstructInfracostArgs(cmd.Flags())
		if costErr != nil {
			fmt.Println(costErr)
		}

		exportArgs, exportErr := graph.ConstructExportArgs(cmd.Flags())
		if exportErr != nil {
			fmt.Println(fmt.Errorf("getting diagram values failed -> %w", exportErr))
			return
		}

		if graphErr := graph.RunGraph(tfArgs, costArgs, exportArgs); graphErr != nil {
			fmt.Println(graphErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
	graphCmd.PersistentFlags().String("title", "", "The title for your diagram, will be displayed in the PDF output")
	graphCmd.PersistentFlags().String("author", "", "The author/creator of the diagram, will be displayed in the PDF output")
	graphCmd.PersistentFlags().String("version", "", "The diagram version, will be displayed in the PDF output")
	graphCmd.PersistentFlags().String("out-dir", "", "The directory the diagram should be exported to")
	graphCmd.PersistentFlags().String("file-name", "", "The name of the exported PDF")
	graphCmd.PersistentFlags().Bool("show-changes", false, "Determines whether the exported PDF highlights changes made in the latest Terraform plan or outputs a general diagram of the infrastructure")
	graphCmd.PersistentFlags().Bool("show-drift", false, "Determines whether the exported PDF highlights resource drift detected by Terraform")
	graphCmd.PersistentFlags().Bool("skip-plan", false, "Generates a diagram without running plan again (needs pluralith state from previous plan run)")
	graphCmd.PersistentFlags().Bool("generate-md", false, "Generate markdown output with exported PDF link and preview image for pull request comments")
	graphCmd.PersistentFlags().StringSlice("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	graphCmd.PersistentFlags().StringSlice("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
	graphCmd.PersistentFlags().String("cost-usage-file", "", "Path to an infracost usage file to be used for the cost breakdown")
	graphCmd.PersistentFlags().Bool("no-costs", false, "If we detect infracost we automatically run a cost breakdown and show it in the diagram. Use this flag to turn that off")
}

package cmd

import (
	"fmt"
	"pluralith/pkg/cost"
	"pluralith/pkg/graph"
	"pluralith/pkg/terraform"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

// stripCmd represents the strip command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Generate a diagram in CI and post it to the Pluralith dashboard",
	Long:  `Generate a diagram in CI and post it to the Pluralith dashboard`,
	Run: func(cmd *cobra.Command, args []string) {
		var runAsCI bool = true // Run as if in CI, no matter what 'IsCI' in state is

		// Print UX head
		ux.PrintFormatted("⠿", []string{"blue", "bold"})
		fmt.Print(" Initiating Run ⇢ Posting To Pluralith Dashboard\n\n")

		tfArgs, tfErr := terraform.ConstructTerraformArgs(cmd.Flags())
		if tfErr != nil {
			fmt.Println(tfErr)
		}

		costArgs, costErr := cost.ConstructInfracostArgs(cmd.Flags())
		if costErr != nil {
			fmt.Println(costErr)
		}

		configValid, configErr := graph.VerifyConfig(false)
		if !configValid {
			return
		}
		if configErr != nil {
			fmt.Println(configErr)
			return
		}

		exportArgs, exportErr := graph.ConstructExportArgs(cmd.Flags(), runAsCI)
		if exportErr != nil {
			if exportErr.Error() == "Handled" {
				return
			}
			fmt.Println(fmt.Errorf("getting diagram values failed -> %w", exportErr))
			return
		}

		if graphErr := graph.RunGraph(tfArgs, costArgs, exportArgs, runAsCI); graphErr != nil {
			fmt.Println(graphErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().String("title", "", "The title for your diagram, will be displayed in the PDF output")
	runCmd.PersistentFlags().String("author", "", "The author/creator of the diagram, will be displayed in the PDF output")
	runCmd.PersistentFlags().String("version", "", "The diagram version, will be displayed in the PDF output")
	runCmd.PersistentFlags().String("out-dir", "", "The directory the diagram should be exported to")
	runCmd.PersistentFlags().String("file-name", "", "The name of the exported PDF")
	runCmd.PersistentFlags().Bool("show-changes", false, "Determines whether the exported PDF highlights changes made in the latest Terraform plan or outputs a general diagram of the infrastructure")
	runCmd.PersistentFlags().Bool("show-drift", false, "Determines whether the exported PDF highlights resource drift detected by Terraform")
	runCmd.PersistentFlags().StringSlice("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	runCmd.PersistentFlags().StringSlice("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
	runCmd.PersistentFlags().String("cost-usage-file", "", "Path to an infracost usage file to be used for the cost breakdown")
	runCmd.PersistentFlags().Bool("no-costs", false, "If we detect infracost we automatically run a cost breakdown and show it in the diagram. Use this flag to turn that off")
}

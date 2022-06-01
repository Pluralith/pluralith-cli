package cmd

import (
	"fmt"
	"math/rand"
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
		// Print UX head
		ux.PrintFormatted("⠿", []string{"blue", "bold"})
		fmt.Print(" Initiating Run ⇢ Posting To Pluralith Dashboard\n\n")

		tfArgs := terraform.ConstructTerraformArgs(cmd.Flags())
		costArgs := cost.ConstructInfracostArgs(cmd.Flags())

		configValid, configErr := graph.VerifyConfig(false)
		if !configValid {
			return
		}
		if configErr != nil {
			fmt.Println(configErr)
			return
		}

		exportArgs := graph.ConstructExportArgs(cmd.Flags())
		exportArgs["title"] = fmt.Sprintf("%07d", rand.Intn(10000000)) // Generate random run id
		exportArgs["show-costs"] = true                                // Always run infracost in CI if infracost is installed

		if graphErr := graph.RunGraph(tfArgs, costArgs, exportArgs, true); graphErr != nil {
			fmt.Println(graphErr)
		}

		if ciError := graph.HandleCIRun(exportArgs); ciError != nil {
			fmt.Println(ciError)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().String("version", "", "The diagram version, will be displayed in the PDF output")
	runCmd.PersistentFlags().String("out-dir", "", "The directory the diagram should be exported to")
	runCmd.PersistentFlags().String("file-name", "", "The name of the exported PDF")
	runCmd.PersistentFlags().Bool("export-pdf", false, "Determines whether a PDF export of the run Diagram is generated locally")
	runCmd.PersistentFlags().Bool("show-changes", false, "Determines whether the exported diagram highlights changes made in the latest Terraform plan or outputs a general diagram of the infrastructure")
	runCmd.PersistentFlags().Bool("show-drift", false, "Determines whether the exported diagram highlights resource drift detected by Terraform")
	runCmd.PersistentFlags().Bool("show-costs", false, "Determines whether the exported diagram includes cost information")
	runCmd.PersistentFlags().StringSlice("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	runCmd.PersistentFlags().StringSlice("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
	runCmd.PersistentFlags().String("cost-usage-file", "", "Path to an infracost usage file to be used for the cost breakdown")
	runCmd.PersistentFlags().Bool("no-costs", false, "If we detect infracost we automatically run a cost breakdown and show it in the diagram. Use this flag to turn that off")
}

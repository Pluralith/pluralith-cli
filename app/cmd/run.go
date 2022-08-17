package cmd

import (
	"fmt"
	"math/rand"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/cost"
	"pluralith/pkg/graph"
	"pluralith/pkg/terraform"
	"pluralith/pkg/ux"
	"strconv"

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
		fmt.Println(" Initiating Run ⇢ Posting To Pluralith Dashboard")

		if auxiliary.StateInstance.Branch != "none" {
			ux.PrintFormatted(" →", []string{"blue", "bold"})
			fmt.Print(" Branch detected: ")
			ux.PrintFormatted(auxiliary.StateInstance.Branch+"\n\n", []string{"blue"})
		} else {
			ux.PrintFormatted(" →", []string{"blue", "bold"})
			fmt.Print(" Could not detect branch\n\n")
		}

		tfArgs := terraform.ConstructTerraformArgs(cmd.Flags())
		costArgs, costErr := cost.ConstructInfracostArgs(cmd.Flags())
		if costErr != nil {
			fmt.Println(costErr)
			return
		}

		exportArgs := graph.ConstructExportArgs(cmd.Flags())
		exportArgs["runId"] = fmt.Sprintf("%07d", rand.Intn(10000000)) // Generate random run id
		costArgs["show-costs"] = true                                  // Always run infracost in CI if infracost is installed

		// Set defaults for export
		if exportArgs["title"] == "" {
			exportArgs["title"] = "Run #" + exportArgs["runId"].(string)
			exportArgs["file-name"] = "Run_" + exportArgs["runId"].(string)
		}

		configValid, projectData, configErr := graph.VerifyConfig(false)
		if !configValid {
			return
		}
		if configErr != nil {
			fmt.Println(configErr)
			return
		}

		projectData = projectData["data"].(map[string]interface{})
		exportArgs["orgId"] = strconv.Itoa(int(projectData["orgId"].(float64)))
		exportArgs["projectId"] = auxiliary.StateInstance.PluralithConfig.ProjectId

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
	runCmd.PersistentFlags().String("title", "", "The title for your diagram, will be displayed in the PDF output")
	runCmd.PersistentFlags().String("version", "", "The diagram version, will be displayed in the PDF output")
	runCmd.PersistentFlags().String("out-dir", "", "The directory the diagram should be exported to")
	runCmd.PersistentFlags().String("file-name", "", "The name of the exported PDF")
	runCmd.PersistentFlags().Bool("export-pdf", false, "Determines whether a PDF export of the run Diagram is generated locally")
	runCmd.PersistentFlags().Bool("show-changes", false, "Determines whether the exported diagram highlights changes made in the latest Terraform plan or outputs a general diagram of the infrastructure")
	runCmd.PersistentFlags().Bool("show-drift", false, "Determines whether the exported diagram highlights resource drift detected by Terraform")
	runCmd.PersistentFlags().Bool("show-costs", false, "Determines whether the exported diagram includes cost information")
	runCmd.PersistentFlags().String("cost-mode", "delta", "Determines which costs are shown. Can be 'delta' or 'total'")
	runCmd.PersistentFlags().String("cost-period", "month", "Determines over which period costs are aggregated. Can be 'hour' or 'month'")
	runCmd.PersistentFlags().StringArray("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	runCmd.PersistentFlags().StringArray("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
	runCmd.PersistentFlags().String("cost-usage-file", "", "Path to an infracost usage file to be used for the cost breakdown")
}

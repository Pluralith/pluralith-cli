package run

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

var RunPlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run terraform plan and push updates to Pluralith",
	Long:  `Run terraform plan and push updates to Pluralith`,
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
		// projectData["orgId"] = ""
		// fmt.Println(int(projectData["orgId"].(float64)))
		exportArgs["projectId"] = auxiliary.StateInstance.PluralithConfig.ProjectId

		if graphErr := graph.RunGraph(tfArgs, costArgs, exportArgs, true); graphErr != nil {
			fmt.Println(graphErr)
		}

		if ciError := graph.HandleCIRun(exportArgs); ciError != nil {
			fmt.Println(ciError)
		}
	},
}

func init() {}

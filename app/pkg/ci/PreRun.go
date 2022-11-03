package ci

import (
	"fmt"
	"math/rand"
	"pluralith/pkg/auth"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/cost"
	"pluralith/pkg/graph"
	"pluralith/pkg/terraform"
	"pluralith/pkg/ux"
	"reflect"
	"strconv"

	"github.com/spf13/pflag"
)

func PreRun(flags *pflag.FlagSet) (map[string]interface{}, map[string]interface{}, map[string]interface{}, error) {
	functionName := "PreRun"

	if auxiliary.StateInstance.Branch != "none" {
		ux.PrintFormatted(" →", []string{"blue", "bold"})
		fmt.Print(" Branch detected: ")
		ux.PrintFormatted(auxiliary.StateInstance.Branch+"\n\n", []string{"blue"})
	} else {
		ux.PrintFormatted(" →", []string{"blue", "bold"})
		fmt.Print(" Could not detect branch\n\n")
	}

	tfArgs := terraform.ConstructTerraformArgs(flags)
	costArgs, costErr := cost.ConstructInfracostArgs(flags)
	if costErr != nil {
		return nil, nil, nil, fmt.Errorf("%v: %w", functionName, costErr)
	}

	exportArgs := graph.ConstructExportArgs(flags)
	exportArgs["runId"] = fmt.Sprintf("%07d", rand.Intn(10000000)) // Generate random run id
	costArgs["show-costs"] = true                                  // Always run infracost in CI if infracost is installed

	// Set defaults for export
	if exportArgs["title"] == "" {
		exportArgs["title"] = "Infrastructure Diagram"     //+ exportArgs["runId"].(string)
		exportArgs["file-name"] = "Infrastructure_Diagram" //+ exportArgs["runId"].(string)
	}

	configValid, projectData, configErr := auth.VerifyConfig(false)
	if !configValid {
		return nil, nil, nil, fmt.Errorf("invalid Pluralith config")
	}
	if configErr != nil {
		return nil, nil, nil, fmt.Errorf("%v: %w", functionName, configErr)
	}

	projectData = projectData["data"].(map[string]interface{})
	// Failsafe for when orgId is a string
	if reflect.TypeOf(projectData["orgId"]).Kind() == reflect.String {
		exportArgs["orgId"] = projectData["orgId"]
	} else {
		exportArgs["orgId"] = strconv.Itoa(int(projectData["orgId"].(float64)))
	}

	exportArgs["projectId"] = auxiliary.StateInstance.PluralithConfig.ProjectId

	return tfArgs, costArgs, exportArgs, nil
}

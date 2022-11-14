package ci

import (
	"fmt"
	"math/rand"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/cost"
	"pluralith/pkg/graph"
	"pluralith/pkg/initialization"
	"pluralith/pkg/terraform"
	"pluralith/pkg/ux"

	"github.com/spf13/pflag"
)

func PreRun(flags *pflag.FlagSet) (map[string]interface{}, map[string]interface{}, map[string]interface{}, error) {
	functionName := "PreRun"

	if auxiliary.StateInstance.Branch != "none" {
		ux.PrintFormatted(" →", []string{"blue", "bold"})
		fmt.Print(" Branch detected: ")
		ux.PrintFormatted(auxiliary.StateInstance.Branch+"\n", []string{"blue"})
	} else {
		ux.PrintFormatted(" →", []string{"blue", "bold"})
		fmt.Print(" Could not detect branch\n")
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

	initData, initErr := initialization.RunInit(false, initialization.InitData{})
	if initErr != nil {
		return nil, nil, nil, fmt.Errorf("%v: %w", functionName, costErr)
	}

	exportArgs["orgId"] = initData.OrgId
	exportArgs["projectId"] = initData.ProjectId
	exportArgs["projectName"] = initData.ProjectName

	return tfArgs, costArgs, exportArgs, nil
}

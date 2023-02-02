package ci

import (
	"fmt"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/cost"
	"pluralith/pkg/graph"
	"pluralith/pkg/initialization"
	"pluralith/pkg/terraform"
	"pluralith/pkg/ux"

	"github.com/spf13/pflag"
)

func PreRun(flags *pflag.FlagSet) (bool, map[string]interface{}, map[string]interface{}, map[string]interface{}, error) {
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
		return false, nil, nil, nil, fmt.Errorf("%v: %w", functionName, costErr)
	}

	exportArgs := graph.ConstructExportArgs(flags)
	costArgs["show-costs"] = true // Always run infracost in CI if infracost is installed

	// Set defaults for export
	if exportArgs["title"] == "" {
		exportArgs["title"] = "Infrastructure Diagram"
		exportArgs["file-name"] = "Infrastructure_Diagram"
	}

	initValid, initData, initErr := initialization.RunInit(true, initialization.InitData{}, false)
	if initErr != nil {
		return false, nil, nil, nil, fmt.Errorf("%v: %w", functionName, initErr)
	}
	if !initValid {
		return false, nil, nil, nil, nil
	}

	exportArgs["org-id"] = initData.OrgId
	exportArgs["project-id"] = initData.ProjectId
	exportArgs["project-name"] = initData.ProjectName

	return true, tfArgs, costArgs, exportArgs, nil
}

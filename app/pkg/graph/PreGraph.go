package graph

import (
	"fmt"
	"pluralith/pkg/cost"
	"pluralith/pkg/initialization"
	"pluralith/pkg/terraform"

	"github.com/spf13/pflag"
)

func PreGraph(flags *pflag.FlagSet) (bool, map[string]interface{}, map[string]interface{}, map[string]interface{}, error) {
	functionName := "PreGraph"

	tfArgs := terraform.ConstructTerraformArgs(flags)
	costArgs, costErr := cost.ConstructInfracostArgs(flags)
	if costErr != nil {
		return false, nil, nil, nil, fmt.Errorf("%v: %w", functionName, costErr)
	}

	exportArgs := ConstructExportArgs(flags)
	costArgs["show-costs"] = true // Always run infracost if infracost is installed

	// Set defaults for export
	if exportArgs["title"] == "" {
		exportArgs["title"] = "Infrastructure Diagram"
		exportArgs["file-name"] = "Infrastructure_Diagram"
	}

	initValid, _, initErr := initialization.RunInit(true, initialization.InitData{}, true)
	if initErr != nil {
		return false, nil, nil, nil, fmt.Errorf("%v: %w", functionName, initErr)
	}
	if !initValid {
		return false, nil, nil, nil, nil
	}

	exportArgs["org-id"] = ""
	exportArgs["project-id"] = "pluralith-local-project"
	exportArgs["project-name"] = "Local Runs"

	return true, tfArgs, costArgs, exportArgs, nil
}

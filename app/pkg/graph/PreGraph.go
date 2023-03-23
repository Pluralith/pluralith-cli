package graph

import (
	"fmt"
	"math/rand"
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
	exportArgs["runId"] = fmt.Sprintf("%07d", rand.Intn(10000000)) // Generate random run id
	costArgs["show-costs"] = true                                  // Always run infracost if infracost is installed

	// Export pdf if local-only flag set
	var localOnly, _ = flags.GetBool("local-only")
	if localOnly {
		exportArgs["export-pdf"] = true
	} else {
		exportArgs["export-pdf"] = false
	}

	// Set defaults for export
	if exportArgs["title"] == "" {
		exportArgs["title"] = "Pluralith Diagram"     //+ exportArgs["runId"].(string)
		exportArgs["file-name"] = "Pluralith_Diagram" //+ exportArgs["runId"].(string)
	}

	initValid, _, initErr := initialization.RunInit(true, initialization.InitData{}, true)
	if initErr != nil {
		return false, nil, nil, nil, fmt.Errorf("%v: %w", functionName, initErr)
	}
	if !initValid {
		return false, nil, nil, nil, nil
	}

	exportArgs["orgId"] = ""
	exportArgs["projectId"] = "pluralith-local-project"
	exportArgs["projectName"] = "Local Runs"

	return true, tfArgs, costArgs, exportArgs, nil
}

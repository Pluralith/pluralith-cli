package cost

import (
	"fmt"
	"pluralith/pkg/auxiliary"

	"github.com/spf13/pflag"
)

func ConstructInfracostArgs(flags *pflag.FlagSet) ([]string, error) {
	functionName := "ConstructInfracostArgs"

	costArgs := []string{}

	// No Costs flag
	noCosts, noCostsErr := flags.GetBool("no-costs")
	if noCostsErr != nil {
		return costArgs, fmt.Errorf("failed to parse no-costs flag -> %v: %w", functionName, noCostsErr)
	}

	if noCosts {
		auxiliary.StateInstance.Infracost = false
		return costArgs, nil
	}

	// Usage file flag
	usageFilePath, usageFileErr := flags.GetString("cost-usage-file")
	if usageFileErr != nil {
		return costArgs, fmt.Errorf("failed to parse cost usage file flag -> %v: %w", functionName, noCostsErr)
	}

	if usageFilePath != "" {
		costArgs = append(costArgs, "--usage-file="+usageFilePath)
	}

	return costArgs, nil
}

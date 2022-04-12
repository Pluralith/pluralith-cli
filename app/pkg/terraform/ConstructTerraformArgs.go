package terraform

import (
	"fmt"

	"github.com/spf13/pflag"
)

func ConstructTerraformArgs(flags *pflag.FlagSet) ([]string, error) {
	functionName := "ConstructTerraformArgs"

	tfArgs := []string{}

	// Get values from flags
	vars, varsErr := flags.GetStringSlice("var")
	if varsErr != nil {
		return tfArgs, fmt.Errorf("failed to parse var flags -> %v: %w", functionName, varsErr)
	}

	varFiles, varFilesErr := flags.GetStringSlice("var-file")
	if varFilesErr != nil {
		return tfArgs, fmt.Errorf("failed to parse var file flags -> %v: %w", functionName, varFilesErr)
	}

	// Construct arg slice for terraform
	for _, varValue := range vars {
		tfArgs = append(tfArgs, "-var='"+varValue+"'")
	}

	for _, varFile := range varFiles {
		tfArgs = append(tfArgs, "-var-file="+varFile)
	}

	return tfArgs, nil
}

package terraform

import (
	"fmt"

	"github.com/spf13/pflag"
)

func ConstructVarArgs(flags *pflag.FlagSet) ([]string, error) {
	functionName := "ConstructVarArgs"

	varArgs := []string{}

	// Get values from flags
	vars, varsErr := flags.GetStringSlice("var")
	if varsErr != nil {
		return varArgs, fmt.Errorf("failed to pars var flags -> %v: %w", functionName, varsErr)
	}

	varFiles, varFilesErr := flags.GetStringSlice("var-file")
	if varFilesErr != nil {
		return varArgs, fmt.Errorf("failed to pars var file flags -> %v: %w", functionName, varFilesErr)
	}

	// Construct arg slice for terraform
	for _, varValue := range vars {
		varArgs = append(varArgs, "-var='"+varValue+"'")
	}

	for _, varFile := range varFiles {
		varArgs = append(varArgs, "-var-file="+varFile)
	}

	return varArgs, nil
}

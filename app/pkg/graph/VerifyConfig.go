package graph

import (
	"fmt"
	"pluralith/pkg/auth"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func VerifyConfig(noProject bool) (bool, error) {
	functionName := "VerifyConfig"

	ux.PrintFormatted("→ ", []string{"blue", "bold"})
	ux.PrintFormatted("Verify", []string{"white", "bold"})
	fmt.Println()

	// Verify API key with backend
	apiKeyValid, apiKeyErr := auth.VerifyAPIKey(auxiliary.StateInstance.APIKey, false)
	if !apiKeyValid {
		return false, nil
	}
	if apiKeyErr != nil {
		return false, fmt.Errorf("verifying API key failed -> %v: %w", functionName, apiKeyErr)
	}

	if !noProject {
		projectValid, projectErr := auth.VerifyProject(auxiliary.StateInstance.PluralithConfig.ProjectId)
		if !projectValid {
			return false, nil
		}
		if projectErr != nil {
			return false, fmt.Errorf("failed to verify project id -> %v: %w", functionName, projectErr)
		}
	}

	return true, nil
}

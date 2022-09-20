package graph

import (
	"fmt"
	"pluralith/pkg/auth"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func VerifyConfig(noProject bool) (bool, map[string]interface{}, error) {
	functionName := "VerifyConfig"

	ux.PrintFormatted("→ ", []string{"blue", "bold"})
	ux.PrintFormatted("Verify\n", []string{"white", "bold"})
	// fmt.Println()

	// Verify API key with backend
	apiKeyValid, apiKeyErr := auth.VerifyAPIKey(auxiliary.StateInstance.APIKey, false)
	if !apiKeyValid {
		return false, nil, nil
	}
	if apiKeyErr != nil {
		return false, nil, fmt.Errorf("verifying API key failed -> %v: %w", functionName, apiKeyErr)
	}

	if !noProject {
		projectData, projectErr := auth.VerifyProject(auxiliary.StateInstance.PluralithConfig.ProjectId)
		if projectData == nil {
			return false, nil, nil
		}
		if projectErr != nil {
			return false, nil, fmt.Errorf("failed to verify project id -> %v: %w", functionName, projectErr)
		}

		return true, projectData, nil
	}

	return true, nil, nil
}

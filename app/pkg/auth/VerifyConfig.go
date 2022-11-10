package auth

import (
	"fmt"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func VerifyConfig(noProject bool) (bool, OrgResponse, error) {
	functionName := "VerifyConfig"
	orgData := OrgResponse{}

	ux.PrintFormatted("→ ", []string{"blue", "bold"})
	ux.PrintFormatted("Verify\n", []string{"white", "bold"})

	// Verify API key with backend
	apiKeyValid, apiKeyErr := VerifyAPIKey(auxiliary.StateInstance.APIKey, false)
	if !apiKeyValid {
		return false, orgData, nil
	}
	if apiKeyErr != nil {
		return false, orgData, fmt.Errorf("verifying API key failed -> %v: %w", functionName, apiKeyErr)
	}

	if !noProject {
		orgData, orgErr := VerifyOrg(auxiliary.StateInstance.PluralithConfig.OrgId)
		if orgData.Data.ID == "" {
			return false, orgData, fmt.Errorf("no project data given -> %v", functionName)
		}
		if orgErr != nil {
			return false, orgData, fmt.Errorf("failed to verify project id -> %v: %w", functionName, orgErr)
		}
	}

	return true, orgData, nil
}

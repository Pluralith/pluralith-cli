package graph

import (
	"fmt"
	"pluralith/pkg/auth"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func VerifyConfig(noProject bool) error {
	functionName := "VerifyConfig"

	ux.PrintFormatted("→ ", []string{"blue", "bold"})
	ux.PrintFormatted("Verify", []string{"white", "bold"})
	fmt.Println()

	// Verify API key with backend
	apiKeyValid, apiKeyErr := auth.VerifyAPIKey(auxiliary.StateInstance.APIKey, false)
	if !apiKeyValid {
		return nil
	}
	if apiKeyErr != nil {
		return fmt.Errorf("verifying API key failed -> %v: %w", functionName, apiKeyErr)
	}

	if !noProject {
		projectValid, projectErr := auth.VerifyProject(auxiliary.StateInstance.PluralithConfig.ProjectId)
		if !projectValid {
			return nil
		}
		if projectErr != nil {
			return fmt.Errorf("failed to verify project id -> %v: %w", functionName, projectErr)
		}
	}

	// if !isValid {
	// 	ux.PrintFormatted("\n✘", []string{"red", "bold"})
	// 	fmt.Print(" Invalid API key → Run ")
	// 	ux.PrintFormatted("pluralith login", []string{"blue"})
	// 	fmt.Println(" again\n")
	// 	return nil
	// }

	return nil
}

package plan

import (
	"fmt"
	"pluralith/pkg/auxiliary"
)

func FetchProviders(jsonString string) ([]string, error) {
	functionName := "FetchProviders"

	// Parse JSON object from string
	parsedJson, parseErr := auxiliary.ParseJson(jsonString)
	if parseErr != nil {
		return make([]string, 0), fmt.Errorf("%v: %w", functionName, parseErr)
	}

	// Get provider config
	configuration := parsedJson["configuration"].(map[string]interface{})
	provider_config := configuration["provider_config"].(map[string]interface{})

	// Find all used providers
	providers := make([]string, 0, len(provider_config))
	for item := range provider_config {
		providers = append(providers, item)
	}

	return providers, nil
}

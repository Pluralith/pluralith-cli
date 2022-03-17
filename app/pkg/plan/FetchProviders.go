package plan

import (
	"fmt"
	"pluralith/pkg/auxiliary"
)

func FetchProviders(jsonString string) ([]string, error) {
	functionName := "FetchProviders"

	// Initialize provider array
	var providers = []string{}

	// Parse JSON object from string
	parsedJson, parseErr := auxiliary.ParseJson(jsonString)
	if parseErr != nil {
		return make([]string, 0), fmt.Errorf("%v: %w", functionName, parseErr)
	}

	// Get provider config
	configuration := parsedJson["configuration"].(map[string]interface{})

	if configuration["provider_config"] != nil {
		provider_config := configuration["provider_config"].(map[string]interface{})

		// Find all used providers
		for item := range provider_config {
			providers = append(providers, item)
		}

		return providers, nil
	}

	// If provider_config missing -> Fall back to fetching providers from resources
	rootModule := configuration["root_module"].(map[string]interface{})
	resources := rootModule["resources"].([]interface{})

	for _, resource := range resources {
		resourceMap := resource.(map[string]interface{})
		providerName := resourceMap["provider_config_key"].(string)
		if !auxiliary.ElementInSlice(providerName, providers) {
			providers = append(providers, providerName)
		}
	}

	return providers, nil
}

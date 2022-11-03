package backends

import (
	"fmt"
)

type TerraformState struct {
	Version int    `json:"version"`
	Serial  int    `json:"serial"`
	Lineage string `json:"lineage"`
	Backend struct {
		Type   string      `json:"type"`
		Config interface{} `json:"config"`
		Hash   int         `json:"hash"`
	} `json:"backend"`
	Modules []struct {
		Path      []string      `json:"path"`
		Outputs   struct{}      `json:"outputs"`
		Resources struct{}      `json:"resources"`
		DependsOn []interface{} `json:"depends_on"`
	} `json:"modules"`
}

func StoreInBackend() error {
	functionName := "PushDiagramToBackend"

	backendConfig, backendErr := LoadBackendConfig()
	if backendErr != nil {
		return fmt.Errorf("could not load backend -> %v: %w", functionName, backendErr)
	}

	// Detect which backend is uses
	if backendConfig.Backend.Type == "s3" {
		if awsErr := PushToAWSBackend(backendConfig); awsErr != nil {
			return fmt.Errorf("failed to push to aws backend -> %v: %w", functionName, awsErr)
		}
	}
	if backendConfig.Backend.Type == "azurerm" {
		if azureErr := PushToAzureBackend(backendConfig); azureErr != nil {
			return fmt.Errorf("failed to push to azure backend -> %v: %w", functionName, azureErr)
		}
	}
	if backendConfig.Backend.Type == "gcs" {
		if googleErr := PushToAzureBackend(backendConfig); googleErr != nil {
			return fmt.Errorf("failed to push to google backend -> %v: %w", functionName, googleErr)
		}
	}

	// fmt.Println(config)
	// TODOs here:
	// - Detect state backend
	// - Get credentials
	// - Upload file to backend
	return nil
}

package backends

import (
	"fmt"
	"pluralith/pkg/ux"
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
		if awsErr := PushToS3Backend(backendConfig); awsErr != nil {
			return fmt.Errorf("failed to push to aws backend -> %v: %w", functionName, awsErr)
		}
		return nil
	}
	if backendConfig.Backend.Type == "azurerm" {
		if azureErr := PushToAzureBackend(backendConfig); azureErr != nil {
			return fmt.Errorf("failed to push to azure backend -> %v: %w", functionName, azureErr)
		}
		return nil
	}
	if backendConfig.Backend.Type == "gcs" {
		if googleErr := PushToGCSBackend(backendConfig); googleErr != nil {
			return fmt.Errorf("failed to push to gcs backend -> %v: %w", functionName, googleErr)
		}
		return nil
	}

	if backendConfig.Backend.Type != "" {
		ux.PrintFormatted("  ✘", []string{"yellow", "bold"})
		fmt.Print(" Backend type ")
		ux.PrintFormatted(backendConfig.Backend.Type, []string{"white", "bold"})
		fmt.Println(" isn't supported. Pluralith currently supports s3, gcs and azurerm.")
	}

	return nil
}

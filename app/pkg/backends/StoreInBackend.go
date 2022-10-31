package backends

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
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

func LoadBackendConfig() (interface{}, error) {
	functionName := "LoadBackendConfig"

	// Read terraform state for backend information
	tfStatePath := filepath.Join(auxiliary.StateInstance.WorkingPath, ".terraform", "terraform.tfstate")
	tfStateByte, configErr := os.ReadFile(tfStatePath)
	if configErr != nil {
		return nil, fmt.Errorf("failed to read working directory config -> %v: %w", functionName, configErr)
	}

	// Parse terraform state
	tfState := TerraformState{}
	jsonErr := json.Unmarshal(tfStateByte, &tfState)
	if jsonErr != nil {
		return nil, fmt.Errorf("failed to unmarshal terraform state information -> %v: %w", functionName, jsonErr)
	}

	// Detect which backend is uses
	if tfState.Backend.Type == "aws" {
		// azureConfig := LoadAzureBackendConfig()
		// return azureConfig, nil
	}

	if tfState.Backend.Type == "azurerm" {
		azureConfig := LoadAzureBackendConfig(tfState)
		return azureConfig, nil
	}

	if tfState.Backend.Type == "google" {
		// azureConfig := LoadAzureBackendConfig()
		// return azureConfig, nil
	}

	return nil, nil
}

func PushDiagramToBackend() error {
	config, err := LoadBackendConfig()
	if err != nil {
		return err
	}

	fmt.Println(config)
	// TODOs here:
	// - Detect state backend
	// - Get credentials
	// - Upload file to backend
	return nil
}

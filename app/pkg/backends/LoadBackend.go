package backends

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"

	"github.com/mitchellh/mapstructure"
)

func MapBackendConfig(tfState TerraformState, configObject interface{}) error {
	functionName := "LoadBackendConfig"

	config := &mapstructure.DecoderConfig{TagName: "json"}
	config.Result = &configObject

	decoder, decodeInitErr := mapstructure.NewDecoder(config)
	if decodeInitErr != nil {
		return fmt.Errorf("creating map decoder failed -> %v: %w", functionName, decodeInitErr)
	}

	decodeErr := decoder.Decode(tfState.Backend.Config)
	if decodeInitErr != nil {
		return fmt.Errorf("decoding backend config failed -> %v: %w", functionName, decodeErr)
	}

	return nil
}

func LoadBackendConfig() (TerraformState, error) {
	functionName := "LoadBackendConfig"
	tfState := TerraformState{}

	// Read terraform state for backend information
	tfStatePath := filepath.Join(auxiliary.StateInstance.WorkingPath, ".terraform", "terraform.tfstate")
	tfStateByte, configErr := os.ReadFile(tfStatePath)
	if configErr != nil {
		return tfState, fmt.Errorf("failed to read working directory config -> %v: %w", functionName, configErr)
	}

	// Parse terraform state
	jsonErr := json.Unmarshal(tfStateByte, &tfState)
	if jsonErr != nil {
		return tfState, fmt.Errorf("failed to unmarshal terraform state information -> %v: %w", functionName, jsonErr)
	}

	return tfState, nil
}

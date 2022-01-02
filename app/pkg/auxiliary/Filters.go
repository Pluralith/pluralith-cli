package auxiliary

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Filters struct {
	Replacement string
	Config      SecretConfig
}

type SecretConfig struct {
	Sensitive []string
}

func (F *Filters) GetSecretConfig() error {
	functionName := "GetSecretConfig"

	// Initialize variables
	var configByte []byte
	var configErr error
	var config SecretConfig

	// Get relevant paths to read config from
	workingConfig := filepath.Join(PathInstance.WorkingPath, "pluralith-config.json")
	defaultConfig := filepath.Join(PathInstance.HomePath, "Pluralith", "pluralith-config.json")

	// Read config file from working directory
	configByte, configErr = os.ReadFile(workingConfig)
	if configErr != nil {
		configErr = nil
		// If no config file in working directroy -> fall back to default config
		configByte, configErr = os.ReadFile(defaultConfig)
		if configErr != nil {
			return fmt.Errorf("failed to get config -> %v: %w", functionName, configErr)
		}
	}

	// Parse config
	parseErr := json.Unmarshal(configByte, &config)
	if parseErr != nil {
		return fmt.Errorf("failed to parse config -> %v: %w", functionName, parseErr)
	}

	// Set config for global access
	FilterInstance.Config = config

	return nil
}

func (F *Filters) InitializeFilters() error {
	// functionName := "InitializeFilters"
	F.Replacement = "gatewatch"

	return nil
}

var FilterInstance = &Filters{}

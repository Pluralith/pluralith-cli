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
	workingConfig := filepath.Join(StateInstance.WorkingPath, "pluralith-config.json")
	defaultConfig := filepath.Join(StateInstance.HomePath, "Pluralith", "pluralith-config.json")

	// Get default config first
	if _, statErr := os.Stat(defaultConfig); !os.IsNotExist(statErr) {
		// Read config file from Pluralith directory
		configByte, configErr = os.ReadFile(defaultConfig)
		if configErr != nil {
			return fmt.Errorf("failed to read working directory config -> %v: %w", functionName, configErr)
		}
	}

	// If current working dir has config -> override default config
	if _, statErr := os.Stat(workingConfig); !os.IsNotExist(statErr) {
		// Read config file from working directory
		configByte, configErr = os.ReadFile(workingConfig)
		if configErr != nil {
			return fmt.Errorf("failed to read working directory config -> %v: %w", functionName, configErr)
		}
	}

	// Parse config if given
	if len(configByte) > 0 {
		parseErr := json.Unmarshal(configByte, &config)
		if parseErr != nil {
			return fmt.Errorf("failed to parse config -> %v: %w", functionName, parseErr)
		}
	}

	// Set config for global access
	FilterInstance.Config = config

	return nil
}

func (F *Filters) InitFilters() error {
	// functionName := "InitializeFilters"
	F.Replacement = "gatewatch"

	return nil
}

var FilterInstance = &Filters{}

package auxiliary

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type PluralithConfig struct {
	ProjectId string `yaml:"project_id"`
	Config    struct {
		SensitiveAttrs []string `yaml:"sensitive_attrs"`
		Vars           []string `yaml:"vars"`
		VarFiles       []string `yaml:"var_files"`
		CostUsageFile  string   `yaml:"cost_usage_file"`
	} `yaml:"config"`
	Export struct {
		Title   string `yaml:"title"`
		Author  string `yaml:"author"`
		Version string `yaml:"version"`
	} `yaml:"export"`
}

func (S *State) GetConfig() error {
	functionName := "GetSecretConfig"

	// Initialize variables
	var configByte []byte
	var configErr error
	var config PluralithConfig

	// Get relevant paths to read config from
	workingConfig := filepath.Join(StateInstance.WorkingPath, "pluralith.yml")
	defaultConfig := filepath.Join(StateInstance.HomePath, "Pluralith", "pluralith.yml")

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
		yamlErr := yaml.Unmarshal(configByte, &config)
		if yamlErr != nil {
			return fmt.Errorf("failed to parse config -> %v: %w", functionName, yamlErr)
		}
	}

	// Set config for global access
	S.PluralithConfig = config

	return nil
}

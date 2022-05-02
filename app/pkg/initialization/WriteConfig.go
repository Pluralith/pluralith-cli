package initialization

import (
	"fmt"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
	"strconv"
)

func WriteConfig(projectId string) error {
	configPath := filepath.Join(auxiliary.StateInstance.WorkingPath, "pluralith.yml")
	configString := EmtpyConfig

	if projectId != "" {
		projectIdInt, convErr := strconv.Atoi(projectId)
		if convErr != nil {
			return fmt.Errorf("converting project id to number failed -> %w", convErr)
		}
		configString = fmt.Sprintf(ConfigTemplate, projectIdInt)
	}

	helperWriteErr := os.WriteFile(configPath, []byte(configString), 0700)
	if helperWriteErr != nil {
		return fmt.Errorf("failed to create config template -> %w", helperWriteErr)
	}

	if projectId != "" {
		ux.PrintFormatted("✔", []string{"blue", "bold"})
		fmt.Print(" Your project has been initialized! Customize your config in ")
		ux.PrintFormatted("pluralith.yml\n\n", []string{"blue"})
	} else {
		ux.PrintFormatted("✔", []string{"blue", "bold"})
		fmt.Print(" Empty config initialized! Customize your config in ")
		ux.PrintFormatted("pluralith.yml\n\n", []string{"blue"})
	}

	return nil
}

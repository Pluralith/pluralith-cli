package initialization

import (
	"fmt"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
	"strings"
)

func WriteConfig(initData InitData) error {
	configPath := filepath.Join(auxiliary.StateInstance.WorkingPath, "pluralith.yml")
	configString := ConfigTemplate

	// Interpolate values
	configString = strings.ReplaceAll(configString, "$PLR_ORG_ID", initData.OrgId)
	configString = strings.ReplaceAll(configString, "$PLR_PROJECT_ID", initData.ProjectId)
	configString = strings.ReplaceAll(configString, "$PLR_PROJECT_NAME", initData.ProjectName)
	configString = strings.ReplaceAll(configString, "$PLR_API_ENDPOINT", initData.PluralithAPIEndpoint)

	helperWriteErr := os.WriteFile(configPath, []byte(configString), 0700)
	if helperWriteErr != nil {
		return fmt.Errorf("failed to create config template -> %w", helperWriteErr)
	}

	if initData.OrgId != "" {
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

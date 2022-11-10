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

	// Insert org id
	// if initData.OrgId != "" && !isEmpty {
	configString = strings.ReplaceAll(ConfigTemplate, "$PLR_ORG_ID", initData.OrgId)
	// } else {
	// 	configString = strings.ReplaceAll(ConfigTemplate, "$PLR_ORG_ID", "null")
	// }

	// Insert project id
	// if initData.ProjectId != "" && !isEmpty {
	configString = strings.ReplaceAll(ConfigTemplate, "$PLR_PROJECT_ID", initData.ProjectId)
	// } else {
	// 	configString = strings.ReplaceAll(ConfigTemplate, "$PLR_PROJECT_ID", "null")
	// }

	// Insert project name
	// if initData.ProjectName != "" && !isEmpty {
	configString = strings.ReplaceAll(ConfigTemplate, "$PLR_PROJECT_NAME", initData.ProjectName)
	// } else {
	// 	configString = strings.ReplaceAll(ConfigTemplate, "$PLR_PROJECT_NAME", "null")
	// }

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

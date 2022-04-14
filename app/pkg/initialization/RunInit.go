package initialization

import (
	"fmt"
	"pluralith/pkg/auth"
	"pluralith/pkg/ux"
)

func RunInit(isEmpty bool, APIKey string, projectId string) error {
	functionName := "RunInit"

	if !isEmpty && projectId == "" {
		fmt.Println("Lets set up your project and get you up and running.\n")
	}

	if isEmpty {
		fmt.Println()
		if writeErr := WriteConfig(projectId); writeErr != nil {
			return fmt.Errorf("failed to create config template -> %v: %w", functionName, writeErr)
		}
		return nil
	}

	if APIKey == "" {
		ux.PrintFormatted("⠿ ", []string{"blue"})
		fmt.Println("We noticed you are not authenticated!")
		ux.PrintFormatted("→", []string{"blue", "bold"})
		fmt.Print(" Enter your API Key: ")

		// Capture user input
		fmt.Scanln(&APIKey)
		loginValid, loginErr := auth.RunLogin(APIKey)
		if !loginValid {
			return nil
		}
		if loginErr != nil {
			fmt.Println(loginErr)
		}
	}

	if projectId == "" {
		ux.PrintFormatted("→", []string{"blue", "bold"})
		fmt.Print(" Enter Project Id: ")

		// Capture user input
		fmt.Scanln(&projectId)
		fmt.Print("  ") // Formatting gimmick
	}

	if writeErr := WriteConfig(projectId); writeErr != nil {
		return fmt.Errorf("failed to create config template -> %v: %w", functionName, writeErr)
	}

	return nil
}

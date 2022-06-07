package initialization

import (
	"fmt"
	"pluralith/pkg/auth"
	"pluralith/pkg/ux"
)

func RunInit(isEmpty bool, APIKey string, projectId string) error {
	functionName := "RunInit"

	if isEmpty {
		// fmt.Println()
		if writeErr := WriteConfig(projectId); writeErr != nil {
			return fmt.Errorf("failed to create config template -> %v: %w", functionName, writeErr)
		}
		return nil
	}

	ux.PrintFormatted("\n→", []string{"blue", "bold"})
	ux.PrintFormatted(" Authentication\n", []string{"white", "bold"})

	if APIKey == "" {
		ux.PrintFormatted("  ⠿ ", []string{"blue"})
		fmt.Println("We noticed you are not authenticated!")
		ux.PrintFormatted("  →", []string{"blue", "bold"})
		fmt.Print(" Enter your API Key: ")

		// Capture user input
		fmt.Scanln(&APIKey)
	}

	// Run login routine and set credentials file
	loginValid, loginErr := auth.RunLogin(APIKey)
	if !loginValid {
		return nil
	}
	if loginErr != nil {
		return fmt.Errorf("failed to authenticate -> %v: %w", functionName, loginErr)
	}

	ux.PrintFormatted("\n→", []string{"blue", "bold"})
	ux.PrintFormatted(" Project Setup\n", []string{"white", "bold"})

	if projectId == "" {
		ux.PrintFormatted("  →", []string{"blue", "bold"})
		fmt.Print(" Enter Project Id: ")

		// Capture user input
		fmt.Scanln(&projectId)
	}

	projectValid, projectErr := auth.VerifyProject(projectId)
	if !projectValid {
		return nil
	}
	if projectErr != nil {
		return fmt.Errorf("failed to verify project id -> %v: %w", functionName, projectErr)
	}

	fmt.Print("  ") // Formatting gimmick

	if writeErr := WriteConfig(projectId); writeErr != nil {
		return fmt.Errorf("failed to create config template -> %v: %w", functionName, writeErr)
	}

	return nil
}

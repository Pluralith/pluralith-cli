package initialization

import (
	"fmt"
	"pluralith/pkg/auth"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func compileInitData(initData InitData) InitData {
	// functionName := "compileInitData"

	// Set init data variables from config or env variables if none passed from flags
	if initData.APIKey == "" {
		initData.APIKey = auxiliary.StateInstance.APIKey
	}
	if initData.OrgId == "" {
		initData.OrgId = auxiliary.StateInstance.PluralithConfig.OrgId
	}
	if initData.ProjectId == "" {
		initData.ProjectId = auxiliary.StateInstance.PluralithConfig.ProjectId
	}
	if initData.ProjectName == "" {
		initData.ProjectName = auxiliary.StateInstance.PluralithConfig.ProjectName
	}

	return initData
}

func RunInit(isEmpty bool, initData InitData) (InitData, error) {
	functionName := "RunInit"

	if isEmpty {
		if writeErr := WriteConfig(initData); writeErr != nil {
			return initData, fmt.Errorf("failed to create config template -> %v: %w", functionName, writeErr)
		}
		return initData, nil
	}

	// Compile init data from various sources
	initData = compileInitData(initData)

	// Handle user inputs

	// Authentication
	ux.PrintFormatted("\n→", []string{"blue", "bold"})
	ux.PrintFormatted(" Authentication\n", []string{"white", "bold"})
	if initData.APIKey == "" {
		fmt.Print("  Enter API Key: ")
		fmt.Scanln(&initData.APIKey) // Capture user input
	}

	// Run login routine and set credentials file
	loginValid, loginErr := auth.RunLogin(initData.APIKey)
	if !loginValid {
		return initData, nil
	}
	if loginErr != nil {
		return initData, fmt.Errorf("failed to authenticate -> %v: %w", functionName, loginErr)
	}

	// Project Setup
	ux.PrintFormatted("\n→", []string{"blue", "bold"})
	ux.PrintFormatted(" Project Setup\n", []string{"white", "bold"})
	if initData.OrgId == "" {
		fmt.Print("  Enter Org Id: ")
		fmt.Scanln(&initData.OrgId) // Capture user input
	}
	auth.VerifyOrg(initData.OrgId)

	if initData.ProjectId == "" {
		fmt.Print("  Enter Project Id: ")
		fmt.Scanln(&initData.ProjectId) // Capture user input
	}
	auth.VerifyProject(initData.OrgId, initData.ProjectId)

	if initData.ProjectName == "" {
		fmt.Print("  Enter Project Name: ")
		fmt.Scanln(&initData.ProjectName) // Capture user input
	}

	// // ask for input if

	// if initData.APIKey == "" {
	// 	ux.PrintFormatted("  ⠿ ", []string{"blue"})
	// 	fmt.Println("We noticed you are not authenticated!")
	// 	ux.PrintFormatted("  →", []string{"blue", "bold"})
	// 	fmt.Print(" Enter your API Key: ")

	// 	// Capture user input
	// 	fmt.Scanln(&initData.APIKey)
	// }

	// // Run login routine and set credentials file
	// loginValid, loginErr := auth.RunLogin(initData.APIKey)
	// if !loginValid {
	// 	return nil
	// }
	// if loginErr != nil {
	// 	return fmt.Errorf("failed to authenticate -> %v: %w", functionName, loginErr)
	// }

	// ux.PrintFormatted("\n→", []string{"blue", "bold"})
	// ux.PrintFormatted(" Project Setup\n", []string{"white", "bold"})

	// if orgId == "" {
	// 	ux.PrintFormatted("  →", []string{"blue", "bold"})
	// 	fmt.Print(" Enter Org Id: ")

	// 	// Capture user input
	// 	fmt.Scanln(&orgId)
	// }

	// if projectId == "" {
	// 	ux.PrintFormatted("  →", []string{"blue", "bold"})
	// 	fmt.Print(" Enter Project Id: ")

	// 	// Capture user input
	// 	fmt.Scanln(&projectId)
	// }

	// orgData, projectErr := auth.VerifyOrg(orgId)
	// if orgData == nil {
	// 	return nil
	// }

	// projectData, projectErr := auth.VerifyProject(orgId, projectId)
	// if projectData == nil {
	// 	return nil
	// }
	// if projectErr != nil {
	// 	return fmt.Errorf("failed to verify project id -> %v: %w", functionName, projectErr)
	// }

	// fmt.Print("  ") // Formatting gimmick

	// if writeErr := WriteConfig(projectId); writeErr != nil {
	// 	return fmt.Errorf("failed to create config template -> %v: %w", functionName, writeErr)
	// }

	return initData, nil
}

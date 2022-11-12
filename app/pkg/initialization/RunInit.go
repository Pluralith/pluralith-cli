package initialization

import (
	"bufio"
	"fmt"
	"os"
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

func RunInit(askInputs bool, initData InitData) (InitData, error) {
	functionName := "RunInit"

	// Compile init data from various sources
	initData = compileInitData(initData)

	// Authentication
	ux.PrintFormatted("\n→", []string{"blue", "bold"})
	ux.PrintFormatted(" Authentication\n", []string{"white", "bold"})
	if initData.APIKey == "" && askInputs {
		ux.PrintFormatted("  ⠿", []string{"blue", "bold"})
		fmt.Print(" Enter API Key: ")
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
	if initData.OrgId == "" && askInputs {
		ux.PrintFormatted("  ⠿", []string{"blue", "bold"})
		fmt.Print(" Enter Org Id: ")
		fmt.Scanln(&initData.OrgId) // Capture user input
	}

	orgFound, orgErr := VerifyOrg(initData.OrgId)
	if orgErr != nil {
		return initData, fmt.Errorf("failed to verify org id -> %v: %w", functionName, orgErr)
	}
	if !orgFound {
		return initData, nil
	}

	if initData.ProjectId == "" && askInputs {
		ux.PrintFormatted("\n  ⠿", []string{"blue", "bold"})
		fmt.Print(" Enter Project Id: ")
		fmt.Scanln(&initData.ProjectId) // Capture user input
	}

	projectValid, projectName, projectErr := VerifyProject(initData.OrgId, initData.ProjectId)
	if projectErr != nil {
		return initData, fmt.Errorf("failed to verify org id -> %v: %w", functionName, projectErr)
	}

	// Handle non-existent project
	if projectValid {
		initData.ProjectName = projectName // Set name in init data if existing project is found

		if initData.ProjectName == "" && askInputs { // If at this point project name is still empty and command run is pluralith init -> ask for user's input
			ux.PrintFormatted("\n  ⠿", []string{"blue", "bold"})
			fmt.Print(" Enter Project Name: ")

			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				initData.ProjectName = scanner.Text() // Capture user input
			}
		}

		if projectName == "" { // If project is not in existence -> create project
			CreateProject(initData)
		}
	}

	// request, _ := http.NewRequest("GET", "https://api.pluralith.com/v1/project/get", nil)

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

	fmt.Println()
	if askInputs {
		if writeErr := WriteConfig(initData); writeErr != nil {
			return initData, fmt.Errorf("failed to create config template -> %v: %w", functionName, writeErr)
		}
	}

	return initData, nil
}

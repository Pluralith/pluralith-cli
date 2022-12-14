package initialization

import (
	"bufio"
	"fmt"
	"os"
	"pluralith/pkg/auth"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func CompileInitData(initData InitData) InitData {
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

func RunInit(noInputs bool, initData InitData, localRun bool) (bool, InitData, error) {
	functionName := "RunInit"

	// Compile init data from various sources
	initData = CompileInitData(initData)

	// Authentication
	ux.PrintFormatted("\n→", []string{"blue", "bold"})
	ux.PrintFormatted(" Authentication\n", []string{"white", "bold"})
	if initData.APIKey == "" && !noInputs {
		ux.PrintFormatted("  ⠿", []string{"blue", "bold"})
		fmt.Print(" Enter API Key: ")
		fmt.Scanln(&initData.APIKey) // Capture user input
	}

	// Run login routine and set credentials file
	loginValid, loginErr := auth.RunLogin(initData.APIKey)
	if loginErr != nil {
		return false, initData, fmt.Errorf("failed to authenticate -> %v: %w", functionName, loginErr)
	}
	if !loginValid {
		return false, initData, nil
	}

	if localRun {
		return true, initData, nil 
	}

	// Project Setup
	ux.PrintFormatted("\n→", []string{"blue", "bold"})
	ux.PrintFormatted(" Project Setup\n", []string{"white", "bold"})
	if initData.OrgId == "" && !noInputs {
		ux.PrintFormatted("  ⠿", []string{"blue", "bold"})
		fmt.Print(" Enter Org Id: ")
		fmt.Scanln(&initData.OrgId) // Capture user input
	}

	orgFound, orgErr := VerifyOrg(initData.OrgId)
	if orgErr != nil {
		return false, initData, fmt.Errorf("failed to verify org id -> %v: %w", functionName, orgErr)
	}
	if !orgFound {
		return false, initData, nil
	}

	if initData.ProjectId == "" && !noInputs {
		ux.PrintFormatted("\n  ⠿", []string{"blue", "bold"})
		fmt.Print(" Enter Project Id: ")

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			initData.ProjectId = "" + scanner.Text() // Capture user input
		}
	}

	projectValid, projectName, projectErr := VerifyProject(initData.OrgId, initData.ProjectId)
	if projectErr != nil {
		return false, initData, fmt.Errorf("failed to verify org id -> %v: %w", functionName, projectErr)
	}

	// Set name in init data if existing project is found
	if projectName != "" {
		initData.ProjectName = projectName
	}

	// Handle non-existent project
	if projectValid && projectName == "" {
		// If at this point project name is still empty and command run is pluralith init -> ask for user's input
		if initData.ProjectName == "" {
			if !noInputs {
				ux.PrintFormatted("\n  ⠿", []string{"blue", "bold"})
				fmt.Print(" Enter Project Name: ")

				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					initData.ProjectName = scanner.Text() // Capture user input
				}
			} else {
				ux.PrintFormatted("  ✘", []string{"red", "bold"})
				fmt.Println(" No Project Name Given → Pass A Name To Create A New Project")
			}
		}

		CreateProject(initData)
	}

	if !noInputs {
		fmt.Println()
		if writeErr := WriteConfig(initData); writeErr != nil {
			return false, initData, fmt.Errorf("failed to create config template -> %v: %w", functionName, writeErr)
		}
	}

	return true, initData, nil
}

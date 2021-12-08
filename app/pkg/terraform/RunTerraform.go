package terraform

import (
	"fmt"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/dblock"
	"pluralith/pkg/ux"
)

func RunTerraform(command string, args []string) error {
	functionName := "RunTerraform"

	// Initialize various components of application
	dblock.LockInstance.GenerateLock()
	auxiliary.PathInstance.GeneratePaths()
	auxiliary.FilterInstance.InitializeFilters()

	// Print running message
	ux.PrintFormatted("⠿", []string{"blue"})
	fmt.Println(RunMessages[command].([]string)[0])

	// Manually parse arg (due to cobra lacking a feature)
	parsedArgs, parsedArgMap := auxiliary.ParseArgs(args, []string{})

	// Add necessary flags if not already given
	if parsedArgMap["auto-approve"] == "" {
		parsedArgs = append(parsedArgs, "-auto-approve")
	}
	if parsedArgMap["json"] == "" {
		parsedArgs = append(parsedArgs, "-json")
	}

	// Launching Pluralith
	auxiliary.LaunchPluralith()

	// Run terraform plan to create execution plan
	planPath, planErr := RunPlan(command)
	if planErr != nil {
		return fmt.Errorf("running terraform plan failed -> %v: %w", functionName, planErr)
	}

	// Add plan path to arguments to run apply on already created execution plan
	parsedArgs = append(parsedArgs, planPath)

	// Run terraform apply on existing execution plan
	applyErr := RunApply(command, parsedArgs)
	if applyErr != nil {
		return fmt.Errorf("running terraform apply failed -> %v: %w", functionName, applyErr)
	}

	return nil
}

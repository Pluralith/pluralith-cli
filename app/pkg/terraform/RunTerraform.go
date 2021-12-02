package terraform

import (
	"fmt"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/dblock"
	"pluralith/pkg/ux"
)

func RunTerraform(command string, args []string) error {
	// Print running message
	ux.PrintFormatted("⠿", []string{"blue"})
	fmt.Println(RunMessages[command].([]string)[0])

	// Initialize comDB locking for this process instance
	dblock.LockInstance.GenerateLock()

	// Manually parse arg (due to cobra lacking a feature)
	parsedArgs, parsedArgMap := auxiliary.ParseArgs(args, []string{})

	// Add necessary flags if not already given
	if parsedArgMap["auto-approve"] == "" {
		parsedArgs = append(parsedArgs, "-auto-approve")
	}
	if parsedArgMap["json"] == "" {
		parsedArgs = append(parsedArgs, "-json")
	}

	// Run terraform plan to create execution plan
	planPath, planErr := RunPlan(command)
	if planErr != nil {
		return planErr
	}

	// Add plan path to arguments to run apply on already created execution plan
	parsedArgs = append(parsedArgs, planPath)

	// Run terraform apply on existing execution plan
	applyErr := RunApply(command, parsedArgs)
	if applyErr != nil {
		return applyErr
	}

	return nil
}

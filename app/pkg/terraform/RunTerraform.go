package terraform

import (
	"fmt"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func RunTerraform(command string, args []string) error {
	functionName := "RunTerraform"

	// Check if "terraform init" has been run
	if !auxiliary.StateInstance.TerraformInit {
		ux.PrintHead()
		ux.PrintFormatted("⠿", []string{"blue"})
		fmt.Print(" No Terraform Initialization found ⇢ Run ")
		ux.PrintFormatted("'terraform init'", []string{"blue", "bold"})
		fmt.Println(" first\n")
		return nil
	}

	// Print running message
	ux.PrintFormatted("⠿", []string{"blue"})
	fmt.Println(RunMessages[command].([]string)[0])

	// Remove old Pluralith state
	removeErr := auxiliary.RemoveOldState()
	if removeErr != nil {
		return fmt.Errorf("deleting old Pluralith state failed -> %v: %w", functionName, removeErr)
	}

	// Launch Pluralith
	launchErr := auxiliary.LaunchPluralith()
	if launchErr != nil {
		return fmt.Errorf("launching Pluralith failed -> %v: %w", functionName, launchErr)
	}

	// Run terraform plan to create execution plan
	planPath, planErr := RunPlan(command, args, false)
	if planErr != nil {
		return fmt.Errorf("running terraform plan failed -> %v: %w", functionName, planErr)
	}

	fmt.Println() // Line separation between plan and apply message prints

	// Add plan path to arguments to run apply on already created execution plan
	args = append(args, planPath)

	// Run terraform apply on existing execution plan
	applyErr := RunApply(command, args)
	if applyErr != nil {
		return fmt.Errorf("running terraform apply failed -> %v: %w", functionName, applyErr)
	}

	return nil
}

package terraform

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/cost"
	"pluralith/pkg/ux"
)

func RunTerraform(command string, tfArgs []string, costArgs []string) error {
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

	// Create Pluralith helper directory (.pluralith)
	_, existErr := os.Stat(filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith"))
	if errors.Is(existErr, os.ErrNotExist) {
		// Create file if it doesn't exist yet
		if mkErr := os.Mkdir(filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith"), 0700); mkErr != nil {
			return fmt.Errorf("%v: %w", functionName, mkErr)
		}
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
	// launchErr := auxiliary.LaunchPluralith()
	// if launchErr != nil {
	// 	return fmt.Errorf("launching Pluralith failed -> %v: %w", functionName, launchErr)
	// }

	// Run terraform plan to create execution plan
	planPath, planErr := RunPlan(command, tfArgs, false)
	if planErr != nil {
		return fmt.Errorf("running terraform plan failed -> %v: %w", functionName, planErr)
	}

	// Run infracost
	if auxiliary.StateInstance.Infracost {
		if costErr := cost.CalculateCost(costArgs); costErr != nil {
			fmt.Println(costErr)
		}
	} else {
		// ux.PrintFormatted("→ ", []string{"bold", "blue"})
		// ux.PrintFormatted("Plan\n", []string{"bold", "white"})
		ux.PrintFormatted("  -", []string{"blue", "bold"})
		fmt.Println(" Cost Calculation Skipped\n")
	}

	fmt.Println() // Line separation between plan and apply message prints

	// Run terraform apply on existing execution plan
	applyErr := RunApply(command, planPath)
	if applyErr != nil {
		return fmt.Errorf("running terraform apply failed -> %v: %w", functionName, applyErr)
	}

	return nil
}

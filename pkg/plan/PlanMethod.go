package plan

import (
	"fmt"

	ux "pluralith/pkg/ux"
)

func PlanMethod(args []string, silent bool) (string, error) {
	if !silent {
		ux.PrintFormatted("⠿", []string{"blue"})
		fmt.Println(" Running plan ⇢ Inspect it in the Pluralith UI\n")
	}

	// Instantiate spinners
	planSpinner := ux.NewSpinner("Generating Execution Plan", "Execution Plan Generated", "Couldn't Generate Execution Plan")
	stripSpinner := ux.NewSpinner("Stripping Secrets", "Secrets Stripped", "Stripping Secrets Failed")

	// Run terraform plan
	planSpinner.Start()
	planExecutionPath, planExecutionErr := GeneratePlan(args)
	if planExecutionErr != nil {
		planSpinner.Fail("Terraform Plan Failed")
		return "", planExecutionErr
	}

	planSpinner.Success("Execution Plan Generated")

	// Run terraform show + strip secrets
	stripSpinner.Start()
	_, planJsonErr := GenerateJson(planExecutionPath)
	if planJsonErr != nil {
		stripSpinner.Fail("JSON State Generation Failed")
		return "", planJsonErr
	}

	stripSpinner.Success("Secrets Stripped")

	return planExecutionPath, nil
}

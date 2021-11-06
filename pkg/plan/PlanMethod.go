package plan

import (
	"fmt"

	ux "pluralith/pkg/ux"
)

func PlanMethod(args []string) {
	// Instantiate spinners
	planSpinner := ux.NewSpinner("Running Terraform Plan", "Terraform Plan Succeeded", "Terraform Plan Failed")
	stripSpinner := ux.NewSpinner("Stripping Secrets", "Secrets Stripped", "Stripping Secrets Failed")

	// Run terraform plan
	planSpinner.Start()
	planOutput, planErr := GeneratePlan(args)
	if planErr != nil {
		planSpinner.Fail("Terraform Plan Failed")
		fmt.Println(planOutput)
		return
	} else {
		planSpinner.Success("Execution Plan Generated")
	}

	// Run terraform show + strip secrets
	stripSpinner.Start()
	_, showErr := GenerateJson(planOutput)
	if showErr != nil {
		stripSpinner.Fail("JSON State Generation Failed")
		fmt.Println(showErr)
		return
	} else {
		stripSpinner.Success("Secrets Stripped")
	}
}

package plan

import (
	"fmt"
	"os"

	"pluralith/pkg/communication"
	ux "pluralith/pkg/ux"
)

func PlanMethod(args []string, silent bool) (string, error) {
	if !silent {
		ux.PrintFormatted("⠿", []string{"blue"})
		fmt.Println(" Running plan ⇢ Inspect it in the Pluralith UI\n")
	}

	// Get working directory
	workingDir, workingErr := os.Getwd()
	if workingErr != nil {
		return "", workingErr
	}

	// Instantiate spinners
	planSpinner := ux.NewSpinner("Generating Execution Plan", "Execution Plan Generated", "Couldn't Generate Execution Plan")
	stripSpinner := ux.NewSpinner("Stripping Secrets", "Secrets Stripped", "Stripping Secrets Failed")

	// Emit plan begin update to UI
	communication.EmitUpdate(communication.UIUpdate{
		Receiver: "UI",
		Command:  "plan",
		Address:  "",
		Path:     workingDir,
		Event:    "begin",
	})

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

	// Emit plan end update to UI -> ask for confirmation
	communication.EmitUpdate(communication.UIUpdate{
		Receiver: "UI",
		Command:  "plan",
		Address:  "",
		Path:     workingDir,
		Event:    "end",
	})

	return planExecutionPath, nil
}

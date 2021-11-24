package plan

import (
	"fmt"
	"os"
	"time"

	"pluralith/pkg/comdb"
	"pluralith/pkg/ux"
)

func PlanMethod(args []string, silent bool) (string, error) {
	command := "apply"

	if !silent {
		command = "plan"
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
	comdb.PushComDBEvent(comdb.Update{
		Receiver:   "UI",
		Timestamp:  time.Now().Unix(),
		Command:    "plan",
		Event:      "begin",
		Address:    "",
		Attributes: make(map[string]interface{}),
		Path:       workingDir,
		Received:   false,
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
	comdb.PushComDBEvent(comdb.Update{
		Receiver:   "UI",
		Timestamp:  time.Now().Unix(),
		Command:    "plan",
		Event:      "end",
		Address:    "",
		Attributes: make(map[string]interface{}),
		Path:       workingDir,
		Received:   false,
	})

	comdb.PushComDBEvent(comdb.Update{
		Receiver:   "UI",
		Timestamp:  time.Now().Unix(),
		Command:    command,
		Event:      "confirm",
		Address:    "",
		Attributes: make(map[string]interface{}),
		Path:       workingDir,
		Received:   false,
	})

	return planExecutionPath, nil
}

package terraform

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"pluralith/pkg/comdb"
	"pluralith/pkg/plan"
	"pluralith/pkg/ux"
	"time"
)

func RunPlan(command string) (string, error) {
	functionName := "RunPlan"

	// Get working directory
	workingDir, workingErr := os.Getwd()
	if workingErr != nil {
		return "", fmt.Errorf("%v: %w", functionName, workingErr)
	}

	// Constructing execution plan path
	workingPlan := path.Join(workingDir, "pluralith.plan")

	// Initialize variables
	planArgs := []string{"-out", workingPlan}
	if command == "destroy" {
		planArgs = append(planArgs, "-destroy")
	}

	// Instantiate spinners
	planSpinner := ux.NewSpinner("Generating Execution Plan", "Execution Plan Generated", "Couldn't Generate Execution Plan")
	stripSpinner := ux.NewSpinner("Stripping Secrets", "Secrets Stripped", "Stripping Secrets Failed")

	planSpinner.Start()
	// Emit plan begin update to UI
	comdb.PushComDBEvent(comdb.Event{
		Receiver:  "UI",
		Timestamp: time.Now().Unix(),
		Command:   "plan",
		Type:      "begin",
		Address:   "",
		Instances: make([]interface{}, 0),
		Path:      workingDir,
		Received:  false,
	})

	// Constructing command to execute
	cmd := exec.Command("terraform", append([]string{"plan"}, planArgs...)...)

	// Defining sinks for std data
	var outputSink bytes.Buffer
	var errorSink bytes.Buffer

	// Redirecting command std data
	cmd.Stdout = &outputSink
	cmd.Stderr = &errorSink
	cmd.Stdin = os.Stdin

	// Run terraform plan
	if err := cmd.Run(); err != nil {
		planSpinner.Fail()
		return errorSink.String(), fmt.Errorf("%v: %w", functionName, err)
	}

	planSpinner.Success()
	stripSpinner.Start()

	_, planJsonErr := plan.CreatePlanJson(workingPlan)
	if planJsonErr != nil {
		stripSpinner.Fail()
		return "", fmt.Errorf("creating terraform plan json failed -> %v: %w", functionName, planJsonErr)
	}

	// Emit plan end update to UI
	comdb.PushComDBEvent(comdb.Event{
		Receiver:  "UI",
		Timestamp: time.Now().Unix(),
		Command:   "plan",
		Type:      "end",
		Address:   "",
		Instances: make([]interface{}, 0),
		Path:      workingDir,
		Received:  false,
	})

	stripSpinner.Success()

	return workingPlan, nil
}

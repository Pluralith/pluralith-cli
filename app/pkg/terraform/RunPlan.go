package terraform

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/comdb"
	"pluralith/pkg/plan"
	"pluralith/pkg/ux"
	"time"
)

func RunPlan(command string) (string, error) {
	functionName := "RunPlan"

	// Constructing execution plan path
	workingPlan := filepath.Join(auxiliary.PathInstance.WorkingPath, "pluralith.plan")

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
	comdb.PushComDBEvent(comdb.ComDBEvent{
		Receiver:  "UI",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Command:   "plan",
		Type:      "begin",
		Path:      auxiliary.PathInstance.WorkingPath,
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
		fmt.Println(errorSink.String())
		return errorSink.String(), fmt.Errorf("%v: %w", functionName, err)
	}

	planSpinner.Success()
	stripSpinner.Start()

	_, providers, planJsonErr := plan.CreatePlanJson(workingPlan)
	if planJsonErr != nil {
		stripSpinner.Fail()
		return "", fmt.Errorf("creating terraform plan json failed -> %v: %w", functionName, planJsonErr)
	}

	// Emit plan end update to UI
	comdb.PushComDBEvent(comdb.ComDBEvent{
		Receiver:  "UI",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Command:   "plan",
		Type:      "end",
		Path:      auxiliary.PathInstance.WorkingPath,
		Received:  false,
		Providers: providers,
	})

	stripSpinner.Success()

	return workingPlan, nil
}

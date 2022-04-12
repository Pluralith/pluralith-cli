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

func RunPlan(command string, tfArgs []string, silent bool) (string, error) {
	functionName := "RunPlan"

	// Constructing execution plan path
	workingPlan := filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.plan.bin")

	// Construct terraform args
	allArgs := []string{
		"plan",
		"-input=false",
		"-out=" + workingPlan,
	}

	allArgs = append(allArgs, tfArgs...)

	if command == "destroy" {
		allArgs = append(allArgs, "-destroy")
	}

	ux.PrintFormatted("→", []string{"blue", "bold"})
	ux.PrintFormatted(" Plan\n", []string{"white", "bold"})

	// Instantiate spinners
	planSpinner := ux.NewSpinner("Generating Execution Plan", "Execution Plan Generated", "Couldn't Generate Execution Plan", true)
	stripSpinner := ux.NewSpinner("Stripping Secrets", "Secrets Stripped", "Stripping Secrets Failed", true)

	planSpinner.Start()
	// Emit plan begin update to UI
	if !silent {
		comdb.PushComDBEvent(comdb.ComDBEvent{
			Receiver:  "UI",
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			Command:   "plan",
			Type:      "begin",
			Path:      auxiliary.StateInstance.WorkingPath,
			Received:  false,
		})
	}

	// Constructing command to execute
	cmd := exec.Command("terraform", allArgs...)

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

		if !silent {
			comdb.PushComDBEvent(comdb.ComDBEvent{
				Receiver:  "UI",
				Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
				Command:   "plan",
				Type:      "failed",
				Error:     errorSink.String(),
				Path:      auxiliary.StateInstance.WorkingPath,
				Received:  false,
			})
		}

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
	if !silent {
		comdb.PushComDBEvent(comdb.ComDBEvent{
			Receiver:  "UI",
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			Command:   "plan",
			Type:      "end",
			Path:      auxiliary.StateInstance.WorkingPath,
			Received:  false,
			Providers: providers,
		})
	}

	stripSpinner.Success()

	return workingPlan, nil
}

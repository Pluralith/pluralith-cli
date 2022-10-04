package terraform

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/comdb"
	"pluralith/pkg/cost"
	"pluralith/pkg/plan"
	"pluralith/pkg/ux"
	"time"
)

func RunPlan(command string, tfArgs map[string]interface{}, costArgs map[string]interface{}, silent bool) (string, error) {
	functionName := "RunPlan"

	// Instantiate spinner
	planSpinner := ux.NewSpinner("Generating Execution Plan", "Execution Plan Generated", "Couldn't Generate Execution Plan", true)

	// Create Pluralith helper directory (.pluralith)
	_, existErr := os.Stat(filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith"))
	if errors.Is(existErr, os.ErrNotExist) {
		// Create file if it doesn't exist yet
		if mkErr := os.Mkdir(filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith"), 0700); mkErr != nil {
			return "", fmt.Errorf("creating .pluralith helper directory failed -> %v: %w", functionName, mkErr)
		}
	}

	// Construct execution plan path
	workingPlan := filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.plan.bin")

	ux.PrintFormatted("\n→", []string{"blue", "bold"})
	ux.PrintFormatted(" Plan\n", []string{"white", "bold"})

	// Check if existing plan was passed
	if tfArgs["plan-file"] != "" {
		workingPlan = tfArgs["plan-file"].(string)

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

		ux.PrintFormatted("  -", []string{"blue", "bold"})
		fmt.Println(" Using Existing Execution Plan")
	} else {
		// Construct terraform args
		allArgs := []string{
			"plan",
			"-input=false",
			"-no-color",
			"-out=" + workingPlan,
		}

		// Construct arg slices for terraform
		for _, varValue := range tfArgs["var"].([]string) {
			allArgs = append(allArgs, "-var="+varValue)
		}

		for _, varFile := range tfArgs["var-file"].([]string) {
			allArgs = append(allArgs, "-var-file="+varFile)
		}

		if command == "destroy" {
			allArgs = append(allArgs, "-destroy")
		}

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

			fmt.Println(errorSink.String())
			return errorSink.String(), fmt.Errorf("%v: %w", functionName, err)
		}

		planSpinner.Success()
	}

	// Create JSON output for graphing
	_, providers, planJsonErr := plan.CreatePlanJson(workingPlan)
	if planJsonErr != nil {
		return "", fmt.Errorf("creating terraform plan json failed -> %v: %w", functionName, planJsonErr)
	}

	// Run Infracost
	if auxiliary.StateInstance.Infracost && costArgs["show-costs"] == true {
		costSpinner := ux.NewSpinner("Calculating Infrastructure Costs", "Costs Calculated", "Couldn't Calculate Costs", true)
		costSpinner.Start()

		if costErr := cost.CalculateCost(costArgs); costErr != nil {
			fmt.Println(costErr)
			costSpinner.Fail()
		}

		costSpinner.Success()
	} else {
		ux.PrintFormatted("  -", []string{"blue", "bold"})
		fmt.Println(" Cost Calculation Skipped")
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

	return workingPlan, nil
}

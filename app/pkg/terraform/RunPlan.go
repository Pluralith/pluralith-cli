package terraform

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/cost"
	"pluralith/pkg/plan"
	"pluralith/pkg/ux"
)

func RunPlan(command string, tfArgs map[string]interface{}, costArgs map[string]interface{}, localRun bool) (string, error) {
	functionName := "RunPlan"

	// Instantiate spinner
	var planSpinner = ux.NewSpinner("Generating Execution Plan", "Execution Plan Generated", "Couldn't Generate Execution Plan", true)
	if(localRun) {
		planSpinner = ux.NewSpinner("Generating Execution Plan locally", "Execution Plan Generated locally", "Couldn't Generate Execution Plan locally", true)
	}
	// Create Pluralith helper directory (.pluralith)
	_, existErr := os.Stat(filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith"))
	if errors.Is(existErr, os.ErrNotExist) {
		// Create file if it doesn't exist yet
		if mkErr := os.Mkdir(filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith"), 0700); mkErr != nil {
			return "", fmt.Errorf("creating .pluralith helper directory failed -> %v: %w", functionName, mkErr)
		}
	}

	// Construct execution plan path
	workingPlanIsJson := false
	workingPlan := filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.plan.bin")

	ux.PrintFormatted("\n→", []string{"blue", "bold"})
	ux.PrintFormatted(" Plan\n", []string{"white", "bold"})

	// Check if existing plan was passed
	if tfArgs["plan-file"] != "" {
		workingPlan = tfArgs["plan-file"].(string)
		ux.PrintFormatted("  -", []string{"blue", "bold"})
		fmt.Println(" Using Existing Execution Plan Binary File")
	} else if tfArgs["plan-file-json"] != "" {
		workingPlanIsJson = true
		workingPlan = tfArgs["plan-file-json"].(string)
		ux.PrintFormatted("  -", []string{"blue", "bold"})
		fmt.Println(" Using Existing Execution Plan JSON File")
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
			return errorSink.String(), fmt.Errorf("%v: %w", functionName, err)
		}

		planSpinner.Success()
	}

	// Create JSON output for graphing
	_, planArray, _, planJsonErr := plan.CreatePlanJson(workingPlan, workingPlanIsJson, localRun)
	if planJsonErr != nil {
		return "", fmt.Errorf("creating terraform plan json failed -> %v: %w", functionName, planJsonErr)
	}

	// Run Infracost
	if auxiliary.StateInstance.Infracost && costArgs["show-costs"] == true {

		var costSpinner = ux.NewSpinner("Calculating Infrastructure Costs", "Costs Calculated", "Couldn't Calculate Costs", true)
		if(localRun) {
			costSpinner = ux.NewSpinner("Calculating Infrastructure Costs locally", "Costs Calculated locally", "Couldn't Calculate Costs locally", true)
		}
		costSpinner.Start()

		if costErr := cost.CalculateCost(costArgs, planArray); costErr != nil {
			fmt.Println(costErr)
			costSpinner.Fail()
		}

		costSpinner.Success()
	} else {
		ux.PrintFormatted("  -", []string{"blue", "bold"})
		fmt.Println(" Cost Calculation Skipped")
	}

	return workingPlan, nil
}

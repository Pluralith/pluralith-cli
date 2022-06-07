package cost

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func CalculateCost(costArgs map[string]interface{}) error {
	functionName := "CalculateCost"

	costSpinner := ux.NewSpinner("Calculating Infrastructure Costs", "Costs Calculated", "Couldn't Calculate Costs", true)
	costSpinner.Start()

	planJsonPath := filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.state.json")

	allArgs := []string{
		"breakdown",
		"--path=" + planJsonPath,
		"--out-file=.pluralith/pluralith.costs.json",
		"--format=json",
	}

	if costArgs["usage-file-path"] != nil {
		allArgs = append(allArgs, "--usage-file="+costArgs["usage-file-path"].(string))
	}

	costCmd := exec.Command("infracost", allArgs...)

	// Defining sinks for std data
	var outputSink bytes.Buffer
	var errorSink bytes.Buffer

	// Redirecting command std data
	costCmd.Stdout = &outputSink
	costCmd.Stderr = &errorSink
	costCmd.Stdin = os.Stdin

	if runErr := costCmd.Run(); runErr != nil {
		costSpinner.Fail()
		fmt.Println(errorSink.String())
		return fmt.Errorf("running infracost breakdown failed -> %v: %w", functionName, runErr)
	}

	costSpinner.Success()
	return nil
}

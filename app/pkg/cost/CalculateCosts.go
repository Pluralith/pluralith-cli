package cost

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
)

func CalculateCost(costArgs map[string]interface{}, planArray []string) error {
	functionName := "CalculateCost"

	costArray := []string{}

	// planJsonPath := filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.state.json")
	cachePath := filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.cache.json")
	costPath := filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.costs.json")

	// Run on each plan in plan array
	for _, plan := range planArray {
		// Write current plan to disk
		if writeErr := os.WriteFile(cachePath, []byte(plan), 0700); writeErr != nil {
			return fmt.Errorf("%v: %w", functionName, writeErr)
		}

		allArgs := []string{
			"breakdown",
			"--path=" + cachePath,
			"--out-file=" + cachePath,
			"--format=json",
		}

		if costArgs["cost-usage-file"] != nil {
			allArgs = append(allArgs, "--usage-file="+costArgs["cost-usage-file"].(string))
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
			fmt.Println(errorSink.String())
			return fmt.Errorf("running infracost breakdown failed -> %v: %w", functionName, runErr)
		}

		// Read cost file
		costFile, readErr := os.ReadFile(cachePath)
		if readErr != nil {
			return fmt.Errorf("failed to open json plan -> %v: %w", functionName, readErr)
		}

		costArray = append(costArray, string(costFile))
	}

	// Merge cost objects
	costString := "["
	for index, costJson := range costArray {
		costString += costJson
		if index < len(costArray)-1 {
			costString += ","
		}
	}
	costString += "]"

	// Write array of cost files to disk
	if writeErr := os.WriteFile(costPath, []byte(costString), 0700); writeErr != nil {
		return fmt.Errorf("%v: %w", functionName, writeErr)
	}

	return nil
}

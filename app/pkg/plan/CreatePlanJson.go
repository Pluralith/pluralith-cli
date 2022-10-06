package plan

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"pluralith/pkg/strip"
	"pluralith/pkg/ux"
)

func ConvertBinaryPlanToJson(planPath string) (string, error) {
	functionName := "ConvertBinaryPlanToJson"

	// Constructing command to execute
	cmd := exec.Command("terraform", append([]string{"show", "-json", planPath})...)

	// Defining sinks for std data
	var outputSink bytes.Buffer
	var errorSink bytes.Buffer

	// Redirecting command std data
	cmd.Stdout = &outputSink
	cmd.Stderr = &errorSink
	cmd.Stdin = os.Stdin

	// Run terraform command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("terraform command failed -> %v: %w", functionName, err)
	}

	return outputSink.String(), nil
}

func CreatePlanJson(planPath string, isJson bool) (string, []string, error) {
	functionName := "CreatePlanJson"

	// Instantiate spinner
	stripSpinner := ux.NewSpinner("Stripping Secrets", "Secrets Stripped", "Stripping Secrets Failed", true)
	stripSpinner.Start()

	// Get working directory
	workingDir, workingErr := os.Getwd()
	if workingErr != nil {
		stripSpinner.Fail()
		return "", []string{}, fmt.Errorf("%v: %w", functionName, workingErr)
	}

	// Construct file path for stripped state
	strippedPath := filepath.Join(workingDir, ".pluralith", "pluralith.state.json")

	planJsonString := ""
	if isJson {
		// Read json file
		file, readErr := ioutil.ReadFile(planPath)
		if readErr != nil {
			stripSpinner.Fail()
			return "", []string{}, fmt.Errorf("failed to open json plan -> %v: %w", functionName, readErr)
		}
		planJsonString = string(file)
	} else {
		// Convert binary plan to json
		planJsonStringTemp, convertErr := ConvertBinaryPlanToJson(planPath)
		if convertErr != nil {
			stripSpinner.Fail()
			return "", []string{}, fmt.Errorf("failed to strip secrets -> %v: %w", functionName, convertErr)
		}
		planJsonString = planJsonStringTemp
	}

	// Strip secrets from plan state json
	strippedState, stripErr := strip.StripSecrets(planJsonString)
	if stripErr != nil {
		stripSpinner.Fail()
		return "", []string{}, fmt.Errorf("failed to strip secrets -> %v: %w", functionName, stripErr)
	}

	// Fetch providers present in state
	providers, providerErr := FetchProviders(planJsonString)
	if providerErr != nil {
		stripSpinner.Fail()
		return "", []string{}, fmt.Errorf("failed to fetch providers -> %v: %w", functionName, providerErr)
	}

	// Write stripped state to file
	if writeErr := ioutil.WriteFile(strippedPath, []byte(strippedState), 0700); writeErr != nil {
		stripSpinner.Fail()
		return "", providers, fmt.Errorf("%v: %w", functionName, writeErr)
	}

	stripSpinner.Success()

	// Return path to execution plan
	return strippedPath, providers, nil
}

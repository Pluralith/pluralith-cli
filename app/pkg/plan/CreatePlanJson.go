package plan

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"pluralith/pkg/strip"
)

func CreatePlanJson(planPath string) (string, []string, error) {
	functionName := "CreatePlanJson"

	// Get working directory
	workingDir, workingErr := os.Getwd()
	if workingErr != nil {
		return "", []string{}, fmt.Errorf("%v: %w", functionName, workingErr)
	}
	// Construct file path for stripped state
	strippedPath := path.Join(workingDir, "pluralith.state.stripped")

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
		return errorSink.String(), []string{}, fmt.Errorf("terraform command failed -> %v: %w", functionName, err)
	}

	// Strip secrets from plan state json
	strippedState, stripErr := strip.StripSecrets(outputSink.String())
	if stripErr != nil {
		return "", []string{}, fmt.Errorf("failed to strip secrets -> %v: %w", functionName, stripErr)
	}

	// Fetch providers present in state
	providers, providerErr := FetchProviders(outputSink.String())
	if providerErr != nil {
		return "", []string{}, fmt.Errorf("failed to fetch providers -> %v: %w", functionName, providerErr)
	}

	// Write stripped state to file
	if writeErr := ioutil.WriteFile(strippedPath, []byte(strippedState), 0700); writeErr != nil {
		return "", providers, fmt.Errorf("%v: %w", functionName, writeErr)
	}

	// Return path to execution plan
	return strippedPath, providers, nil
}

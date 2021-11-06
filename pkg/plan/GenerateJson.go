package plan

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	strip "pluralith/pkg/strip"
)

func GenerateJson(planPath string) (string, error) {
	// Get working directory
	workingDir, workingErr := os.Getwd()
	if workingErr != nil {
		return "", workingErr
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
		return errorSink.String(), errors.New("terraform command failed")
	}

	// Strip secrets from plan state json
	strippedState, stripErr := strip.StripSecrets(outputSink.String(), []string{}, "gatewatch")
	if stripErr != nil {
		return "", stripErr
	}

	// Write stripped state to file
	ioutil.WriteFile(strippedPath, []byte(strippedState), 0644)

	// Return path to execution plan
	return strippedPath, nil
}

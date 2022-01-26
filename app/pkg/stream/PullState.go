package stream

import (
	"bytes"
	"fmt"
	"os/exec"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/strip"
)

func PullState() (map[string]interface{}, error) {
	functionName := "FetchState"

	// Constructing command to execute
	cmd := exec.Command("terraform", "state", "pull")

	// Defining sinks for std data
	var outputSink bytes.Buffer
	var errorSink bytes.Buffer
	// Redirecting command std data
	cmd.Stdout = &outputSink
	cmd.Stderr = &errorSink

	// Running terraform command
	if err := cmd.Run(); err != nil {
		fmt.Println(errorSink.String())
		// Handling error
		return make(map[string]interface{}), fmt.Errorf("%v: %w", functionName, err)
	}

	// Parse state received from terraform state pull
	parsedState, parseErr := auxiliary.ParseJson(outputSink.String()) // FAILS ON DESTROY -> FIX
	if parseErr != nil {
		return make(map[string]interface{}), fmt.Errorf("parsing terraform state failed -> %v: %w", functionName, parseErr)
	}

	// Filter secrets according to config
	strip.ReplaceSensitive(parsedState)

	return parsedState, nil
}
